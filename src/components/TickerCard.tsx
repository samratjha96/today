import { cn } from "@/lib/utils";
interface TickerCardProps {
  ticker: string;
  todaysPrice: number;
  dayChange: number;
  weekChange: number;
  yearChange: number;
}
export const TickerCard = ({
  ticker,
  todaysPrice,
  dayChange,
  weekChange,
  yearChange,
}: TickerCardProps) => {
  return (
    <div className="border border-terminal-accent/20 rounded-lg p-4 bg-terminal-bg/50 hover:bg-terminal-secondary/50 transition-colors">
      <div className="flex items-center justify-between mb-2">
        <span className="text-terminal-text font-mono text-lg">{ticker}</span>
        <span className="text-terminal-text font-mono text-lg">
          ${todaysPrice.toFixed(2)}
        </span>
      </div>
      <div className="grid grid-cols-3 gap-2 text-xs font-mono">
        <div>
          <span className="text-terminal-text/60 block">24h</span>
          <span
            className={cn(dayChange >= 0 ? "text-green-400" : "text-red-400")}
          >
            {dayChange.toFixed(2)}%
          </span>
        </div>
        <div>
          <span className="text-terminal-text/60 block">5d</span>
          <span
            className={cn(weekChange >= 0 ? "text-green-400" : "text-red-400")}
          >
            {weekChange.toFixed(2)}%
          </span>
        </div>
        <div>
          <span className="text-terminal-text/60 block">1y</span>
          <span
            className={cn(yearChange >= 0 ? "text-green-400" : "text-red-400")}
          >
            {yearChange.toFixed(2)}%
          </span>
        </div>
      </div>
    </div>
  );
};
