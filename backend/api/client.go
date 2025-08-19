package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jannin2/stock-app/backend/handlers"
	"github.com/jannin2/stock-app/backend/models"
)

const (
	KARENAI_API_URL        = "https://api.karenai.click/swechallenge/list"
	FINNHUB_BASE_URL       = "https://finnhub.io/api/v1"
	ALPHA_VANTAGE_BASE_URL = "https://www.alphavantage.co/query"
)

func SetupRouter(r *chi.Mux, stockHandlers *handlers.StockHandlers) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/stocks", func(r chi.Router) {
			r.Get("/", stockHandlers.GetStocks)
			r.Get("/{id}", stockHandlers.GetStockByID)
			r.Get("/recommended", stockHandlers.GetRecommendedStocks)

		})
	})
}

type karenaiResponse struct {
	Items    []models.Stock `json:"items"`
	NextPage string         `json:"next_page"`
}

// Structs for Finnhub responses
type FinnhubMetricResponse struct {
	Metric struct {
		PeExclExtraTTM   float64 `json:"peExclExtraTTM"`
		PeRatio          float64 `json:"peRatio"`
		DividendYield    float64 `json:"dividendYieldAnnually"`
		DividendYieldAlt float64 `json:"dividendYield"`
		MarketCap        float64 `json:"marketCapitalization"`
	} `json:"metric"`
}

type FinnhubQuoteResponse struct {
	CurrentPrice float64 `json:"c"`
	Timestamp    int64   `json:"t"`
}

// Consolidated struct for Finnhub data
type FinnhubData struct {
	PE_Ratio             float64
	DividendYield        float64
	MarketCapitalization float64
	CurrentPrice         float64
	LatestTradingDay     time.Time
	Error                error
}

// Consolidated struct for Alpha Vantage data
type AlphaVantageData struct {
	Alpha            float64
	LatestTradingDay time.Time
	Error            error
}

func GetRecommendationsFromKarenai() ([]models.Stock, error) {
	karenaiAPIKey := os.Getenv("KARENAI_API_KEY")
	if karenaiAPIKey == "" {
		return nil, fmt.Errorf("KARENAI_API_KEY no está configurada en las variables de entorno. Necesaria para Karenai.click API.")
	}

	log.Println("DEBUG: Intentando obtener recomendaciones de Karenai.click desde:", KARENAI_API_URL)

	req, err := http.NewRequest("GET", KARENAI_API_URL, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear la solicitud HTTP a Karenai.click: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+karenaiAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR HTTP (Karenai.click): Falló la solicitud: %v", err)
		return nil, fmt.Errorf("error al realizar la solicitud HTTP a Karenai.click: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("DEBUG: Karenai.click API - Estado de respuesta HTTP: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta de Karenai.click: %w", err)
	}

	log.Printf("DEBUG: Karenai.click API - Cuerpo de respuesta RAW: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Karenai.click API devolvió un estado de error: %s - Cuerpo de respuesta: %s", resp.Status, string(body))
	}

	var karenaiResp karenaiResponse
	err = json.Unmarshal(body, &karenaiResp)
	if err != nil {
		log.Printf("ERROR JSON (Karenai.click): Falló la decodificación del JSON. Error: %v. Cuerpo recibido: %s", err, string(body))
		return nil, fmt.Errorf("error al decodificar la respuesta JSON de Karenai.click: %w", err)
	}

	log.Printf("DEBUG: Karenai.click API - %d stocks decodificados correctamente.", len(karenaiResp.Items))
	return karenaiResp.Items, nil
}

func GetFinnhubMetricsAndQuote(ticker string) (FinnhubData, error) {
	finnhubAPIKey := os.Getenv("FINNHUB_API_KEY")
	if finnhubAPIKey == "" {
		return FinnhubData{Error: fmt.Errorf("FINNHUB_API_KEY no está configurada")}, fmt.Errorf("FINNHUB_API_KEY no está configurada")
	}

	var finnhubData FinnhubData

	// --- 1. Fetch Metrics (PE Ratio, Dividend Yield, Market Cap) ---
	metricURL := fmt.Sprintf("%s/stock/metric?symbol=%s&metricType=all&token=%s", FINNHUB_BASE_URL, ticker, finnhubAPIKey)
	log.Printf("DEBUG: Finnhub API (metrics) - Intentando obtener métricas para %s desde: %s", ticker, metricURL)

	respMetrics, err := http.Get(metricURL)
	if err != nil {
		finnhubData.Error = fmt.Errorf("error al consultar métricas de Finnhub para %s: %w", ticker, err)
		log.Printf("ERROR: Finnhub API (metrics) - Error al hacer la solicitud para %s: %v", ticker, err)
	} else {
		defer respMetrics.Body.Close()
		bodyMetrics, err := io.ReadAll(respMetrics.Body)
		if err != nil {
			finnhubData.Error = fmt.Errorf("error al leer el cuerpo de la respuesta de Finnhub metrics: %w", err)
			log.Printf("ERROR: Finnhub API (metrics) - Error al leer el cuerpo de la respuesta para %s: %v", ticker, err)
		} else {
			log.Printf("DEBUG: Finnhub API (metrics) - Estado HTTP para %s: %d %s", ticker, respMetrics.StatusCode, respMetrics.Status)
			log.Printf("DEBUG: Finnhub API (metrics) - Cuerpo RAW para %s: %s", ticker, string(bodyMetrics))

			if respMetrics.StatusCode != http.StatusOK {
				finnhubData.Error = fmt.Errorf("Finnhub métricas API devolvió estado de error para %s: %s - Cuerpo: %s", ticker, respMetrics.Status, string(bodyMetrics))
				log.Printf("ADVERTENCIA: %v", finnhubData.Error)
			} else {
				var metricData FinnhubMetricResponse
				err = json.Unmarshal(bodyMetrics, &metricData)
				if err != nil {
					finnhubData.Error = fmt.Errorf("error al decodificar JSON de métricas de Finnhub para %s: %w", ticker, err)
					log.Printf("ERROR: %v. Cuerpo: %s", finnhubData.Error, string(bodyMetrics))
				} else {

					if metricData.Metric.PeExclExtraTTM != 0 {
						finnhubData.PE_Ratio = metricData.Metric.PeExclExtraTTM
					} else {
						finnhubData.PE_Ratio = metricData.Metric.PeRatio
					}
					if metricData.Metric.DividendYield != 0 {
						finnhubData.DividendYield = metricData.Metric.DividendYield
					} else {
						finnhubData.DividendYield = metricData.Metric.DividendYieldAlt
					}
					finnhubData.MarketCapitalization = metricData.Metric.MarketCap
				}
			}
		}
	}

	quoteURL := fmt.Sprintf("%s/quote?symbol=%s&token=%s", FINNHUB_BASE_URL, ticker, finnhubAPIKey)
	log.Printf("DEBUG: Finnhub API (quote) - Intentando obtener cotización para %s desde: %s", ticker, quoteURL)

	respQuote, err := http.Get(quoteURL)
	if err != nil {
		finnhubData.Error = fmt.Errorf("error al consultar cotización de Finnhub para %s: %w. %v", ticker, err, finnhubData.Error) // Combine errors
		log.Printf("ERROR: Finnhub API (quote) - Error al hacer la solicitud para %s: %v", ticker, err)
	} else {
		defer respQuote.Body.Close()
		bodyQuote, err := io.ReadAll(respQuote.Body)
		if err != nil {
			finnhubData.Error = fmt.Errorf("error al leer el cuerpo de la respuesta de Finnhub quote: %w. %v", err, finnhubData.Error) // Combine errors
			log.Printf("ERROR: Finnhub API (quote) - Error al leer el cuerpo de la respuesta de cotización para %s: %v", ticker, err)
		} else {
			log.Printf("DEBUG: Finnhub API (quote) - Estado HTTP para %s: %d %s", ticker, respQuote.StatusCode, respQuote.Status)
			log.Printf("DEBUG: Finnhub API (quote) - Cuerpo RAW para %s: %s", ticker, string(bodyQuote))

			if respQuote.StatusCode != http.StatusOK {
				finnhubData.Error = fmt.Errorf("Finnhub cotización API devolvió estado de error para %s: %s - Cuerpo: %s. %v", ticker, respQuote.Status, string(bodyQuote), finnhubData.Error) // Combine errors
				log.Printf("ADVERTENCIA: %v", finnhubData.Error)
			} else {
				var quoteData FinnhubQuoteResponse
				err = json.Unmarshal(bodyQuote, &quoteData)
				if err != nil {
					finnhubData.Error = fmt.Errorf("error al decodificar JSON de cotización de Finnhub para %s: %w. %v", ticker, err, finnhubData.Error) // Combine errors
					log.Printf("ERROR: %v. Cuerpo: %s", finnhubData.Error, string(bodyQuote))
				} else {
					finnhubData.CurrentPrice = quoteData.CurrentPrice

					if quoteData.Timestamp != 0 {
						finnhubData.LatestTradingDay = time.Unix(quoteData.Timestamp, 0)
					}
				}
			}
		}
	}

	return finnhubData, finnhubData.Error
}

func GetAlphaAndLatestTradingDayFromAlphaVantage(ticker string) (AlphaVantageData, error) {
	alphaVantageAPIKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if alphaVantageAPIKey == "" {
		return AlphaVantageData{Error: fmt.Errorf("ALPHA_VANTAGE_API_KEY no está configurada")}, fmt.Errorf("ALPHA_VANTAGE_API_KEY no está configurada")
	}

	time.Sleep(15 * time.Second)

	url := fmt.Sprintf("%s?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", ALPHA_VANTAGE_BASE_URL, ticker, alphaVantageAPIKey)
	log.Printf("DEBUG: Alpha Vantage API - Intentando obtener datos para %s desde: %s", ticker, url)

	var avData AlphaVantageData

	resp, err := http.Get(url)
	if err != nil {
		avData.Error = fmt.Errorf("error al consultar Alpha Vantage para %s: %w", ticker, err)
		log.Printf("ERROR: Alpha Vantage API - Error al hacer la solicitud para %s: %v", ticker, err)
		return avData, avData.Error
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		avData.Error = fmt.Errorf("error al leer el cuerpo de la respuesta de Alpha Vantage: %w", err)
		log.Printf("ERROR: Alpha Vantage API - Error al leer el cuerpo de la respuesta para %s: %v", ticker, err)
		return avData, avData.Error
	}

	log.Printf("DEBUG: Alpha Vantage API - Estado HTTP para %s: %d %s", ticker, resp.StatusCode, resp.Status)
	log.Printf("DEBUG: Alpha Vantage API - Cuerpo RAW para %s: %s", ticker, string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		avData.Error = fmt.Errorf("Alpha Vantage API devolvió estado de error para %s: %s. Cuerpo: %s", ticker, resp.Status, string(bodyBytes))
		log.Printf("ADVERTENCIA: %v", avData.Error)
		return avData, avData.Error
	}

	var avResponse map[string]interface{}
	err = json.Unmarshal(bodyBytes, &avResponse)
	if err != nil {
		avData.Error = fmt.Errorf("error al decodificar respuesta JSON de Alpha Vantage para %s: %w", ticker, err)
		log.Printf("ERROR: %v. Cuerpo: %s", avData.Error, string(bodyBytes))
		return avData, avData.Error
	}

	if errorMessage, ok := avResponse["Error Message"].(string); ok {
		avData.Error = fmt.Errorf("Alpha Vantage API error: %s", errorMessage)
		log.Printf("ADVERTENCIA: %v. Se usarán 0.0 para Alpha y fecha inválida.", avData.Error)
		return avData, avData.Error
	}
	if note, ok := avResponse["Note"].(string); ok {
		avData.Error = fmt.Errorf("Alpha Vantage API note/warning: %s", note)
		log.Printf("ADVERTENCIA: %v. Se usarán 0.0 para Alpha y fecha inválida.", avData.Error)
		return avData, avData.Error
	}

	// Parse Global Quote data
	if globalQuote, ok := avResponse["Global Quote"].(map[string]interface{}); ok {
		// Extract Latest Trading Day (still here for completeness, but Finnhub is prioritized in enricher)
		if ltDayStr, found := globalQuote["07. latest trading day"].(string); found && ltDayStr != "" {
			parsedTime, parseErr := time.Parse("2006-01-02", ltDayStr) // Alpha Vantage format: YYYY-MM-DD
			if parseErr != nil {
				log.Printf("ERROR: Alpha Vantage API - Error al parsear fecha '%s' para %s: %v", ltDayStr, ticker, parseErr)
				avData.Error = fmt.Errorf("error al parsear '07. latest trading day': %w", parseErr)
			} else {
				avData.LatestTradingDay = parsedTime
			}
		} else {
			log.Printf("ADVERTENCIA: Alpha Vantage API - '07. latest trading day' no encontrado o vacío para %s.", ticker)
		}

	} else {
		avData.Error = fmt.Errorf("Alpha Vantage API - 'Global Quote' no encontrado en la respuesta para %s", ticker)
		log.Printf("ADVERTENCIA: %v", avData.Error)
	}

	return avData, avData.Error
}
