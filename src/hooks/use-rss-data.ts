import { useQuery } from "@tanstack/react-query";
import { BACKEND_URL } from "@/lib/constants";

export interface RSSItem {
  source: string;
  title: string;
  link: string;
}

const fetchRSSData = async (): Promise<RSSItem[]> => {
  const url = `${BACKEND_URL}/news`;

  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Error fetching RSS data: ${response.statusText}`);
    }
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Error fetching RSS data:", error);
    return [];
  }
};

export const useRSSData = () => {
  return useQuery({
    queryKey: ["rss"],
    queryFn: fetchRSSData,
    refetchInterval: 300000, // 5 minutes
  });
};
