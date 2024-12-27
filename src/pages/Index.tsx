import { TickerTable } from "@/components/TickerTable";
import { RSSFeed } from "@/components/RSSFeed";
import { useTickerData, getMarketSentiment } from "@/hooks/use-ticker-data";
import { GithubTrending } from "@/components/GithubTrending";
import { HackerNews } from "@/components/HackerNews";
import { cn } from "@/lib/utils";

const Index = () => {
  const { data } = useTickerData();
  const marketSentiment = data ? getMarketSentiment(data) : null;

  return (
    <div className="min-h-screen bg-terminal-bg text-terminal-text p-4 md:p-6">
      <div className="max-w-7xl mx-auto">
        <div className="border-b border-terminal-accent/20 pb-4 mb-6">
          <div className="flex items-center justify-center space-x-2">
            <div className="w-3 h-3 rounded-full bg-terminal-accent animate-pulse" />
            <h1 className="font-mono text-xl">
              Today is looking like a{" "}
              <span className="inline-flex items-center">
                {marketSentiment ? (
                  <span
                    className={cn(
                      marketSentiment === "bullish" ? "text-green-400" : "text-red-400"
                    )}
                  >
                    {marketSentiment === "bullish" ? "W" : "L"}
                  </span>
                ) : null}
                <span className="text-green-400 inline-block w-[2px] h-7 ml-[1px] align-middle animate-subtleBlink">|</span>
              </span>
            </h1>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-12 gap-6">
          {/* Left column - Ticker Cards */}
          <div className="md:col-span-3 space-y-4">
            <section>
              <div className="mb-4">
                <h2 className="font-mono text-sm text-terminal-text/60 flex items-center gap-2 mb-4">
                  📈 Market Data
                </h2>
              </div>
              <TickerTable />
            </section>
          </div>

          {/* Right column - News and Updates */}
          <div className="md:col-span-9 space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <section>
                <h2 className="font-mono text-sm text-terminal-text/60 flex items-center gap-2 mb-4">
                  🐙 Github Trending
                </h2>
                <GithubTrending />
              </section>
              <section>
                <h2 className="font-mono text-sm text-terminal-text/60 flex items-center gap-2 mb-4">
                  🔥 Hacker News
                </h2>
                <HackerNews />
              </section>
            </div>
            <section>
              <h2 className="font-mono text-sm text-terminal-text/60 flex items-center gap-2 mb-4">
                💻 Tech News
              </h2>
              <RSSFeed />
            </section>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Index;
