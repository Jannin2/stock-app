package database

import "github.com/jannin2/stock-app/backend/models"

// StockDB define las operaciones que cualquier base de datos de stocks debe implementar.
// Esto permite que el código que interactúa con la base de datos sea independiente de la implementación específica.
type StockDB interface {
	GetAllStocks(opts StockQueryOptions) ([]models.Stock, error)
	GetStockByID(id string) (models.Stock, error)
	UpsertStocks(stocks []models.Stock) error
	GetStockCount(searchQuery string) (int, error)
	GetRecommendedStocks(limit int) ([]models.Stock, error)
}

// StockQueryOptions define los parámetros para consultar stocks.
// Incluye opciones de búsqueda, ordenamiento y paginación.
type StockQueryOptions struct {
	Search string // Término de búsqueda para filtrar por ticker o compañía
	SortBy string // Columna por la cual ordenar (ej. "ticker", "current_price")
	Order  string // Orden del sort: "asc" (ascendente) o "desc" (descendente)
	Limit  int    // Número máximo de resultados a devolver
	Offset int    // Número de resultados a omitir (para paginación)
}
