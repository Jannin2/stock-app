// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import StockDetailView from '../views/StockDetailView.vue' 

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    
    {
      path: '/stocks/:id', // Dynamic route for stock ID
      name: 'stock-details',
      component: StockDetailView 
    },
    
  ]
})

export default router