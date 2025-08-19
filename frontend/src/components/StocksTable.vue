<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useStockStore } from '../stores/stocks';
import type { Stock } from '../types/stock'; // Make sure you have this import and the type is correct

const stockStore = useStockStore();

onMounted(() => {
    stockStore.fetchStocks();
});

const searchQuery = ref('');
const filteredStocks = computed(() => {
    if (!stockStore.stocks || stockStore.stocks.length === 0) {
        return [];
    }
    const query = searchQuery.value.toLowerCase();
    return stockStore.stocks.filter(stock => {
        // Use snake_case for consistency with backend if not transforming
        const ticker = stock.ticker ? stock.ticker.toLowerCase() : '';
        const company = stock.company ? stock.company.toLowerCase() : '';
        const brokerage = stock.brokerage ? stock.brokerage.toLowerCase() : '';

        return (
            ticker.includes(query) ||
            company.includes(query) ||
            brokerage.includes(query)
        );
    });
});

const paginatedStocks = computed(() => {
    // For now, no actual pagination logic, just returns filtered stocks
    return filteredStocks.value;
});

const isLoading = computed(() => stockStore.loading);
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
/* Tu estilo existente */
.stocks-table-container {
    padding: 20px;
    background-color: #f9f9f9;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 20px;
}

h2 {
    color: #333;
    margin-bottom: 15px;
    text-align: center;
}

.search-input {
    width: 100%;
    padding: 10px;
    margin-bottom: 20px;
    border: 1px solid #ddd;
    border-radius: 4px;
    box-sizing: border-box;
}

.stocks-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 20px;
}

.stocks-table th,
.stocks-table td {
    border: 1px solid #e0e0e0;
    padding: 12px 15px;
    text-align: left;
}

.stocks-table th {
    background-color: #007bff;
    color: white;
    font-weight: bold;
    text-transform: uppercase;
    font-size: 0.9em;
}

.stocks-table tr:nth-child(even) {
    background-color: #f2f2f2;
}

.stocks-table tr:hover {
    background-color: #e9e9e9;
}

.loading-message, .error-message, .no-data-message {
    text-align: center;
    padding: 20px;
    font-size: 1.1em;
    color: #555;
}

.error-message {
    color: #d9534f;
    font-weight: bold;
}
</style>