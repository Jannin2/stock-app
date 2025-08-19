package enricher

import (
	"database/sql"
	"log"
	"time"

	"github.com/jannin2/stock-app/backend/api"
	"github.com/jannin2/stock-app/backend/database"
	"github.com/jannin2/stock-app/backend/models"
)

// Enricher handles fetching and updating stock data periodically.
type Enricher struct {
	dbClient database.StockDB // This is where your database interface is held
}

// NewEnricher creates a new Enricher instance.
// It receives the StockDB interface as a dependency.
func NewEnricher(dbClient database.StockDB) *Enricher {
	return &Enricher{
		dbClient: dbClient,
	}
}

// StartFetching initiates the cron job to fetch and update stock data.
// This is the entry point for the periodic task.
func (e *Enricher) StartFetching() {
	// Execute immediately once at startup
	log.Println("🔄 Starting initial stock data enrichment...")
	e.fetchAndEnrichStocks() // Calls the method that contains all the logic

	// Then, execute on each ticker tick (e.g., every 24 hours)
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("⏰ Executing scheduled stock data enrichment...")
		e.fetchAndEnrichStocks() // Calls the method that contains all the logic
	}
}

// fetchAndEnrichStocks contains the logic to fetch data from external APIs and update it in the DB.
// This method is now part of the Enricher, allowing it to access e.dbClient.
func (e *Enricher) fetchAndEnrichStocks() {
	log.Println("Starting stock data enrichment...")

	stocksFromKarenai, err := api.GetRecommendationsFromKarenai()
	if err != nil {
		log.Printf("Error getting recommendations from Karenai.click: %v", err)
		return
	}
	log.Printf("Received %d recommendations from Karenai.click", len(stocksFromKarenai))

	for i := range stocksFromKarenai {
		ticker := stocksFromKarenai[i].Ticker
		log.Printf("Enriching data for ticker: %s", ticker)

		// --- Get Current Price and Finnhub Metrics ---
		finnhubMetrics, err := api.GetFinnhubMetricsAndQuote(ticker)
		if err != nil {
			log.Printf("Error getting metrics/price from Finnhub for %s: %v. Assigning null/default values.", ticker, err)
			stocksFromKarenai[i].PERatio = models.NullFloat64{sql.NullFloat64{Valid: false}}
			stocksFromKarenai[i].DividendYield = models.NullFloat64{sql.NullFloat64{Valid: false}}
			stocksFromKarenai[i].MarketCapitalization = models.NullFloat64{sql.NullFloat64{Valid: false}}
			stocksFromKarenai[i].CurrentPrice = 0.0
			stocksFromKarenai[i].LatestTradingDay = models.NullTime{NullTime: sql.NullTime{Valid: false}}
		} else {
			stocksFromKarenai[i].PERatio = models.NullFloat64{sql.NullFloat64{Float64: finnhubMetrics.PE_Ratio, Valid: true}}
			stocksFromKarenai[i].DividendYield = models.NullFloat64{sql.NullFloat64{Float64: finnhubMetrics.DividendYield, Valid: true}}
			stocksFromKarenai[i].MarketCapitalization = models.NullFloat64{sql.NullFloat64{Float64: finnhubMetrics.MarketCapitalization, Valid: true}}
			stocksFromKarenai[i].CurrentPrice = finnhubMetrics.CurrentPrice

			if !finnhubMetrics.LatestTradingDay.IsZero() {
				stocksFromKarenai[i].LatestTradingDay = models.NullTime{NullTime: sql.NullTime{Time: finnhubMetrics.LatestTradingDay, Valid: true}}
			} else {
				stocksFromKarenai[i].LatestTradingDay = models.NullTime{NullTime: sql.NullTime{Valid: false}}
			}

			log.Printf("Finnhub data for %s: Price: %.2f, PE: %.2f, Div Yield: %.4f, Market Cap: %.2f, Trading Day (Finnhub): %v",
				ticker, stocksFromKarenai[i].CurrentPrice, finnhubMetrics.PE_Ratio, finnhubMetrics.DividendYield, finnhubMetrics.MarketCapitalization, stocksFromKarenai[i].LatestTradingDay.Time.Format("2006-01-02"))
		}

		// --- Alpha Vantage Alpha ---
		alphaVantageData, err := api.GetAlphaAndLatestTradingDayFromAlphaVantage(ticker)
		if err != nil {
			log.Printf("Error getting Alpha from Alpha Vantage for %s: %v. Assigning null value.", ticker, err)
			stocksFromKarenai[i].Alpha = models.NullFloat64{sql.NullFloat64{Valid: false}}
		} else {
			stocksFromKarenai[i].Alpha = models.NullFloat64{sql.NullFloat64{Float64: alphaVantageData.Alpha, Valid: true}}
			log.Printf("Alpha Vantage data for %s: Alpha: %.4f", ticker, alphaVantageData.Alpha)
		}

		// --- Calculate Recommendation Score ---
		scoreVal := CalculateRecommendationScore(stocksFromKarenai[i])

		stocksFromKarenai[i].RecommendationScore = models.NullFloat64{sql.NullFloat64{Float64: scoreVal, Valid: true}}
		log.Printf("Recommendation score calculated for %s: %.2f", ticker, scoreVal)

		stocksFromKarenai[i].UpdatedAt = time.Now()

		log.Printf("Processed and Enriched %s: Price: %.2f, PE: %.2f (Valid: %t), Div Yield: %.4f (Valid: %t), Market Cap: %.2f (Valid: %t), Alpha: %.4f (Valid: %t), Rec Score: %.2f (Valid: %t), Trading Day: %v (Valid: %t)",
			ticker, stocksFromKarenai[i].CurrentPrice,
			stocksFromKarenai[i].PERatio.Float64, stocksFromKarenai[i].PERatio.Valid,
			stocksFromKarenai[i].DividendYield.Float64, stocksFromKarenai[i].DividendYield.Valid,
			stocksFromKarenai[i].MarketCapitalization.Float64, stocksFromKarenai[i].MarketCapitalization.Valid,
			stocksFromKarenai[i].Alpha.Float64, stocksFromKarenai[i].Alpha.Valid,
			stocksFromKarenai[i].RecommendationScore.Float64, stocksFromKarenai[i].RecommendationScore.Valid,
			func() string {
				if stocksFromKarenai[i].LatestTradingDay.Valid {
					return stocksFromKarenai[i].LatestTradingDay.Time.Format("2006-01-02")
				}
				return "0001-01-01"
			}(), stocksFromKarenai[i].LatestTradingDay.Valid)
	}

	// ✅ THE KEY CORRECTION: Call UpsertStocks via the dbClient instance
	err = e.dbClient.UpsertStocks(stocksFromKarenai)
	if err != nil {
		log.Printf("Error saving/updating stocks in the database: %v", err)
		return
	}
	log.Println("Stock data enriched and saved to the database successfully.")
}

// CalculateRecommendationScore remains an auxiliary function that does not require the DB instance.
func CalculateRecommendationScore(stock models.Stock) float64 {
	scoreVal := 0.0

	// Condition 1: Based on the action
	if stock.Action == "Buy" || stock.Action == "Strong Buy" {
		scoreVal += 5.0
	}

	// Condition 2: Based on target price vs. current price
	// Ensure CurrentPrice is positive to avoid division by zero or nonsensical logic
	if stock.CurrentPrice > 0 && stock.TargetTo.Valid && stock.TargetTo.Float64 > stock.CurrentPrice*1.1 {
		scoreVal += 3.0
	}

	// Alpha contribution removed as per discussion.
	// If you ever integrate a real Alpha, re-add it here.

	return scoreVal
}
