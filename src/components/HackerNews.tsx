import { useQuery } from "@tanstack/react-query";
interface HNPost {
  title: string;
  points: number;
  author: string;
  comments: number;
  url: string;
}
const mockHNData: HNPost[] = [
  {
    title: "Rust Is The Future of JavaScript Infrastructure",
    points: 789,
    author: "rustdev",
    comments: 234,
    url: "https://example.com/rust-js",
  },
  {
    title: "The Future of Web Development: WASM and Beyond",
    points: 567,
    author: "webdev",
    comments: 189,
    url: "https://example.com/wasm",
  },
  {
    title: "New AI Model Achieves Human-Level Performance",
    points: 432,
    author: "airesearcher",
    comments: 156,
    url: "https://example.com/ai-model",
  },
  {
    title: "Understanding Modern CPU Architecture",
    points: 345,
    author: "cpuexpert",
    comments: 123,
    url: "https://example.com/cpu",
  },
];
const fetchHNData = async (): Promise<HNPost[]> => {
  return Promise.resolve(mockHNData);
};
export const HackerNews = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["hackernews"],
    queryFn: fetchHNData,
    refetchInterval: 300000,
  });
  if (isLoading) {
    return (
      <div className="animate-pulse p-4 bg-terminal-secondary rounded-lg">
        <div className="h-40 bg-terminal-muted rounded" />
      </div>
    );
  }
  if (error || !data || data.length === 0) {
    return (
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
  }
  return (
    <div className="border border-terminal-accent/20 rounded-lg overflow-hidden">
      <div className="divide-y divide-terminal-accent/20">
        {data.map((post, index) => (
          <div
            key={index}
            className="p-4 bg-terminal-bg/50 hover:bg-terminal-secondary/50 transition-colors"
          >
            <div className="flex flex-col space-y-2">
              <a
                href={post.url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-[#33C3F0] font-mono text-sm hover:text-[#1EAEDB] transition-colors line-clamp-2"
              >
                {post.title}
              </a>
              <div className="flex items-center space-x-4 text-terminal-text/60 font-mono text-xs">
                <span>👤 {post.author}</span>
                <span>⭐ {post.points}</span>
                <span>💬 {post.comments}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
