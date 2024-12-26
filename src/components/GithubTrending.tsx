import { useQuery } from "@tanstack/react-query";
import { cn } from "@/lib/utils";
interface GithubRepo {
  name: string;
  author: string;
  description: string;
  stars: number;
  language: string;
}
const mockGithubData: GithubRepo[] = [
  {
    name: "awesome-ai",
    author: "microsoft",
    description: "A curated list of AI tools and frameworks",
    stars: 12500,
    language: "Python"
  },
  {
    name: "next-auth",
    author: "nextauthjs",
    description: "Authentication for Next.js",
    stars: 9800,
    language: "TypeScript"
  },
  {
    name: "rust-book",
    author: "rust-lang",
    description: "The Rust Programming Language",
    stars: 8900,
    language: "Rust"
  },
  {
    name: "deno",
    author: "denoland",
    description: "A modern runtime for JavaScript and TypeScript",
    stars: 7600,
    language: "TypeScript"
  }
];
const fetchGithubData = async (): Promise<GithubRepo[]> => {
  return Promise.resolve(mockGithubData);
};
export const GithubTrending = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["github"],
    queryFn: fetchGithubData,
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
          <span className="text-xl">🐙</span>
        </div>
        <h3 className="text-terminal-text font-mono text-lg mb-2">No Github Trends</h3>
        <p className="text-terminal-text/60 font-mono text-sm">
          Unable to fetch trending repositories. Please try again later.
        </p>
      </div>
    );
  }
  return (
    <div className="border border-terminal-accent/20 rounded-lg overflow-hidden">
      <div className="divide-y divide-terminal-accent/20">
        {data.map((repo, index) => (
          <div
            key={index}
            className="p-4 bg-terminal-bg/50 hover:bg-terminal-secondary/50 transition-colors"
          >
            <div className="flex flex-col space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-[#33C3F0] font-mono text-sm">
                  {repo.author}/{repo.name}
                </span>
                <span className="text-terminal-text/60 font-mono text-xs">
                  ⭐ {repo.stars.toLocaleString()}
                </span>
              </div>
              <p className="text-terminal-text/80 font-mono text-xs line-clamp-2">
                {repo.description}
              </p>
              <span className="text-terminal-text/40 font-mono text-xs">
                {repo.language}
              </span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};