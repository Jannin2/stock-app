package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jannin2/stock-app/backend/database"
)

// StockHandlers contiene la interfaz de la base de datos.
type StockHandlers struct {
	dbClient database.StockDB
}

// NewStockHandlers crea una nueva instancia de StockHandlers.
// Recibe la interfaz StockDB como dependencia.
func NewStockHandlers(dbClient database.StockDB) *StockHandlers {
	return &StockHandlers{dbClient: dbClient}
}

// GetStocks maneja la obtención de una lista de stocks con paginación, búsqueda y ordenamiento.
func (h *StockHandlers) GetStocks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	searchQuery := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	order := r.URL.Query().Get("order")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Límite por defecto
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Offset por defecto
	}

	opts := database.StockQueryOptions{ // Usa database.StockQueryOptions
		Search: searchQuery,
		SortBy: sortBy,
		Order:  strings.ToLower(order),
		Limit:  limit,
		Offset: offset,
	}

	// Llama a los métodos de la interfaz StockDB a través de h.dbClient
	stocks, err := h.dbClient.GetAllStocks(opts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener stocks: %v", err), http.StatusInternalServerError)
		return
	}

	totalCount, err := h.dbClient.GetStockCount(searchQuery)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener el conteo de stocks: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", strconv.Itoa(totalCount))
	json.NewEncoder(w).Encode(stocks)
}

// GetStockByID maneja la obtención de un stock por su ID.
func (h *StockHandlers) GetStockByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Se requiere el ID del stock", http.StatusBadRequest)
		return
	}

	// Llama al método de la interfaz StockDB a través de h.dbClient
	stock, err := h.dbClient.GetStockByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Stock no encontrado: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

// GetRecommendedStocks maneja la obtención de stocks recomendados.
func (h *StockHandlers) GetRecommendedStocks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5 // Límite por defecto para stocks recomendados
	}

	// Llama al método de la interfaz StockDB a través de h.dbClient
	stocks, err := h.dbClient.GetRecommendedStocks(limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener stocks recomendados: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}
