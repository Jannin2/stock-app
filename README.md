Aplicación de stock
¡Bienvenido a la Stock App ! Esta es una aplicación web full-stack diseñada para ofrecer una visión integral del mercado de valores. Permite a los usuarios explorar datos de acciones, obtener recomendaciones en tiempo real y visualizar métricas financieras clave.

✨ Características Principales
Listado de Stocks: Visualiza una lista completa de acciones con opciones de búsqueda, filtrado y paginación.
Recomendaciones: Obtén recomendaciones de stocks basadas en el análisis de una API externa.
Actualización de Datos en Tiempo Real (simulada): Un job de cron en el backend se encarga de enriquecer los datos de stocks con información actualizada de diversas fuentes externas.
🚀 Tecnologías Utilizadas
Este proyecto se construye con un enfoque en la eficiencia y la modularidad, utilizando un stack tecnológico moderno:

Backend (Ir)
GoLang: El lenguaje de programación principal, conocido por su concurrencia y rendimiento, ideal para servicios de backend.
Chi: Un enrutador HTTP ligero y rápido que facilita la creación de APIs RESTful.
PostgreSQL (CockroachDB): Una base de datos distribuida, compatible con PostgreSQL, utilizada para un almacenamiento de datos robusto y escalable.
godotenv: Para la carga segura de variables de entorno en el entorno de desarrollo local, manteniendo las credenciales fuera del control de versiones.
API externas:
Karenai.click: Fuente principal para la obtención inicial de stocks y sus recomendaciones.
Finnhub.io: Provee métricas financieras detalladas (PE Ratio, Dividend Yield, Market Capitalization) y datos de cotización actuales.
Alpha Vantage: Utilizada para obtener datos de mercado adicionales, como el valor "Alpha" y el último día de negociación.
Interfaz (Vue.js)
Vue.js 3: Un framework JavaScript progresivo para construir interfaces de usuario interactivas y reactivas.
Vite: Un build tool de nueva generación que ofrece un entorno de desarrollo extremadamente rápido y optimiza el proceso de empaquetado para producción.
Tailwind CSS: Un framework CSS de utilidad que permite un diseño rápido y altamente personalizable directamente en el marcado HTML.
⚙️ Configuración y Ejecución Local
Sigue estos pasos para levantar y ejecutar la aplicación en tu máquina local.

Prerrequisitos
Antes de empezar, asegúrate de tener instalados los siguientes componentes:

Go (versión 1.20 o superior)

Node.js (versión LTS recomendada) y su gestor de paquetes preferido (se recomienda pnpm, pero npm o yarn también son válidos).

Docker y Docker Compose: Esencial para ejecutar la base de datos CockroachDB en un contenedor.

API Keys: Necesitarás obtener tus propias claves API de las siguientes plataformas. Son fundamentales para que la aplicación obtenga datos externos:

Karenai.click: Visitahttps://api.karenai.click/ para obtener tu clave.
Finnhub.io: Regístrate en https://finnhub.io/ para obtener tu token API.
Alpha Vantage: Consigue tu clave API en https://www.alphavantage.co/.
1. Clonar el Repositorio
Primero, clona este repositorio en tu máquina local:

Intento

git clone https://github.com/jannin2/stock-app.git # Asegúrate de usar la URL correcta de tu repositorio
cd stock-app
2. Configurar Variables de Entorno
Crea un nuevo archivo llamado .env dentro del directorio backend/ . Este archivo almacenará tus claves API y la configuración de la base de datos. Rellénalo con los valores que obtuviste de las APIs:

backend/.env
DATABASE_URL="postgres://root@localhost:26257/defaultdb?sslmode=disable" KARENAI_API_KEY="TU_CLAVE_KARENAI_AQUI" FINNHUB_API_KEY="TU_CLAVE_FINNHUB_AQUI" ALPHA_VANTAGE_API_KEY="TU_CLAVE_ALPHA_VANTAGE_AQUI" PORT="8081" # Puerto en el que el backend escuchará

🚨 ¡Importante! El archivo .env está en el .gitignore del proyecto y nunca debe ser subido a tu repositorio público por razones de seguridad.

3. Iniciar la Base de Datos (CockroachDB con Docker Compose)
Desde la raíz de tu proyecto (stock-app/), ejecuta el siguiente comando para iniciar el contenedor de CockroachDB:

Intento

docker-compose up -d cockroachdb
Esto levantará la base de datos en segundo plano. Podrás acceder a la interfaz de administración de CockroachDB en tu navegador a través de http://localhost:8080.

4. Ejecutar el Backend (Go)
Ahora, navega al directorio backend e inicia la aplicación Go:

Intento

cd backend
go mod tidy # Asegúrate de que todas las dependencias estén instaladas y actualizadas
go run main.go
El backend se conectará a CockroachDB, inicializará el esquema de la base de datos (creando las tablas necesarias si no existen) y comenzará a escuchar las peticiones HTTP en http://localhost:8081. Verás mensajes de log indicando el inicio del servidor y el job de cron para el enriquecimiento de datos.

5. Ejecutar el Frontend (Vue.js)
En una nueva terminal , navega al directorio frontend. Instala las dependencias y luego inicia el servidor de desarrollo de Vue.js:

Intento

cd ../frontend # Si estás en el directorio 'backend', vuelve a la raíz del proyecto y luego a 'frontend'
pnpm install   # O usa 'npm install' o 'yarn install'
pnpm dev       # O usa 'npm run dev' o 'yarn dev'
El servidor de desarrollo del frontend se iniciará, generalmente en http://localhost:5173. Abre esta URL en tu navegador web.

¡Con ambos servicios en ejecución, la aplicación estará completamente funcional y lista para ser explorada!

🧪 Ejecución de Pruebas
El proyecto incluye pruebas para asegurar la calidad y el correcto funcionamiento de ambos servicios.

Pruebas del Backend (Go)
Desde el directorio backend:

Intento

go test ./...
Pruebas del Frontend (Vue.js)
Desde el directorio frontend (asumiendo que uses Vitest, Jest u otra herramienta de testing configurada):

Intento

pnpm test # O 'npm run test' o 'yarn test'
