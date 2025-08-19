package database

import (
	"database/sql"

	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jannin2/stock-app/backend/models"
)

// Helper function to create sql.NullFloat64 from float64
func newNullFloat64(f float64) models.NullFloat64 {
	return models.NullFloat64{NullFloat64: sql.NullFloat64{Float64: f, Valid: true}}
}

// Helper function to create sql.NullTime from time.Time
func newNullTime(t time.Time) models.NullTime {
	return models.NullTime{NullTime: sql.NullTime{Time: t, Valid: true}}
}

func TestConnectDB(t *testing.T) {
	// This test usually requires a real DB connection or a very advanced mock.

	t.Skip("Skipping ConnectDB test, typically requires real DB or more complex mocking.")
}

func TestInitSchema(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Adjust CREATE TABLE SQL regex to match exactly, removing unnecessary leading/trailing newlines for robustness
	createTableSQL := `CREATE TABLE IF NOT EXISTS stocks (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        ticker VARCHAR(10) NOT NULL UNIQUE,
        company TEXT,
        brokerage TEXT,
        action TEXT,
        rating_from TEXT,
        rating_to TEXT,
        target_from NUMERIC(10, 2)NULL,
        target_to NUMERIC(10, 2)NULL,
        current_price DECIMAL(10, 2),
        pe_ratio DECIMAL(10, 2),
        dividend_yield DECIMAL(10, 4),
        market_capitalization DECIMAL(20, 2),
        alpha DECIMAL(10, 4),
        latest_trading_day TIMESTAMP WITH TIME ZONE,
        recommendation_score DECIMAL(5, 2),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
    );`
	mock.ExpectExec(regexp.QuoteMeta(createTableSQL)).WillReturnResult(sqlmock.NewResult(0, 0))

	// Expect the ALTER TABLE ADD CONSTRAINT statement
	mock.ExpectExec(regexp.QuoteMeta(`ALTER TABLE stocks ADD CONSTRAINT IF NOT EXISTS stocks_ticker_key UNIQUE (ticker);`)).WillReturnResult(sqlmock.NewResult(0, 0))

	// Expect all ALTER TABLE ADD COLUMN statements
	mock.ExpectExec(regexp.QuoteMeta(`ALTER TABLE stocks ADD COLUMN IF NOT EXISTS pe_ratio DECIMAL(10, 2);`)).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta(`ALTER TABLE stocks ADD COLUMN IF NOT EXISTS dividend_yield DECIMAL(10, 4);`)).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta(`ALTER TABLE stocks ADD COLUMN IF NOT EXISTS market_capitalization DECIMAL(20, 2);`)).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta(`ALTER TABLE stocks ADD COLUMN IF NOT EXISTS alpha DECIMAL(10, 4);`)).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta(`ALTER TABLE stocks ADD COLUMN IF NOT EXISTS recommendation_score DECIMAL(5, 2);`)).WillReturnResult(sqlmock.NewResult(0, 0))

	// Call InitSchema with the MOCKED database connection
	err = InitSchema(db)
	if err != nil {
		t.Errorf("❌ error inesperado al inicializar el esquema: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("⚠️ expectativas no cumplidas en TestInitSchema: %s", err)
	}
}

func TestUpsertStocks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sdb := NewStockDB(db)

	mockTime := time.Now()
	testStocks := []models.Stock{
		{
			Ticker:               "TEST1",
			Company:              "Test Company 1",
			Brokerage:            "BrokerX",
			Action:               "Buy",
			RatingFrom:           "Strong Buy",
			RatingTo:             "Buy",
			TargetFrom:           models.NullFloat64{},
			TargetTo:             newNullFloat64(100.50),
			CurrentPrice:         100.00,
			PERatio:              newNullFloat64(20.0),
			DividendYield:        newNullFloat64(0.015),
			MarketCapitalization: newNullFloat64(1.0e9),
			Alpha:                newNullFloat64(0.005),
			LatestTradingDay:     newNullTime(mockTime),
			RecommendationScore:  newNullFloat64(4.0),
		},
		{
			Ticker:               "TEST2",
			Company:              "Test Company 2",
			Brokerage:            "BrokerY",
			Action:               "Hold",
			RatingFrom:           "Buy",
			RatingTo:             "Hold",
			TargetFrom:           newNullFloat64(50.0),
			TargetTo:             models.NullFloat64{},
			CurrentPrice:         50.00,
			PERatio:              newNullFloat64(15.0),
			DividendYield:        newNullFloat64(0.02),
			MarketCapitalization: newNullFloat64(5.0e8),
			Alpha:                newNullFloat64(-0.002),
			LatestTradingDay:     newNullTime(mockTime),
			RecommendationScore:  newNullFloat64(3.0),
		},
	}

	// Expect a transaction begin *before* preparing the statement
	mock.ExpectBegin()

	// The regex for PrepareContext must match the exact string, including comments/newlines
	expectedSQL := ` // Use context for prepare
        INSERT INTO stocks (
            ticker, company, brokerage, action, rating_from, rating_to,
            target_from, target_to, current_price, pe_ratio, dividend_yield,
            market_capitalization, alpha, latest_trading_day, recommendation_score, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, now(), now()
        )
        ON CONFLICT (ticker) DO UPDATE SET
            company = EXCLUDED.company,
            brokerage = EXCLUDED.brokerage,
            action = EXCLUDED.action,
            rating_from = EXCLUDED.rating_from,
            rating_to = EXCLUDED.rating_to,
            target_from = EXCLUDED.target_from,
            target_to = EXCLUDED.target_to,
            current_price = EXCLUDED.current_price,
            pe_ratio = EXCLUDED.pe_ratio,
            dividend_yield = EXCLUDED.dividend_yield,
            market_capitalization = EXCLUDED.market_capitalization,
            alpha = EXCLUDED.alpha,
            latest_trading_day = EXCLUDED.latest_trading_day,
            recommendation_score = EXCLUDED.recommendation_score,
            updated_at = now();
    `
	mock.ExpectPrepare(regexp.QuoteMeta(expectedSQL))

	// Expect each Exec call for the prepared statement
	for _, s := range testStocks {
		mock.ExpectExec(regexp.QuoteMeta(expectedSQL)). // Match the prepared statement regex
								WithArgs(
				s.Ticker, s.Company, s.Brokerage, s.Action, s.RatingFrom, s.RatingTo,
				s.TargetFrom.NullFloat64,
				s.TargetTo.NullFloat64,
				s.CurrentPrice,
				s.PERatio.NullFloat64,
				s.DividendYield.NullFloat64,
				s.MarketCapitalization.NullFloat64,
				s.Alpha.NullFloat64,
				s.LatestTradingDay.NullTime,
				s.RecommendationScore.NullFloat64,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	// Expect a commit
	mock.ExpectCommit()

	err = sdb.UpsertStocks(testStocks)
	if err != nil {
		t.Errorf("❌ error inesperado al upsertar stocks: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("⚠️ expectativas no cumplidas en TestUpsertStocks: %s", err)
	}
}

func TestGetAllStocks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sdb := NewStockDB(db)

	opts := StockQueryOptions{
		Limit:  10,
		Offset: 0,
		SortBy: "ticker",
		Order:  "asc",
		Search: "Test Company", // Set a search term to trigger the WHERE clause in GetStockCount
	}

	// FIX: Expect the COUNT(*) query first, as GetStockCount is called first in GetAllStocks
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM stocks WHERE ticker ILIKE $1 OR company ILIKE $2")). // Updated to $1 and $2
															WithArgs("%"+opts.Search+"%", "%"+opts.Search+"%"). // Two arguments for $1 and $2
															WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	columns := []string{"id", "ticker", "company", "brokerage", "action", "rating_from", "rating_to", "target_from", "target_to", "current_price", "pe_ratio", "dividend_yield", "market_capitalization", "alpha", "latest_trading_day", "recommendation_score", "created_at", "updated_at"}
	mockTime := time.Now()

	rows := sqlmock.NewRows(columns).
		AddRow(uuid.New().String(), "TEST1", "Test Company 1", "BrokerX", "Buy", "Strong Buy", "Buy", nil, 100.50, 100.00, 20.0, 0.015, 1.0e9, 0.005, mockTime, 4.0, mockTime, mockTime)

	// Then expect the main SELECT query for GetAllStocks
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, current_price, pe_ratio, dividend_yield, market_capitalization, alpha, latest_trading_day, recommendation_score, created_at, updated_at FROM stocks WHERE ticker ILIKE $1 OR company ILIKE $2 ORDER BY ticker ASC LIMIT $3 OFFSET $4")).
		WithArgs(
			"%"+opts.Search+"%", "%"+opts.Search+"%", // Args for search (for $1 and $2)
			opts.Limit, opts.Offset, // Args for pagination ($3 and $4)
		).
		WillReturnRows(rows)

	stocks, err := sdb.GetAllStocks(opts)
	if err != nil {
		t.Errorf("❌ error inesperado al obtener stocks: %v", err)
	}

	if len(stocks) != 1 {
		t.Errorf("❌ se esperaban 1 stock, se obtuvieron %d", len(stocks))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("⚠️ expectativas no cumplidas en TestGetAllStocks: %s", err)
	}
}

func TestGetStockByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sdb := NewStockDB(db)

	testID := uuid.New()

	mockTime := time.Now()

	columns := []string{"id", "ticker", "company", "brokerage", "action", "rating_from", "rating_to", "target_from", "target_to", "current_price", "pe_ratio", "dividend_yield", "market_capitalization", "alpha", "latest_trading_day", "recommendation_score", "created_at", "updated_at"}
	rows := sqlmock.NewRows(columns).
		AddRow(testID.String(), "MSFT", "Microsoft", "BrokerC", "Buy", "Hold", "Buy", nil, nil, 405.10, 32.0, 0.007, 3.2e12, 0.008, mockTime, 4.7, mockTime, mockTime)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, current_price, pe_ratio, dividend_yield, market_capitalization, alpha, latest_trading_day, recommendation_score, created_at, updated_at FROM stocks WHERE id = $1`)).
		WithArgs(testID.String()).
		WillReturnRows(rows)

	stock, err := sdb.GetStockByID(testID.String())
	if err != nil {
		t.Errorf("❌ error inesperado al obtener stock por ID: %v", err)
	}

	if stock.ID.String() != testID.String() {
		t.Errorf("❌ ID de stock inesperado: se esperaba %s, se obtuvo %s", testID.String(), stock.ID.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("⚠️ expectativas no cumplidas en TestGetStockByID: %s", err)
	}
}

func TestGetRecommendedStocks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sdb := NewStockDB(db)
	limit := 2
	mockTime := time.Now()

	columns := []string{"id", "ticker", "company", "brokerage", "action", "rating_from", "rating_to", "target_from", "target_to", "current_price", "pe_ratio", "dividend_yield", "market_capitalization", "alpha", "latest_trading_day", "recommendation_score", "created_at", "updated_at"}
	rows := sqlmock.NewRows(columns).
		AddRow(uuid.New().String(), "MSFT", "Microsoft", "BrokerC", "Buy", "Hold", "Buy", nil, nil, 405.10, 32.0, 0.007, 3.2e12, 0.008, mockTime, 4.7, mockTime, mockTime).
		AddRow(uuid.New().String(), "AAPL", "Apple", "BrokerA", "Buy", "Neutral", "Buy", nil, nil, 195.50, 28.5, 0.005, 3.0e12, 0.01, mockTime, 4.5, mockTime, mockTime)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, current_price, pe_ratio, dividend_yield, market_capitalization, alpha, latest_trading_day, recommendation_score, created_at, updated_at FROM stocks ORDER BY recommendation_score DESC NULLS LAST LIMIT $1`)).
		WithArgs(limit).
		WillReturnRows(rows)

	stocks, err := sdb.GetRecommendedStocks(limit)
	if err != nil {
		t.Errorf("❌ error inesperado al obtener stocks recomendados: %v", err)
	}

	if len(stocks) != limit {
		t.Errorf("❌ se esperaban %d stocks, se obtuvieron %d", limit, len(stocks))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("⚠️ expectativas no cumplidas en TestGetRecommendedStocks: %s", err)
	}
}
