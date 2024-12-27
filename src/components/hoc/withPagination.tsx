import { cn } from "@/lib/utils";
import { ReactNode, useState } from "react";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "../ui/pagination";

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
}

export function withPagination<T>(
  WrappedComponent: React.ComponentType<WithPaginationProps<T>>
) {
  return function PaginatedComponent(props: WithPaginationProps<T>) {
    const [internalCurrentPage, setInternalCurrentPage] = useState(1);
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
    } = props;

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

    return (
      <div className={containerClassName}>
        <div className={itemsContainerClassName}>
          {currentItems.map((item, index) => renderItem(item))}
        </div>
        {totalPages > 1 && (
          <PaginationControls
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={onPageChange}
            variant={paginationVariant}
          />
        )}
      </div>
    );
  };
}
