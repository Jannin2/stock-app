package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/jannin2/stock-app/backend/api"
	enricher "github.com/jannin2/stock-app/backend/cron"
	"github.com/jannin2/stock-app/backend/database"
	"github.com/jannin2/stock-app/backend/handlers"
)

func main() {
	// 0. Cargar el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env. Asegúrate de que las variables de entorno estén configuradas o se usarán los valores por defecto.")
	}

	// 1. Conectar a la base de datos
	dbConn, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Error al conectar a la base de datos: %v", err)
	}
	defer database.CloseDB(dbConn)

	// 2. Inicializar el esquema de la base de datos
	if err = database.InitSchema(dbConn); err != nil {
		log.Fatalf("❌ Error al inicializar el esquema de la base de datos: %v", err)
	}

	// 3. Crear una instancia del cliente de base de datos
	dbClient := database.NewStockDB(dbConn)

	// 4. Inicializar los manejadores de HTTP
	stockHandlers := handlers.NewStockHandlers(dbClient)

	// 5. Inicializar el job de cron
	enricherJob := enricher.NewEnricher(dbClient)
	go enricherJob.StartFetching()

	// 6. Configurar el router HTTP
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)


	router.Use(cors.Handler(cors.Options{
		
		AllowedOrigins:   []string{"http://localhost:5173", "https://stock-app-1-0lc4.onrender.com/"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// --- Rutas corregidas para evitar el 404 ---

	// Handler para el chequeo de estado de Render
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Handler para la ruta principal
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("¡Hola! La API de Stocks está activa. Para la interfaz de usuario, conéctate a la URL del frontend."))
	})

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
