package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jannin2/stock-app/backend/models"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// StockDB interface defines the methods for stock-related database operations.

// DB is kept for compatibility if other parts of your app still rely on a global DB.
// However, the `cockroachDB` methods below will use their encapsulated `c.db`.
var DB *sql.DB

// cockroachDB implements the StockDB interface.
type cockroachDB struct {
	db *sql.DB // The actual database connection encapsulated within the struct
}

// NewStockDB creates a new instance of StockDB.
// It returns a pointer to cockroachDB, which implements the StockDB interface.
func NewStockDB(dbConn *sql.DB) StockDB {
	return &cockroachDB{db: dbConn}
}

// ConnectDB establishes a connection to the PostgreSQL database.
// This function still sets the global DB for general use, but can be adapted
// to return a StockDB interface instead if preferred for the main application.
func ConnectDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Println("DATABASE_URL no está configurada, usando valor por defecto.")
	}

	db, err := sql.Open("postgres", connStr) // Use a local variable
	if err != nil {
		return nil, fmt.Errorf("error al abrir la conexión a la base de datos: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		db.Close() // Close on ping failure
		return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
	}

	log.Println("Conexión a la base de datos establecida correctamente.")
	return db, nil // <--- MODIFIED: Return the db connection
}

// CloseDB closes the database connection.
func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
		log.Println("Conexión a la base de datos cerrada.")
	}
}

// InitSchema inicializa el esquema de la base de datos.
// This function still uses the global DB.
func InitSchema(dbConn *sql.DB) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS stocks (
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

	_, err := dbConn.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error al crear/verificar la tabla 'stocks': %w", err)
	}

	addUniqueConstraintSQL := `
    ALTER TABLE stocks ADD CONSTRAINT IF NOT EXISTS stocks_ticker_key UNIQUE (ticker);`

	_, err = dbConn.Exec(addUniqueConstraintSQL)
	if err != nil {
		log.Printf("Advertencia/Error al añadir o verificar la restricción UNIQUE a 'ticker': %v", err)
	}

	alterTableSQLs := []string{
		`ALTER TABLE stocks ADD COLUMN IF NOT EXISTS pe_ratio DECIMAL(10, 2);`,
		`ALTER TABLE stocks ADD COLUMN IF NOT EXISTS dividend_yield DECIMAL(10, 4);`,
		`ALTER TABLE stocks ADD COLUMN IF NOT EXISTS market_capitalization DECIMAL(20, 2);`,
		`ALTER TABLE stocks ADD COLUMN IF NOT EXISTS alpha DECIMAL(10, 4);`,
		`ALTER TABLE stocks ADD COLUMN IF NOT EXISTS recommendation_score DECIMAL(5, 2);`,
	}

	for _, sql := range alterTableSQLs {
		_, err := dbConn.Exec(sql)
		if err != nil {
			log.Printf("Advertencia: No se pudo añadir/alterar columna con SQL: %s, Error: %v", sql, err)
		}
	}

	log.Println("Esquema de la base de datos inicializado (tabla 'stocks' y columnas verificadas/creadas).")
	return nil
}

// --- Métodos de *cockroachDB que implementan la interfaz StockDB ---

// GetStockCount returns the total count of stocks, optionally filtered by a search query.
func (c *cockroachDB) GetStockCount(searchQuery string) (int, error) {
	query := "SELECT COUNT(*) FROM stocks"
	args := []interface{}{}
	if searchQuery != "" {
		query += " WHERE ticker ILIKE $1 OR company ILIKE $2"         // FIX: Changed $1 to $2 for company
		args = append(args, "%"+searchQuery+"%", "%"+searchQuery+"%") // FIX: Added the second argument
	}

	var count int
	err := c.db.QueryRowContext(context.Background(), query, args...).Scan(&count) // Use c.db and context
	if err != nil {
		return 0, fmt.Errorf("error al obtener el recuento de stocks: %w", err)
	}
	return count, nil
}

// GetAllStocks fetches all stocks from the database with pagination, search, and sorting.
func (c *cockroachDB) GetAllStocks(opts StockQueryOptions) ([]models.Stock, error) {
	// First, get the total count for pagination metadata (if needed by your API response)
	// This call will execute the COUNT(*) query
	_, err := c.GetStockCount(opts.Search) // Execute GetStockCount here
	if err != nil {
		log.Printf("Advertencia: No se pudo obtener el recuento de stocks: %v", err)
		// Decide if this should be a fatal error or just logged.
		// For now, it's just a warning, but if count is essential for your API, return error.
	}

	query := "SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, current_price, pe_ratio, dividend_yield, market_capitalization, alpha, latest_trading_day, recommendation_score, created_at, updated_at FROM stocks"
	args := []interface{}{}
	argCounter := 1 // Start counter for positional arguments

	// Add search filter
	if opts.Search != "" {
		query += fmt.Sprintf(" WHERE ticker ILIKE $%d OR company ILIKE $%d", argCounter, argCounter+1)
		args = append(args, "%"+opts.Search+"%", "%"+opts.Search+"%")
		argCounter += 2
	}

	// Add sorting
	if opts.SortBy != "" {
		validSortColumns := map[string]bool{
			"ticker": true, "company": true, "current_price": true,
			"action": true, "recommendation_score": true, "pe_ratio": true,
			"dividend_yield": true, "market_capitalization": true, "alpha": true,
		}
		if !validSortColumns[opts.SortBy] {
			opts.SortBy = "ticker" // Default to a safe column
		}

		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", opts.SortBy, order)
	} else {
		query += " ORDER BY ticker ASC" // Default sort
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, opts.Limit, opts.Offset)

	rows, err := c.db.QueryContext(context.Background(), query, args...) // Use c.db and context
	if err != nil {
		return nil, fmt.Errorf("error al consultar todos los stocks: %w", err)
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var s models.Stock
		var latestTradingDay sql.NullTime
		var targetFrom, targetTo, peRatio, dividendYield, marketCap, alpha, recScore sql.NullFloat64 // Define here for scanning
		err := rows.Scan(
			&s.ID, &s.Ticker, &s.Company, &s.Brokerage, &s.Action,
			&s.RatingFrom, &s.RatingTo, &targetFrom, &targetTo, &s.CurrentPrice,
			&peRatio, &dividendYield, &marketCap, &alpha,
			&latestTradingDay, &recScore,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear fila de stock: %w", err)
		}
		// Assign to models.Null* types
		s.TargetFrom = models.NullFloat64{NullFloat64: targetFrom}
		s.TargetTo = models.NullFloat64{NullFloat64: targetTo}
		s.PERatio = models.NullFloat64{NullFloat64: peRatio}
		s.DividendYield = models.NullFloat64{NullFloat64: dividendYield}
		s.MarketCapitalization = models.NullFloat64{NullFloat64: marketCap}
		s.Alpha = models.NullFloat64{NullFloat64: alpha}
		s.LatestTradingDay = models.NullTime{NullTime: latestTradingDay}
		s.RecommendationScore = models.NullFloat64{NullFloat64: recScore}

		stocks = append(stocks, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error después de iterar filas: %w", err)
	}

	return stocks, nil
}

// GetStockByID fetches a single stock by its ID.
func (c *cockroachDB) GetStockByID(id string) (models.Stock, error) {
	query := `SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, current_price, pe_ratio, dividend_yield, market_capitalization, alpha, latest_trading_day, recommendation_score, created_at, updated_at FROM stocks WHERE id = $1`
	var s models.Stock
	var latestTradingDay sql.NullTime
	var targetFrom, targetTo, peRatio, dividendYield, marketCap, alpha, recScore sql.NullFloat64

	err := c.db.QueryRowContext(context.Background(), query, id).Scan( // Use c.db and context
		&s.ID, &s.Ticker, &s.Company, &s.Brokerage, &s.Action,
		&s.RatingFrom, &s.RatingTo, &targetFrom, &targetTo, &s.CurrentPrice,
		&peRatio, &dividendYield, &marketCap, &alpha,
		&latestTradingDay, &recScore,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Stock{}, fmt.Errorf("stock con ID %s no encontrado", id)
		}
		return models.Stock{}, fmt.Errorf("error al obtener stock por ID %s: %w", id, err)
	}
	s.TargetFrom = models.NullFloat64{NullFloat64: targetFrom}
	s.TargetTo = models.NullFloat64{NullFloat64: targetTo}
	s.PERatio = models.NullFloat64{NullFloat64: peRatio}
	s.DividendYield = models.NullFloat64{NullFloat64: dividendYield}
	s.MarketCapitalization = models.NullFloat64{NullFloat64: marketCap}
	s.Alpha = models.NullFloat64{NullFloat64: alpha}
	s.LatestTradingDay = models.NullTime{NullTime: latestTradingDay}
	s.RecommendationScore = models.NullFloat64{NullFloat64: recScore}

	return s, nil
}

// GetRecommendedStocks fetches a limited number of stocks ordered by recommendation_score.
func (c *cockroachDB) GetRecommendedStocks(limit int) ([]models.Stock, error) {
	query := `SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, current_price, pe_ratio, dividend_yield, market_capitalization, alpha, latest_trading_day, recommendation_score, created_at, updated_at FROM stocks ORDER BY recommendation_score DESC NULLS LAST LIMIT $1`

	rows, err := c.db.QueryContext(context.Background(), query, limit) // Use c.db and context
	if err != nil {
		return nil, fmt.Errorf("error al consultar stocks recomendados: %w", err)
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var s models.Stock
		var latestTradingDay sql.NullTime
		var targetFrom, targetTo, peRatio, dividendYield, marketCap, alpha, recScore sql.NullFloat64
		err := rows.Scan(
			&s.ID, &s.Ticker, &s.Company, &s.Brokerage, &s.Action,
			&s.RatingFrom, &s.RatingTo, &targetFrom, &targetTo, &s.CurrentPrice,
			&peRatio, &dividendYield, &marketCap, &alpha,
			&latestTradingDay, &recScore,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear fila de stock recomendado: %w", err)
		}
		s.TargetFrom = models.NullFloat64{NullFloat64: targetFrom}
		s.TargetTo = models.NullFloat64{NullFloat64: targetTo}
		s.PERatio = models.NullFloat64{NullFloat64: peRatio}
		s.DividendYield = models.NullFloat64{NullFloat64: dividendYield}
		s.MarketCapitalization = models.NullFloat64{NullFloat64: marketCap}
		s.Alpha = models.NullFloat64{NullFloat64: alpha}
		s.LatestTradingDay = models.NullTime{NullTime: latestTradingDay}
		s.RecommendationScore = models.NullFloat64{NullFloat64: recScore}

		stocks = append(stocks, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error después de iterar filas recomendadas: %w", err)
	}

	return stocks, nil
}

// UpsertStocks inserts new stocks or updates existing ones based on their ticker.
func (c *cockroachDB) UpsertStocks(stocks []models.Stock) error {
	if len(stocks) == 0 {
		return nil // Nothing to upsert
	}

	tx, err := c.db.BeginTx(context.Background(), nil) // Use c.db and context for transaction
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción para upsert: %w", err)
	}
	defer tx.Rollback() // Rollback on error or if commit fails

	stmt, err := tx.PrepareContext(context.Background(), ` // Use context for prepare
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
    `)
	if err != nil {
		return fmt.Errorf("error al preparar la declaración upsert: %w", err)
	}
	defer stmt.Close()

	for _, s := range stocks {
		_, err := stmt.ExecContext(context.Background(), // Use context for exec
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
		)
		if err != nil {
			log.Printf("ERROR UPSERT para ticker %s: %v. Valores de depuración: TargetFrom.Float64=%.2f (Valid:%t), TargetTo.Float64=%.2f (Valid:%t), CurrentPrice=%.2f, PERatio.Float64=%.2f (Valid:%t), DividendYield.Float64=%.4f (Valid:%t), MarketCapitalization.Float64=%.2f (Valid:%t), Alpha.Float64=%.4f (Valid:%t), LatestTradingDay.Time=%v (Valid:%t), RecommendationScore.Float64=%.2f (Valid:%t)",
				s.Ticker, err,
				s.TargetFrom.Float64, s.TargetFrom.Valid,
				s.TargetTo.Float64, s.TargetTo.Valid,
				s.CurrentPrice,
				s.PERatio.Float64, s.PERatio.Valid,
				s.DividendYield.Float64, s.DividendYield.Valid,
				s.MarketCapitalization.Float64, s.MarketCapitalization.Valid,
				s.Alpha.Float64, s.Alpha.Valid,
				s.LatestTradingDay.Time, s.LatestTradingDay.Valid,
				s.RecommendationScore.Float64, s.RecommendationScore.Valid)
			return fmt.Errorf("error al ejecutar upsert para el ticker %s: %w", s.Ticker, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar la transacción upsert: %w", err)
	}

	return nil
}
