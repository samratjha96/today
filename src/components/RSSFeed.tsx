import { useRSSData } from "@/hooks/use-rss-data";
import { withPagination } from "@/components/hoc/withPagination";

interface RSSItem {
  title: string;
  link: string;
  source: string;
}

const PaginatedRSSFeed = withPagination<RSSItem>(() => null);

export const RSSFeed = () => {
  const { data, isLoading, error } = useRSSData();

  const renderItem = (item: RSSItem) => (
    <div
      key={item.link}
      className="px-4 py-2 bg-terminal-bg/50 hover:bg-terminal-secondary/50 transition-colors min-h-[60px] flex items-center"
    >
      <div className="flex items-center justify-between w-full">
        <a
          href={item.link}
          target="_blank"
          rel="noopener noreferrer"
          className="text-[#33C3F0] font-mono text-sm hover:underline line-clamp-1 flex-1 mr-4"
        >
          {item.title}
        </a>
        <span className="text-terminal-text/60 font-mono text-xs whitespace-nowrap">
          {item.source}
        </span>
      </div>
    </div>
  );

  const errorComponent = (
    <div className="p-8 text-center border border-terminal-accent/20 rounded-lg bg-terminal-bg/50">
      <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-terminal-secondary mb-4">
        <span className="text-xl">📰</span>
      </div>
      <h3 className="text-terminal-text font-mono text-lg mb-2">
        No News Available
      </h3>
      <p className="text-terminal-text/60 font-mono text-sm">
        Unable to fetch news articles at this time. Please check back later.
      </p>
    </div>
  );

  const loadingComponent = (
    <div className="animate-pulse p-4 bg-terminal-secondary rounded-lg">
      <div className="h-40 bg-terminal-muted rounded" />
    </div>
  );

  return (
    <PaginatedRSSFeed
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
