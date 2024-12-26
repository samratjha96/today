import { useState } from "react";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "./ui/pagination";
import { cn } from "@/lib/utils";
import { useRSSData } from "@/hooks/use-rss-data";

const ITEMS_PER_PAGE = 5;

export const RSSFeed = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const { data, isLoading, error } = useRSSData();

  if (isLoading) {
    return (
      <div className="animate-pulse p-4 bg-terminal-secondary rounded-lg">
        <div className="h-60 bg-terminal-muted rounded" />
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
          No News Available
        </h3>
        <p className="text-terminal-text/60 font-mono text-sm">
          Unable to fetch news articles at this time. Please check back later.
        </p>
      </div>
    );
  }

  const totalPages = Math.ceil((data?.length || 0) / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const paginatedData = data?.slice(startIndex, startIndex + ITEMS_PER_PAGE);

  return (
    <div className="border border-terminal-accent/20 rounded-lg overflow-hidden animate-fadeIn">
      <div className="divide-y divide-terminal-accent/20">
        {paginatedData?.map((item, index) => (
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
                className="text-[#33C3F0] font-mono text-sm hover:text-[#1EAEDB] transition-colors line-clamp-2"
              >
                {item.title}
              </a>
            </div>
          </div>
        ))}
      </div>
      {totalPages > 1 && (
        <div className="bg-terminal-secondary px-6 py-3">
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  onClick={() =>
                    currentPage > 1 && setCurrentPage((p) => Math.max(1, p - 1))
                  }
                  className={cn(
                    "text-[#33C3F0] hover:text-[#1EAEDB] cursor-pointer",
                    currentPage === 1 && "pointer-events-none opacity-50",
                  )}
                />
              </PaginationItem>
              {Array.from({ length: totalPages }).map((_, i) => (
                <PaginationItem key={i}>
                  <PaginationLink
                    onClick={() => setCurrentPage(i + 1)}
                    isActive={currentPage === i + 1}
                    className={cn(
                      "text-[#33C3F0] hover:text-[#1EAEDB] cursor-pointer",
                      currentPage === i + 1 &&
                        "border border-[#33C3F0] bg-transparent",
                    )}
                  >
                    {i + 1}
                  </PaginationLink>
                </PaginationItem>
              ))}
              <PaginationItem>
                <PaginationNext
                  onClick={() =>
                    currentPage < totalPages &&
                    setCurrentPage((p) => Math.min(totalPages, p + 1))
                  }
                  className={cn(
                    "text-[#33C3F0] hover:text-[#1EAEDB] cursor-pointer",
                    currentPage === totalPages &&
                      "pointer-events-none opacity-50",
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
