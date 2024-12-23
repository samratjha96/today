from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import yfinance as yf
import pandas as pd
import feedparser
from bs4 import BeautifulSoup
from fastapi.middleware.cors import CORSMiddleware
import os

app = FastAPI()

# Get allowed origins from environment variable or use default
ALLOWED_HOSTS = os.getenv("ALLOWED_HOSTS", "today.techbrohomelab.xyz").split(",")
origins = [
    f"https://{host.strip()}" for host in ALLOWED_HOSTS
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    allow_headers=["*"],
    expose_headers=["*"],
    max_age=3600,  # Cache preflight requests for 1 hour
)

class TickerData(BaseModel):
    ticker: str
    todaysPrice: float | None
    dayChange: float | None
    weekChange: float | None
    yearChange: float | None

class RSSItem(BaseModel):
    source: str
    title: str
    link: str

RSS_FEEDS = {
    "TechCrunch": "http://feeds.feedburner.com/TechCrunch/",
    "Wired": "https://www.wired.com/feed/rss",
    "The Verge": "https://www.theverge.com/rss/index.xml",
    "Ars Technica": "http://feeds.arstechnica.com/arstechnica/index",
}

DEFAULT_TICKERS = [
    "SPY",
    "QQQ",
    "VTI",
    "VT",
    "SCHD",
    "REIT",
    "IAU",
]

def get_etf_data(etf_tickers):
    try:
        data = yf.download(etf_tickers, period="1y", progress=False)
        return data
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error downloading data: {e}")

def calculate_changes(data, ticker):
    if data.empty or len(data) < 2:
        return None, None, None

    try:
        today = data["Close"][ticker].iloc[-1]
        yesterday = data["Close"][ticker].iloc[-2]
        first_day = data["Close"][ticker].iloc[0]
        five_days_ago = data["Close"][ticker].iloc[-6]
    except (KeyError, IndexError):
        return None, None, None

    change_24h_percent = ((today - yesterday) / yesterday) * 100 if pd.notna(yesterday) and yesterday != 0 else None
    change_1y_percent = ((today - first_day) / first_day) * 100 if pd.notna(first_day) and first_day != 0 else None
    change_5d_percent = ((today - five_days_ago) / five_days_ago) * 100 if pd.notna(five_days_ago) and five_days_ago != 0 else None

    return today, change_24h_percent, change_5d_percent, change_1y_percent

def get_top_news_from_rss(rss_url):
    try:
        feed = feedparser.parse(rss_url)
        if feed.status == 200:
            return feed.entries[:5]
        else:
            raise HTTPException(status_code=feed.status, detail=f"Error fetching RSS feed: HTTP status {feed.status} for {rss_url}")
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error parsing RSS feed for {rss_url}: {e}")

def clean_html(html_string):
    if html_string is None:
        return ""
    soup = BeautifulSoup(html_string, "html.parser")
    return soup.get_text(separator=" ", strip=True)

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

@app.get("/tickers", response_model=list[TickerData])
async def get_tickers():
    ticker_list = DEFAULT_TICKERS
    data = get_etf_data(ticker_list)

    ticker_data_list = []
    for ticker in ticker_list:
        if ("Close", ticker) not in data.columns or data[[("Close", ticker)]].empty or data[[("Close", ticker)]].isnull().all().all():
            ticker_data_list.append(TickerData(ticker=ticker, todaysPrice=None, dayChange=None, weekChange=None, yearChange=None))
            continue
        today_price, change_24h, change_5d, change_1y = calculate_changes(data, ticker)
        ticker_data_list.append(TickerData(ticker=ticker, todaysPrice=today_price, dayChange=change_24h, weekChange=change_5d, yearChange=change_1y))

    return ticker_data_list

@app.get("/news", response_model=list[RSSItem])
async def get_news():
    all_news = []
    for source, url in RSS_FEEDS.items():
        news = get_top_news_from_rss(url)
        if news:
            for entry in news:
                title = entry.get("title", "No title available").strip()
                link = entry.get("link")
                if link: # only append if a link is available
                    all_news.append(RSSItem(source=source, title=title, link=link))

    return all_news

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
