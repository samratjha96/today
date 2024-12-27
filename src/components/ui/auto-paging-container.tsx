import { cn } from "@/lib/utils";
import { ReactNode, useEffect, useState } from "react";

interface AutoPagingContainerProps {
  children: ReactNode[];
  itemsPerPage?: number;
  interval?: number;
  className?: string;
}

export const AutoPagingContainer = ({
  children,
  itemsPerPage = 5,
  interval = 15000,
  className,
}: AutoPagingContainerProps) => {
  const [currentPage, setCurrentPage] = useState(1);

  const totalItems = children.length;
  const totalPages = Math.ceil(totalItems / itemsPerPage);

  // Auto-cycle pages
  useEffect(() => {
    if (totalItems <= itemsPerPage) return;
    
    const timer = setInterval(() => {
      setCurrentPage((current) => (current % totalPages) + 1);
    }, interval);

    return () => clearInterval(timer);
  }, [totalItems, itemsPerPage, totalPages, interval]);

  // Create arrays of items for each page
  const pages = Array.from({ length: totalPages }, (_, i) => {
    const startIndex = i * itemsPerPage;
    return children.slice(startIndex, startIndex + itemsPerPage);
  });

  return (
    <div className={cn("space-y-4", className)}>
      <div className="relative overflow-hidden">
        <div
          className={cn(
            "transition-transform duration-1000 ease-in-out transform",
            "will-change-transform"
          )}
          style={{
            transform: `translateY(${(currentPage - 1) * -100}%)`,
          }}
        >
          {pages.map((pageItems, pageIndex) => (
            <div key={pageIndex} className="relative">
              {pageItems}
            </div>
          ))}
        </div>
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
