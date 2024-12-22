import { TickerTable } from "@/components/TickerTable";
import { RSSFeed } from "@/components/RSSFeed";

const Index = () => {
  return (
    <div className="min-h-screen bg-terminal-bg text-terminal-text p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="border-b border-terminal-accent/20 pb-4">
          <div className="flex items-center space-x-2">
            <div className="w-3 h-3 rounded-full bg-terminal-accent animate-pulse" />
            <h1 className="font-mono text-xl">Today's Dashboard</h1>
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
            <RSSFeed />
          </section>
        </div>
      </div>
    </div>
  );
};

export default Index;