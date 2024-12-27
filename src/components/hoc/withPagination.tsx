import { cn } from "@/lib/utils";
import { ReactNode, useState, useEffect } from "react";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  variant?: "simple" | "numbered";
}

const PaginationControls = ({
  currentPage,
  totalPages,
  onPageChange,
  variant = "simple",
}: PaginationProps) => {
  if (variant === "numbered") {
    return (
      <div className="flex justify-center pt-2">
        <Pagination>
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                onClick={() =>
                  currentPage > 1 && onPageChange(Math.max(1, currentPage - 1))
                }
                className={cn(
                  "text-terminal-text hover:text-terminal-accent cursor-pointer",
                  currentPage === 1 && "pointer-events-none opacity-50"
                )}
              />
            </PaginationItem>
            {Array.from({ length: totalPages }).map((_, i) => (
              <PaginationItem key={i}>
                <PaginationLink
                  onClick={() => onPageChange(i + 1)}
                  isActive={currentPage === i + 1}
                  className={cn(
                    "text-terminal-text hover:text-terminal-accent cursor-pointer",
                    currentPage === i + 1 &&
                    "border border-terminal-accent bg-transparent"
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
                  onPageChange(Math.min(totalPages, currentPage + 1))
                }
                className={cn(
                  "text-terminal-text hover:text-terminal-accent cursor-pointer",
                  currentPage === totalPages && "pointer-events-none opacity-50"
                )}
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-between p-4 border-t border-terminal-accent/20 bg-terminal-bg/50">
      <button
        onClick={() => onPageChange(Math.max(1, currentPage - 1))}
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
        onClick={() => onPageChange(Math.min(totalPages, currentPage + 1))}
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
  );
};

export interface WithPaginationProps<T> {
  data: T[] | undefined;
  isLoading: boolean;
  error: unknown;
  itemsPerPage?: number;
  renderItem: (item: T) => ReactNode;
  loadingComponent?: ReactNode;
  errorComponent?: ReactNode;
  containerClassName?: string;
  itemsContainerClassName?: string;
  paginationVariant?: "simple" | "numbered";
  currentPage?: number;
  onPageChange?: (page: number) => void;
  autoPaginate?: boolean;
}

export function withPagination<T>(
  WrappedComponent: React.ComponentType<WithPaginationProps<T>>
) {
  return function PaginatedComponent(props: WithPaginationProps<T>) {
    const [internalCurrentPage, setInternalCurrentPage] = useState(1);
    const [direction, setDirection] = useState<'up' | 'down' | null>(null);
    const [isAnimating, setIsAnimating] = useState(false);
    const {
      data,
      isLoading,
      error,
      itemsPerPage = 5,
      renderItem,
      loadingComponent,
      errorComponent,
      containerClassName = "border border-terminal-accent/20 rounded-lg overflow-hidden",
      itemsContainerClassName = "divide-y divide-terminal-accent/20 max-h-[600px] overflow-y-auto",
      paginationVariant = "simple",
      currentPage = internalCurrentPage,
      onPageChange = setInternalCurrentPage,
      autoPaginate = true,
    } = props;

    useEffect(() => {
      if (isAnimating) {
        const timer = setTimeout(() => {
          setIsAnimating(false);
          setDirection(null);
        }, 500);
        return () => clearTimeout(timer);
      }
    }, [isAnimating]);

    useEffect(() => {
      let interval: NodeJS.Timeout | null = null;

      if (autoPaginate && data && data.length > itemsPerPage) {
        interval = setInterval(() => {
          const totalPages = Math.ceil(data.length / itemsPerPage);
          const nextPage = currentPage < totalPages ? currentPage + 1 : 1;
          handlePageChange(nextPage);
        }, 15000);
      }

      return () => {
        if (interval) {
          clearInterval(interval);
        }
      };
    }, [autoPaginate, currentPage, data, itemsPerPage]);

    const handlePageChange = (newPage: number) => {
      setDirection(newPage > currentPage ? 'up' : 'down');
      setIsAnimating(true);
      onPageChange(newPage);
    };

    if (isLoading) {
      return (
        loadingComponent || (
          <div className="animate-pulse p-4 bg-terminal-secondary rounded-lg">
            <div className="h-40 bg-terminal-muted rounded" />
          </div>
        )
      );
    }

    if (error || !data || data.length === 0) {
      return (
        errorComponent || (
          <div className="p-8 text-center border border-terminal-accent/20 rounded-lg bg-terminal-bg/50">
            <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-terminal-secondary mb-4">
              <span className="text-xl">❌</span>
            </div>
            <h3 className="text-terminal-text font-mono text-lg mb-2">
              No Data Available
            </h3>
            <p className="text-terminal-text/60 font-mono text-sm">
              Unable to fetch data. Please try again later.
            </p>
          </div>
        )
      );
    }

    const totalPages = Math.ceil(data.length / itemsPerPage);
    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    const currentItems = data.slice(startIndex, endIndex);
    const shouldAnimate = totalPages > 1;

    return (
      <div className={containerClassName}>
        <style>
          {`
            @keyframes scrollUp {
              0% {
                transform: translateY(100%);
              }
              100% {
                transform: translateY(0);
              }
            }
            @keyframes scrollDown {
              0% {
                transform: translateY(-100%);
              }
              100% {
                transform: translateY(0);
              }
            }
            .scroll-up {
              animation: scrollUp 0.5s cubic-bezier(0.4, 0.0, 0.2, 1) forwards;
            }
            .scroll-down {
              animation: scrollDown 0.5s cubic-bezier(0.4, 0.0, 0.2, 1) forwards;
            }
            .hide-scrollbar {
              -ms-overflow-style: none;
              scrollbar-width: none;
            }
            .hide-scrollbar::-webkit-scrollbar {
              display: none;
            }
            .content-container {
              position: relative;
              overflow: hidden;
            }
          `}
        </style>
        <div className={cn(
          itemsContainerClassName,
          shouldAnimate && isAnimating && "hide-scrollbar overflow-hidden"
        )}>
          <div className={shouldAnimate ? "content-container" : undefined}>
            <div className={cn(
              "divide-y divide-terminal-accent/20",
              shouldAnimate && direction === 'up' && isAnimating && 'scroll-up',
              shouldAnimate && direction === 'down' && isAnimating && 'scroll-down'
            )}>
              {currentItems.map((item, index) => renderItem(item))}
            </div>
          </div>
        </div>
        {totalPages > 1 && (
          <PaginationControls
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={handlePageChange}
            variant={paginationVariant}
          />
        )}
      </div>
    );
  };
}
