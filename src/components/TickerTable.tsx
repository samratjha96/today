import { useTickerData } from "@/hooks/use-ticker-data";
import { withPagination } from "./hoc/withPagination";
import { TickerCard } from "./TickerCard";

const PaginatedTickerTable = withPagination<any>(() => null);

export const TickerTable = () => {
  const { data, isLoading, error } = useTickerData();

  const loadingComponent = (
    <div className="animate-pulse space-y-4">
      {[...Array(3)].map((_, i) => (
        <div key={i} className="h-24 bg-terminal-secondary rounded-lg" />
      ))}
    </div>
  );

  const errorComponent = (
    <div className="p-8 text-center border border-terminal-accent/20 rounded-lg bg-terminal-bg/50">
      <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-terminal-secondary mb-4">
        <span className="text-xl">📊</span>
      </div>
      <h3 className="text-terminal-text font-mono text-lg mb-2">
        No Market Data Available
      </h3>
      <p className="text-terminal-text/60 font-mono text-sm">
        Unable to fetch market data at this time. Please try again later.
      </p>
    </div>
  );

  return (
    <PaginatedTickerTable
      data={data}
      isLoading={isLoading}
      error={error}
      itemsPerPage={8}
      renderItem={(ticker) => (
        <TickerCard
          key={ticker.ticker}
          ticker={ticker.ticker}
          todaysPrice={ticker.todaysPrice}
          dayChange={ticker.dayChange}
          weekChange={ticker.weekChange}
          yearChange={ticker.yearChange}
        />
      )}
      errorComponent={errorComponent}
      loadingComponent={loadingComponent}
      containerClassName="animate-fadeIn space-y-6"
      itemsContainerClassName="space-y-4"
      paginationVariant="numbered"
    />
  );
};
