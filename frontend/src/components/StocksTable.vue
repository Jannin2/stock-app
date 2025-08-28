<script setup lang="ts">
// Importa las funciones reactivas de Vue 3: ref, computed, y onMounted.
import { ref, computed, onMounted } from 'vue';
// Importa el 'store' de acciones para manejar el estado global.
import { useStockStore } from '../stores/stocks';
// Importa la interfaz de tipo 'Stock' para garantizar un tipado estricto.
import type { Stock } from '../types/stock';

// Inicializa el 'store' para poder acceder a su estado y sus acciones.
const stockStore = useStockStore();

// Llama al 'hook' del ciclo de vida 'onMounted'.
onMounted(() => {
    // Cuando el componente se monta, llama a la acción 'fetchStocks' para obtener los datos de las acciones.
    stockStore.fetchStocks();
});

// Crea una referencia reactiva para el texto de búsqueda.
const searchQuery = ref('');
// Define una propiedad 'computed' para filtrar las acciones.
const filteredStocks = computed(() => {
    // Si no hay stocks o la lista está vacía, devuelve un array vacío.
    if (!stockStore.stocks || stockStore.stocks.length === 0) {
        return [];
    }
    // Convierte el valor de búsqueda a minúsculas para una comparación sin distinción de mayúsculas y minúsculas.
    const query = searchQuery.value.toLowerCase();
    // Filtra la lista de stocks.
    return stockStore.stocks.filter(stock => {
        // Obtiene las propiedades de stock y las convierte a minúsculas, manejando posibles valores nulos.
        const ticker = stock.ticker ? stock.ticker.toLowerCase() : '';
        const company = stock.company ? stock.company.toLowerCase() : '';
        const brokerage = stock.brokerage ? stock.brokerage.toLowerCase() : '';

        // Retorna 'true' si el texto de búsqueda se encuentra en el ticker, la compañía o la correduría.
        return (
            ticker.includes(query) ||
            company.includes(query) ||
            brokerage.includes(query)
        );
    });
});

// Define una propiedad 'computed' para manejar la paginación.
const paginatedStocks = computed(() => {
    // Por ahora, solo devuelve la lista de stocks filtrada.
    return filteredStocks.value;
});

// Define una propiedad 'computed' para el estado de carga del 'store'.
const isLoading = computed(() => stockStore.loading);
// Define una propiedad 'computed' para el mensaje de error del 'store'.
const errorMessage = computed(() => stockStore.error);

</script>

<template>
    <div class="stocks-table-container">
        <h2>All Stocks</h2>
        <input type="text" v-model="searchQuery" placeholder="Search by ticker, company, or brokerage..." class="search-input" />

        <div v-if="isLoading" class="loading-message">
            Loading all stocks... ⏳
        </div>
        <div v-else-if="errorMessage" class="error-message">
            Error: {{ errorMessage }} 🔴
        </div>
        <div v-else-if="filteredStocks.length === 0" class="no-data-message">
            No stocks found. Please try a different search or wait for data. 🤷‍♀️
        </div>
        <table v-else class="stocks-table">
            <thead>
                <tr>
                    <th>Ticker</th>
                    <th>Company</th>
                    <th>Brokerage</th>
                    <th>Action</th>
                    <th>Rating From</th>
                    <th>Rating To</th>
                    <th>Target From</th>
                    <th>Target To</th>
                    <th>Current Price</th>
                    <th>PE Ratio</th>
                    <th>Dividend Yield</th>
                    <th>Market Cap</th>
                    <th>Alpha</th>
                    <th>Latest Trading Day</th>
                    <th>Rec. Score</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="stock in paginatedStocks" :key="stock.id">
                    <td>{{ stock.ticker }}</td>
                    <td>{{ stock.company }}</td>
                    <td>{{ stock.brokerage }}</td>
                    <td>{{ stock.action }}</td>
                    <td>
                        {{ stock.rating_from && stock.rating_from !== "" ? stock.rating_from : 'N/A' }}
                    </td>
                    <td>
                        {{ stock.rating_to && stock.rating_to !== "" ? stock.rating_to : 'N/A' }}
                    </td>
                    <td>{{ stockStore.formatCurrency(stock.target_from) }}</td>
                    <td>{{ stockStore.formatCurrency(stock.target_to) }}</td>
                    <td>{{ stockStore.formatCurrency(stock.current_price) }}</td>
                    <td>
                        {{ typeof stock.pe_ratio === 'number' && !isNaN(stock.pe_ratio) ? stock.pe_ratio.toFixed(2) : 'N/A' }}
                    </td>
                    <td>{{ stockStore.formatPercentage(stock.dividend_yield) }}</td>
                    <td>{{ stockStore.formatMarketCap(stock.market_capitalization) }}</td>
                    <td>{{ stockStore.formatPercentage(stock.alpha) }}</td>
                    <td>{{ stockStore.formatDate(stock.latest_trading_day) }}</td>
                    <td>
                        {{ typeof stock.recommendation_score === 'number' && !isNaN(stock.recommendation_score) ? stock.recommendation_score.toFixed(2) : 'N/A' }}
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<style scoped>
/* Contenedor principal de la tabla. */
.stocks-table-container {
    padding: 20px;
    background-color: #f9f9f9;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 20px;
}

/* Estilo para el título de la sección. */
h2 {
    color: #333;
    margin-bottom: 15px;
    text-align: center;
}

/* Estilo para el campo de búsqueda. */
.search-input {
    width: 100%;
    padding: 10px;
    margin-bottom: 20px;
    border: 1px solid #ddd;
    border-radius: 4px;
    box-sizing: border-box;
}

/* Estilo principal de la tabla. */
.stocks-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 20px;
}

/* Estilos de las celdas y encabezados de la tabla. */
.stocks-table th,
.stocks-table td {
    border: 1px solid #e0e0e0;
    padding: 12px 15px;
    text-align: left;
}

/* Estilo para los encabezados de la tabla. */
.stocks-table th {
    background-color: #007bff;
    color: white;
    font-weight: bold;
    text-transform: uppercase;
    font-size: 0.9em;
}

/* Estilo para las filas pares de la tabla para mejorar la legibilidad. */
.stocks-table tr:nth-child(even) {
    background-color: #f2f2f2;
}

/* Estilo de la fila al pasar el cursor sobre ella. */
.stocks-table tr:hover {
    background-color: #e9e9e9;
}

/* Estilo para los mensajes de estado (carga, error, sin datos). */
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
