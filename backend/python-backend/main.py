from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import pandas as pd
import feedparser
from bs4 import BeautifulSoup
from fastapi.middleware.cors import CORSMiddleware
import os
from typing import Optional

app = FastAPI()

# Get allowed origins from environment variable or use default
ALLOWED_HOSTS = os.getenv("ALLOWED_HOSTS", "localhost, today.techbrohomelab.xyz").split(",")
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
