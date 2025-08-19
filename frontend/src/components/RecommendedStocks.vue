<script setup lang="ts">
import { onMounted } from 'vue';
import { useStockStore } from '../stores/stocks';

const stockStore = useStockStore();

onMounted(() => {
    stockStore.fetchRecommendedStocks();
});

const getCardClass = (action: string) => {
    switch (action.toLowerCase()) {
        case 'buy':
        case 'strong buy':
            return 'card-buy';
        case 'sell':
        case 'strong sell':
            return 'card-sell';
        case 'hold':
            return 'card-hold';
        default:
            // This case handles 'target lowered by', 'reiterated by', 'target raised by', etc.
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
/* Tu estilo existente */
.recommended-stocks-container {
    padding: 20px;
    background-color: #f9f9f9;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-top: 20px;
}

h2 {
    color: #333;
    margin-bottom: 20px;
    text-align: center;
}

.recommended-cards {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
}

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

.stock-card h3 {
    margin-top: 0;
    color: #007bff;
    font-size: 1.2em;
    margin-bottom: 10px;
}

.stock-card p {
    margin: 5px 0;
    color: #555;
    font-size: 0.95em;
}

.action-text {
    font-weight: bold;
}

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