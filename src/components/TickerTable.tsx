import { cn } from "@/lib/utils";
import { useState } from "react";
import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from "./ui/pagination";
import { useTickerData } from "@/hooks/use-ticker-data";

const ITEMS_PER_PAGE = 5;

export const TickerTable = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const { data, isLoading, error } = useTickerData();

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
    <div className="overflow-x-auto animate-fadeIn">
      <div className="inline-block min-w-full align-middle">
        <div className="border border-terminal-accent/20 rounded-lg overflow-hidden">
          <table className="min-w-full divide-y divide-terminal-accent/20">
            <thead className="bg-terminal-secondary">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-mono text-terminal-text">Ticker</th>
                <th className="px-6 py-3 text-right text-xs font-mono text-terminal-text">Today's Price</th>
                <th className="px-6 py-3 text-right text-xs font-mono text-terminal-text">24h Change (%)</th>
                <th className="px-6 py-3 text-right text-xs font-mono text-terminal-text">5d Change (%)</th>
                <th className="px-6 py-3 text-right text-xs font-mono text-terminal-text">1y Change (%)</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-terminal-accent/20 bg-terminal-bg/50">
              {paginatedData?.map((ticker) => (
                <tr key={ticker.ticker} className="hover:bg-terminal-secondary/50 transition-colors">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-terminal-text">
                    {ticker.ticker}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-right text-terminal-text">
                    {ticker.todaysPrice.toFixed(2)}
                  </td>
                  <td className={cn(
                    "px-6 py-4 whitespace-nowrap text-sm font-mono text-right",
                    ticker.dayChange >= 0 ? "text-green-400" : "text-red-400"
                  )}>
                    {ticker.dayChange.toFixed(2)}%
                  </td>
                  <td className={cn(
                    "px-6 py-4 whitespace-nowrap text-sm font-mono text-right",
                    ticker.weekChange >= 0 ? "text-green-400" : "text-red-400"
                  )}>
                    {ticker.weekChange.toFixed(2)}%
                  </td>
                  <td className={cn(
                    "px-6 py-4 whitespace-nowrap text-sm font-mono text-right",
                    ticker.yearChange >= 0 ? "text-green-400" : "text-red-400"
                  )}>
                    {ticker.yearChange.toFixed(2)}%
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {totalPages > 1 && (
            <div className="bg-terminal-secondary px-6 py-3">
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
      </div>
    </div>
  );
};
