<script setup lang="ts">
// Importa la función 'onMounted' de Vue. Este es un 'hook' del ciclo de vida que se ejecuta cuando el componente se ha añadido al DOM.
import { onMounted } from 'vue';
// Importa 'useStockStore', una función que te permite acceder a un 'store' (almacén de estado) global.
import { useStockStore } from '../stores/stocks';

// Inicializa el 'store' para poder acceder a sus datos y métodos.
const stockStore = useStockStore();

// Llama a la función 'onMounted'. El código dentro de esta función se ejecutará automáticamente después de que el componente sea montado.
onMounted(() => {
    // Llama a la acción 'fetchRecommendedStocks' del 'store' para obtener los datos de las acciones desde la API.
    stockStore.fetchRecommendedStocks();
});

// Define una función para determinar qué clase de estilo CSS aplicar a una tarjeta de acción.
const getCardClass = (action: string) => {
    // Utiliza una declaración 'switch' para decidir la clase según el valor de 'action'.
    switch (action.toLowerCase()) {
        // Para las acciones de compra, devuelve la clase 'card-buy'.
        case 'buy':
        case 'strong buy':
            return 'card-buy';
        // Para las acciones de venta, devuelve la clase 'card-sell'.
        case 'sell':
        case 'strong sell':
            return 'card-sell';
        // Para las acciones de 'mantener', devuelve la clase 'card-hold'.
        case 'hold':
            return 'card-hold';
        default:
            // Para cualquier otra acción no especificada, devuelve una clase de estilo neutral.
            return 'card-neutral';
    }
};

</script>

<template>
    <div class="recommended-stocks-container">
        <h2>Recommended Stocks</h2>

        <div v-if="stockStore.loading" class="loading-message">
            Loading recommended stocks... ⏳
        </div>
        <div v-else-if="stockStore.error" class="error-message">
            Error: {{ stockStore.error }} 🔴
        </div>
        <div v-else-if="(stockStore.recommendedStocks || []).length === 0" class="no-data-message">
            No recommended stocks found. 🤷‍♂️
        </div>
        <div v-else class="recommended-cards">
            <div v-for="stock in stockStore.recommendedStocks" :key="stock.id" :class="['stock-card', getCardClass(stock.action)]">
                <h3>{{ stock.ticker }} ({{ stock.company }})</h3>
                <p><strong>Brokerage:</strong> {{ stock.brokerage }}</p>
                <p><strong>Action:</strong> <span class="action-text">{{ stock.action }}</span></p>
                <p><strong>Rating:</strong>
                    {{ stock.rating_from && stock.rating_from !== "" ? stock.rating_from : 'N/A' }} to
                    {{ stock.rating_to && stock.rating_to !== "" ? stock.rating_to : 'N/A' }}
                </p>
                <p><strong>Target Price:</strong>
                    {{ stockStore.formatCurrency(stock.target_from) }} -
                    {{ stockStore.formatCurrency(stock.target_to) }}
                </p>
                <p><strong>Current Price:</strong> {{ stockStore.formatCurrency(stock.current_price) }}</p>

                <p><strong>PE Ratio:</strong>
                    {{ typeof stock.pe_ratio === 'number' && !isNaN(stock.pe_ratio) ? stock.pe_ratio.toFixed(2) : 'N/A' }}
                </p>

                <p><strong>Dividend Yield:</strong> {{ stockStore.formatPercentage(stock.dividend_yield) }}</p>
                <p><strong>Market Cap:</strong> {{ stockStore.formatMarketCap(stock.market_capitalization) }}</p>

                <p><strong>Alpha:</strong> {{ stockStore.formatPercentage(stock.alpha) }}</p>

                <p><strong>Latest Trading Day:</strong> {{ stockStore.formatDate(stock.latest_trading_day) }}</p>
                <p><strong>Recommendation Score:</strong>
                    {{ typeof stock.recommendation_score === 'number' && !isNaN(stock.recommendation_score) ? stock.recommendation_score.toFixed(2) : 'N/A' }}
                </p>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Contenedor principal de la sección. */
.recommended-stocks-container {
    padding: 20px;
    background-color: #f9f9f9;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-top: 20px;
}

/* Estilo para el título de la sección. */
h2 {
    color: #333;
    margin-bottom: 20px;
    text-align: center;
}

/* Contenedor que utiliza un 'grid' para mostrar las tarjetas. */
.recommended-cards {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
}

/* Estilo para cada tarjeta individual. */
.stock-card {
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
    padding: 20px;
    border-left: 5px solid;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}

/* Estilo para el título de la tarjeta. */
.stock-card h3 {
    margin-top: 0;
    color: #007bff;
    font-size: 1.2em;
    margin-bottom: 10px;
}

/* Estilo para los párrafos de la tarjeta. */
.stock-card p {
    margin: 5px 0;
    color: #555;
    font-size: 0.95em;
}

/* Estilo para el texto de la acción. */
.action-text {
    font-weight: bold;
}

/* Clases dinámicas que cambian el color del borde izquierdo según la recomendación. */
.card-buy {
    border-left-color: #28a745;
}

.card-sell {
    border-left-color: #dc3545;
}

.card-hold {
    border-left-color: #ffc107;
}

.card-neutral {
    border-left-color: #6c757d;
}

/* Estilos para los mensajes de estado. */
.loading-message, .error-message, .no-data-message {
    text-align: center;
    padding: 20px;
    font-size: 1.1em;
    color: #555;
}

/* Estilo específico para el mensaje de error. */
.error-message {
    color: #d9534f;
    font-weight: bold;
}
</style>
