<script setup lang="ts">
// Importa la función 'useStockStore' para interactuar con el almacén de estado global.
import { useStockStore } from '../stores/stocks';
// Importa la interfaz 'Stock' para el tipado estricto de las propiedades.
import type { Stock } from '../types/stock';

// Define la forma (tipos y nombres) de las propiedades que el componente espera recibir.
interface Props {
  stock: Stock;
}

// Declara las propiedades del componente. 'stock' será de tipo 'Props'.
const props = defineProps<Props>();
// Inicializa el almacén para que el componente pueda acceder a los datos y métodos globales.
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
      <span class="font-semibold">Rating:</span>
      {{ stock.rating_from && stock.rating_from !== "" ? stock.rating_from : 'N/A' }} to
      {{ stock.rating_to && stock.rating_to !== "" ? stock.rating_to : 'N/A' }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Target Price:</span>
      {{ stockStore.formatCurrency(stock.target_from) }} -
      {{ stockStore.formatCurrency(stock.target_to) }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Current Price:</span> {{ stockStore.formatCurrency(stock.current_price) }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">PE Ratio:</span>
      {{ stock.pe_ratio !== null ? stock.pe_ratio.toFixed(2) : '0.00' }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Dividend Yield:</span>
      {{ stockStore.formatPercentage(stock.dividend_yield) }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Capitalización de Mercado:</span> {{ stockStore.formatMarketCap(stock.market_capitalization) }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Última Actualización:</span> {{ stockStore.formatDate(stock.latest_trading_day) }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Alpha:</span> {{ stockStore.formatPercentage(stock.alpha) }}
    </p>
    <p class="text-gray-700 mb-2">
      <span class="font-semibold">Score Recom.:</span>
      {{ stock.recommendation_score !== null ? stock.recommendation_score.toFixed(2) : '0.00' }}
    </p>
  </div>
</template>

<style scoped>

</style>
