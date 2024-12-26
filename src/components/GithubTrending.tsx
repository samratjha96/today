import { useQuery } from "@tanstack/react-query";
import { cn } from "@/lib/utils";
import { useState } from "react";

interface GithubContributor {
  username: string;
  href: string;
  avatar: string;
}

interface GithubRepo {
  author: string;
  name: string;
  avatar: string;
  url: string;
  description: string;
  language: string;
  languageColor: string;
  stars: number;
  forks: number;
  currentPeriodStars: number;
  builtBy: GithubContributor[];
}

const fetchGithubData = async (): Promise<GithubRepo[]> => {
  const baseUrl = import.meta.env.VITE_BACKEND_URL || "/api";
  const response = await fetch(`${baseUrl}/github/trending`);
  if (!response.ok) {
    throw new Error("Failed to fetch GitHub trending repositories");
  }
  return response.json();
};

const ITEMS_PER_PAGE = 5;

export const GithubTrending = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const { data, isLoading, error } = useQuery({
    queryKey: ["github"],
    queryFn: fetchGithubData,
    refetchInterval: 300000, // Refetch every 5 minutes
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
        <h3 className="text-terminal-text font-mono text-lg mb-2">
          No Github Trends
        </h3>
        <p className="text-terminal-text/60 font-mono text-sm">
          Unable to fetch trending repositories. Please try again later.
        </p>
      </div>
    );
  }

  const totalPages = Math.ceil(data.length / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const endIndex = startIndex + ITEMS_PER_PAGE;
  const currentItems = data.slice(startIndex, endIndex);

  return (
    <div className="border border-terminal-accent/20 rounded-lg overflow-hidden">
      <div className="divide-y divide-terminal-accent/20 max-h-[600px] overflow-y-auto">
        {currentItems.map((repo, index) => (
          <div
            key={index}
            className="p-4 bg-terminal-bg/50 hover:bg-terminal-secondary/50 transition-colors h-[120px] flex flex-col justify-between"
          >
            <div className="flex flex-col h-full">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-2 min-w-0 flex-1">
                  <img
                    src={repo.avatar}
                    alt={`${repo.author}'s avatar`}
                    className="w-4 h-4 rounded-full flex-shrink-0"
                  />
                  <a
                    href={repo.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-[#33C3F0] font-mono text-sm hover:underline truncate"
                  >
                    {repo.author}/{repo.name}
                  </a>
                </div>
                <div className="flex items-center gap-3 text-terminal-text/60 font-mono text-xs flex-shrink-0 ml-2">
                  <span>⭐ {repo.stars.toLocaleString()}</span>
                  <span>🔱 {repo.forks.toLocaleString()}</span>
                  <span className="text-green-400">
                    +{repo.currentPeriodStars}
                  </span>
                </div>
              </div>
              <p className="text-terminal-text/80 font-mono text-xs line-clamp-2 mb-2 flex-1">
                {repo.description}
              </p>
              <div className="flex items-center justify-between mt-auto">
                <div className="flex items-center gap-1.5">
                  <span
                    className="w-3 h-3 rounded-full flex-shrink-0"
                    style={{ backgroundColor: repo.languageColor }}
                  />
                  <span className="text-terminal-text/40 font-mono text-xs">
                    {repo.language}
                  </span>
                </div>
                <div className="flex -space-x-2">
                  {repo.builtBy.slice(0, 3).map((contributor, idx) => (
                    <a
                      key={idx}
                      href={contributor.href}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="hover:z-10"
                    >
                      <img
                        src={contributor.avatar}
                        alt={`${contributor.username}'s avatar`}
                        className="w-5 h-5 rounded-full border border-terminal-bg"
                      />
                    </a>
                  ))}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
      {totalPages > 1 && (
        <div className="flex items-center justify-between p-4 border-t border-terminal-accent/20 bg-terminal-bg/50">
          <button
            onClick={() => setCurrentPage((prev) => Math.max(1, prev - 1))}
            disabled={currentPage === 1}
            className={cn(
              "px-3 py-1 rounded font-mono text-xs",
              currentPage === 1
                ? "text-terminal-text/40 cursor-not-allowed"
                : "text-terminal-text hover:bg-terminal-secondary"
            )}
          >
            Previous
          </button>
          <span className="text-terminal-text/60 font-mono text-xs">
            Page {currentPage} of {totalPages}
          </span>
          <button
            onClick={() =>
              setCurrentPage((prev) => Math.min(totalPages, prev + 1))
            }
            disabled={currentPage === totalPages}
            className={cn(
              "px-3 py-1 rounded font-mono text-xs",
              currentPage === totalPages
                ? "text-terminal-text/40 cursor-not-allowed"
                : "text-terminal-text hover:bg-terminal-secondary"
            )}
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
};
