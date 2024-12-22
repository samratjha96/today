import { useQuery } from "@tanstack/react-query";

interface RSSItem {
  source: string;
  title: string;
  link: string;
}

const mockRSSData: RSSItem[] = [
  {
    source: "TechCrunch",
    title: "LimeWire AI Studio Review 2023: Details, Pricing & Features",
    link: "https://techcruncher.blogspot.com/2023/12/limewire-ai-studio-review-2023-details.html"
  },
  {
    source: "TechCrunch",
    title: "Top 10 AI Tools in 2023 That Will Make Your Life Easier",
    link: "https://techcruncher.blogspot.com/2023/01/top-10-ai-tools-in-2023-that-will-make.html"
  },
  {
    source: "Wired",
    title: "9 Best French Presses (2024): Plastic, Glass, Stainless Steel, Travel",
    link: "https://www.wired.com/gallery/best-french-presses/"
  },
  {
    source: "The Verge",
    title: "10 great shows to stream on Amazon Prime Video from 2024",
    link: "https://www.theverge.com/24302668/amazon-prime-video-best-2024-shows-streaming"
  },
  {
    source: "Ars Technica",
    title: "Ars Technica's top 20 video games of 2024",
    link: "https://arstechnica.com/gaming/2024/12/ars-technicas-top-20-video-games-of-2024/"
  }
];

const fetchRSSData = async (): Promise<RSSItem[]> => {
  // Simulating API call
  return new Promise((resolve) => {
    setTimeout(() => resolve(mockRSSData), 1000);
  });
};

export const RSSFeed = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["rss"],
    queryFn: fetchRSSData,
    refetchInterval: 300000, // Refetch every 5 minutes
  });

  if (isLoading) {
    return (
      <div className="animate-pulse p-4 bg-terminal-secondary rounded-lg">
        <div className="h-60 bg-terminal-muted rounded" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4 text-red-500 bg-terminal-secondary rounded-lg">
        Error loading RSS feed
      </div>
    );
  }

  return (
    <div className="border border-terminal-accent/20 rounded-lg overflow-hidden animate-fadeIn">
      <div className="bg-terminal-secondary px-6 py-3">
        <h2 className="text-terminal-text font-mono text-sm">Top Tech and Finance News (from RSS):</h2>
      </div>
      <div className="divide-y divide-terminal-accent/20">
        {data?.map((item, index) => (
          <div
            key={index}
            className="px-6 py-4 bg-terminal-bg/50 hover:bg-terminal-secondary/50 transition-colors"
          >
            <div className="flex items-start space-x-4">
              <span className="text-terminal-text/60 font-mono text-xs whitespace-nowrap">
                {item.source}
              </span>
              <a
                href={item.link}
                target="_blank"
                rel="noopener noreferrer"
                className="text-terminal-text font-mono text-sm hover:text-terminal-accent transition-colors line-clamp-2"
              >
                {item.title}
              </a>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};