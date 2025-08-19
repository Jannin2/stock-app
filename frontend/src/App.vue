<script setup lang="ts">
// CORRECTED: Use default imports for components defined with <script setup>
import RecommendedStocks from './components/RecommendedStocks.vue';
import StocksTable from './components/StocksTable.vue';
</script>

<template>
  <div id="app" class="min-h-screen bg-gray-100 text-gray-800">
    <header class="bg-green-600 text-white p-6 shadow-lg text-center">
      <h1 class="text-4xl font-extrabold tracking-tight">Mercado de Acciones</h1>
    </header>

    <main class="container mx-auto p-6 mt-8 bg-white shadow-xl rounded-lg">
      <RecommendedStocks />
      <hr class="my-8 border-gray-200" />
      <StocksTable />
    </main>
  </div>
</template>

<style scoped>
/* No es necesario mucho CSS aquí si usas Tailwind para todo */
/* Sin embargo, puedes añadir estilos específicos que no sean de Tailwind si lo necesitas */
</style>