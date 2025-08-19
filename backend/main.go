package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"   // Import the cors package
	"github.com/joho/godotenv" // Import godotenv

	"github.com/jannin2/stock-app/backend/api"
	enricher "github.com/jannin2/stock-app/backend/cron"
	"github.com/jannin2/stock-app/backend/database"
	"github.com/jannin2/stock-app/backend/handlers"
)

func main() {
	// 0. Load .env file at the very beginning of main()
	// This makes environment variables available to subsequent calls like os.Getenv
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env. Asegúrate de que las variables de entorno estén configuradas o se usarán los valores por defecto.")
	}

	// 1. Conectar a la base de datos
	// `err` is already declared by godotenv.Load(), so use `=`
	dbConn, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Error al conectar a la base de datos: %v", err)
	}
	defer database.CloseDB(dbConn)

	// 2. Inicializar el esquema de la base de datos (crear tablas si no existen)
	if err = database.InitSchema(dbConn); err != nil {
		log.Fatalf("❌ Error al inicializar el esquema de la base de datos: %v", err)
	}

	// 3. Crear una instancia del cliente de base de datos que implementa StockDB
	dbClient := database.NewStockDB(dbConn)

	// 4. Inicializar los manejadores de HTTP con la instancia de dbClient
	stockHandlers := handlers.NewStockHandlers(dbClient)

	// 5. Inicializar el job de cron con la instancia de dbClient
	enricherJob := enricher.NewEnricher(dbClient)
	go enricherJob.StartFetching() // Inicia el job de cron en una goroutine

	// 6. Configurar el router HTTP
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// --- Add CORS middleware here. This should be placed BEFORE any specific routes ---
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Allow your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rutas de la API (asumiendo que SetupRouter las define)
	api.SetupRouter(router, stockHandlers)

	// Iniciar el servidor HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("🚀 Servidor escuchando en http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
