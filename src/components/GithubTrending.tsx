import { useQuery } from "@tanstack/react-query";
import { withPagination } from "./hoc/withPagination";

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

const PaginatedGithubTrending = withPagination<GithubRepo>(() => null);

export const GithubTrending = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["github"],
    queryFn: fetchGithubData,
    refetchInterval: 300000, // Refetch every 5 minutes
  });

  const renderItem = (repo: GithubRepo) => (
    <div
      key={`${repo.author}/${repo.name}`}
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
            <span className="text-green-400">+{repo.currentPeriodStars}</span>
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
  );

  const errorComponent = (
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

  const loadingComponent = (
    <div className="animate-pulse p-4 bg-terminal-secondary rounded-lg">
      <div className="h-40 bg-terminal-muted rounded" />
    </div>
  );

  return (
    <PaginatedGithubTrending
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
