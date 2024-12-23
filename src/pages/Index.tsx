import { TickerTable } from "@/components/TickerTable";
import { RSSFeed } from "@/components/RSSFeed";
import { useTickerData, getMarketSentiment } from "@/hooks/use-ticker-data";
import { cn } from "@/lib/utils";

const Index = () => {
  const { data } = useTickerData();
  const marketSentiment = data ? getMarketSentiment(data) : null;

  return (
    <div className="min-h-screen bg-terminal-bg text-terminal-text p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="border-b border-terminal-accent/20 pb-4">
          <div className="flex items-center space-x-2">
            <div className="w-3 h-3 rounded-full bg-terminal-accent animate-pulse" />
            <h1 className="font-mono text-xl">
              Today is looking like a {" "}
              {marketSentiment && (
                <span className={cn(
                  marketSentiment === "bullish" ? "text-green-400" : "text-red-400"
                )}>
                  {marketSentiment === "bullish" ? "W" : "L"}
                </span>
              )}
            </h1>
          </div>
        </div>
        
        <div className="space-y-6">
          <section>
            <div className="mb-4">
              <h2 className="font-mono text-sm text-terminal-text/60 flex items-center gap-2">
                📈 Ticker Prices
              </h2>
            </div>
            <TickerTable />
          </section>

          <section>
            <div className="mb-4">
              <h2 className="font-mono text-sm text-terminal-text/60 flex items-center gap-2">
                💻 💰 Top Tech and Finance News
              </h2>
            </div>
            <RSSFeed />
          </section>
        </div>
      </div>
    </div>
  );
};

export default Index;
