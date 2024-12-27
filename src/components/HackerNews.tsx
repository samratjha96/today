import { useQuery } from "@tanstack/react-query";
import { withPagination } from "./hoc/withPagination";

interface HNStory {
  by: string;
  descendants: number;
  id: number;
  score: number;
  time: number;
  title: string;
  type: string;
  url: string;
}

const fetchHNData = async (): Promise<HNStory[]> => {
  const baseUrl = import.meta.env.VITE_BACKEND_URL || "/api";
  const response = await fetch(`${baseUrl}/hackernews/top`);
  if (!response.ok) {
    throw new Error("Failed to fetch Hacker News stories");
  }
  return response.json();
};

const PaginatedHackerNews = withPagination<HNStory>(() => null);

export const HackerNews = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["hackernews"],
    queryFn: fetchHNData,
    refetchInterval: 300000, // Refetch every 5 minutes
  });

  const renderItem = (story: HNStory) => (
    <div
      key={story.id}
      className="p-4 bg-terminal-bg/50 hover:bg-terminal-secondary/50 transition-colors min-h-[120px] flex flex-col justify-between"
    >
      <div className="flex flex-col space-y-4">
        <div className="flex items-start justify-between">
          <a
            href={story.url}
            target="_blank"
            rel="noopener noreferrer"
            className="text-[#33C3F0] font-mono text-sm hover:underline line-clamp-2 flex-1 mr-4"
          >
            {story.title}
          </a>
          <div className="flex items-center gap-3 text-terminal-text/60 font-mono text-xs whitespace-nowrap">
            <span>⭐ {story.score}</span>
            <span>💬 {story.descendants}</span>
          </div>
        </div>
        <div className="flex items-center justify-between mt-auto">
          <div className="flex items-center gap-2 text-terminal-text/40 font-mono text-xs">
            <span>by {story.by}</span>
            <span>•</span>
            <span>{new Date(story.time * 1000).toLocaleString()}</span>
          </div>
          <a
            href={`https://news.ycombinator.com/item?id=${story.id}`}
            target="_blank"
            rel="noopener noreferrer"
            className="text-terminal-text/60 font-mono text-xs hover:text-terminal-text"
          >
            discuss →
          </a>
        </div>
      </div>
    </div>
  );

  const errorComponent = (
    <div className="p-8 text-center border border-terminal-accent/20 rounded-lg bg-terminal-bg/50">
      <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-terminal-secondary mb-4">
        <span className="text-xl">📰</span>
      </div>
      <h3 className="text-terminal-text font-mono text-lg mb-2">
        No HN Stories
      </h3>
      <p className="text-terminal-text/60 font-mono text-sm">
        Unable to fetch Hacker News stories. Please try again later.
      </p>
    </div>
  );

  const loadingComponent = (
    <div className="animate-pulse p-4 bg-terminal-secondary rounded-lg">
      <div className="h-40 bg-terminal-muted rounded" />
    </div>
  );

  return (
    <PaginatedHackerNews
      data={data}
      isLoading={isLoading}
      error={error}
      itemsPerPage={5}
      renderItem={renderItem}
      errorComponent={errorComponent}
      loadingComponent={loadingComponent}
    />
  );
};
