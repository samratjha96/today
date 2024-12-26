import { cn } from "@/lib/utils";
import { useState } from "react";
import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from "./ui/pagination";
import { useTickerData } from "@/hooks/use-ticker-data";
import { TickerCard } from "./TickerCard";

const ITEMS_PER_PAGE = 8;

export const TickerTable = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const { data, isLoading, error } = useTickerData();

  if (isLoading) {
    return (
      <div className="animate-pulse space-y-4">
        {[...Array(3)].map((_, i) => (
          <div key={i} className="h-24 bg-terminal-secondary rounded-lg" />
        ))}
      </div>
    );
  }

  if (error || !data || data.length === 0) {
    return (
      <div className="p-8 text-center border border-terminal-accent/20 rounded-lg bg-terminal-bg/50">
        <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-terminal-secondary mb-4">
          <span className="text-xl">📊</span>
        </div>
        <h3 className="text-terminal-text font-mono text-lg mb-2">No Market Data Available</h3>
        <p className="text-terminal-text/60 font-mono text-sm">
          Unable to fetch market data at this time. Please try again later.
        </p>
      </div>
    );
  }

  const totalPages = Math.ceil((data?.length || 0) / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const paginatedData = data?.slice(startIndex, startIndex + ITEMS_PER_PAGE);

  return (
    <div className="animate-fadeIn space-y-6">
      <div className="space-y-4">
        {paginatedData?.map((ticker) => (
          <TickerCard
            key={ticker.ticker}
            ticker={ticker.ticker}
            todaysPrice={ticker.todaysPrice}
            dayChange={ticker.dayChange}
            weekChange={ticker.weekChange}
            yearChange={ticker.yearChange}
          />
        ))}
      </div>

      {totalPages > 1 && (
        <div className="flex justify-center pt-2">
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious 
                  onClick={() => currentPage > 1 && setCurrentPage(p => Math.max(1, p - 1))}
                  className={cn(
                    "text-terminal-text hover:text-terminal-accent cursor-pointer",
                    currentPage === 1 && "pointer-events-none opacity-50"
                  )}
                />
              </PaginationItem>
              {Array.from({ length: totalPages }).map((_, i) => (
                <PaginationItem key={i}>
                  <PaginationLink
                    onClick={() => setCurrentPage(i + 1)}
                    isActive={currentPage === i + 1}
                    className={cn(
                      "text-terminal-text hover:text-terminal-accent cursor-pointer",
                      currentPage === i + 1 && "border border-terminal-accent bg-transparent"
                    )}
                  >
                    {i + 1}
                  </PaginationLink>
                </PaginationItem>
              ))}
              <PaginationItem>
                <PaginationNext
                  onClick={() => currentPage < totalPages && setCurrentPage(p => Math.min(totalPages, p + 1))}
                  className={cn(
                    "text-terminal-text hover:text-terminal-accent cursor-pointer",
                    currentPage === totalPages && "pointer-events-none opacity-50"
                  )}
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </div>
      )}
    </div>
  );
};
