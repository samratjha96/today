import { useQuery } from "@tanstack/react-query";
import { BACKEND_URL } from "../lib/constants";

export interface TickerData {
  ticker: string;
  todaysPrice: number;
  dayChange: number;
  weekChange: number;
  yearChange: number;
}

const MOCK_DATA: TickerData[] = [
  {
    ticker: "SPY",
    todaysPrice: 478.25,
    dayChange: 0.75,
    weekChange: 2.15,
    yearChange: 24.32
  },
  {
    ticker: "QQQ",
    todaysPrice: 424.75,
    dayChange: -0.45,
    weekChange: 1.85,
    yearChange: 55.18
  },
  {
    ticker: "VTI",
    todaysPrice: 235.90,
    dayChange: 0.55,
    weekChange: 1.95,
    yearChange: 22.45
  },
  {
    ticker: "VT",
    todaysPrice: 98.45,
    dayChange: -0.25,
    weekChange: 1.15,
    yearChange: 18.75
  },
  {
    ticker: "SCHD",
    todaysPrice: 76.80,
    dayChange: 0.35,
    weekChange: 1.45,
    yearChange: 15.25
  },
  {
    ticker: "REIT",
    todaysPrice: 22.15,
    dayChange: -0.85,
    weekChange: -1.25,
    yearChange: -5.45
  },
  {
    ticker: "IAU",
    todaysPrice: 38.90,
    dayChange: 0.15,
    weekChange: 2.85,
    yearChange: 8.75
  }
];

const fetchTickerData = async (): Promise<TickerData[]> => {
  // Check if we should use mock data
  if (import.meta.env.VITE_API_MODE === "mock") {
    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 500));
    return MOCK_DATA;
  }

  const url = `${BACKEND_URL}/tickers`;

  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Error fetching Tickers data: ${response.statusText}`);
    }
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Error fetching Tickers data:", error);
    return [];
  }
};

export const useTickerData = () => {
  return useQuery({
    queryKey: ["tickers"],
    queryFn: fetchTickerData,
    refetchInterval: 30000, // 30 seconds
    staleTime: 25000, // 25 seconds - data considered fresh for 25 seconds
    gcTime: 120000, // 2 minutes - keep unused data in cache for 2 minutes
    retry: 2, // Retry failed requests twice
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000), // Exponential backoff
    refetchOnWindowFocus: false, // Don't refetch when window regains focus
    refetchOnReconnect: true, // Refetch on reconnection
  });
};

export const getMarketSentiment = (data: TickerData[]) => {
  if (!data || data.length === 0) return null;
  
  const negativeCount = data.filter(ticker => ticker.dayChange < 0).length;
  return negativeCount > data.length / 2 ? "bearish" : "bullish";
};
