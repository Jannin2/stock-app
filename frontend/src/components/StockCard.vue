<script setup lang="ts">
import { useStockStore } from '../stores/stocks';
import type { Stock } from '../types/stock';

interface Props {
  stock: Stock;
}

const props = defineProps<Props>();
const stockStore = useStockStore();

</script>

<template>
  <div class="bg-white rounded-lg shadow-lg p-6 mb-6">
    <h2 class="text-2xl font-bold text-gray-800 mb-4">{{ stock.company }} ({{ stock.ticker }})</h2>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Acción:</span> {{ stock.action }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Brokerage:</span> {{ stock.brokerage }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Target Hasta:</span>
      {{ stockStore.formatCurrency(stock.targetTo === 0 ? null : stock.targetTo) }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Precio Actual:</span> {{ stockStore.formatCurrency(stock.currentPrice) }}
    </p>

    <p class="text-gray-700 mb-2">
      <span class="font-semibold">PE Ratio:</span>
      {{ stock.peRatio !== null && stock.peRatio !== 0 ? stock.peRatio.toFixed(2) : 'N/A' }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Dividend Yield:</span>
      {{ stockStore.formatPercentage(stock.dividendYield) }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Score Recom.:</span>
      {{ stock.recommendationScore !== null ? stock.recommendationScore.toFixed(2) : 'N/A' }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Capitalización de Mercado:</span> {{ stockStore.formatMarketCap(stock.marketCap) }}
      </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Última Actualización:</span> {{ stockStore.formatDate(stock.latestTradingDay) }}
    </p>
  </div>
</template>

<style scoped>

</style>