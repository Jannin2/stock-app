// frontend/src/stores/stocks.ts
// Importa 'defineStore' de Pinia, la herramienta de gestión de estado global para Vue.
import { defineStore } from 'pinia';
// Importa la interfaz 'Stock' para el tipado estricto de los datos.
import type { Stock } from '../types/stock';

// Obtiene la URL de la API desde las variables de entorno de Vite.
const API_URL = import.meta.env.VITE_API_URL;

// Define una interfaz para el estado del 'store', garantizando tipos de datos consistentes.
interface StockState {
    stocks: Stock[]; // Array para todas las acciones.
    recommendedStocks: Stock[]; // Array para las acciones recomendadas.
    selectedStock: Stock | null; // Objeto para la acción seleccionada, puede ser nulo.
    loading: boolean; // Indicador de estado de carga.
    error: string | null; // Mensaje de error, puede ser nulo.
}

// Define y exporta el 'store' de acciones.
export const useStockStore = defineStore('stock', {
    state: (): StockState => ({
        stocks: [],
        recommendedStocks: [],
        selectedStock: null,
        loading: false,
        error: null,
    }),
    actions: {
        async fetchStocks() {
            this.loading = true;
            this.error = null;
            try {
                // Usa la variable de entorno para la URL base.
                const response = await fetch(`${API_URL}/api/v1/stocks`);
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const data: Stock[] = await response.json();
                this.stocks = data;
            } catch (e: any) {
                this.error = e.message;
                console.error('Error fetching stocks:', e);
            } finally {
                this.loading = false;
            }
        },
        async fetchRecommendedStocks() {
            this.loading = true;
            this.error = null;
            try {
                // Usa la variable de entorno para la URL base.
                const response = await fetch(`${API_URL}/api/v1/stocks/recommended`);
                if (!response.ok) {
                    const errorText = await response.text();
                    throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
                }
                const data: Stock[] = await response.json();
                this.recommendedStocks = data;
            } catch (e: any) {
                this.error = e.message;
                console.error('Error fetching recommended stocks:', e);
            } finally {
                this.loading = false;
            }
        },
        async fetchStockDetails(id: string) {
            this.loading = true;
            this.error = null;
            this.selectedStock = null;
            try {
                // Usa la variable de entorno para la URL base.
                const response = await fetch(`${API_URL}/api/v1/stocks/${id}`);
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const data: Stock = await response.json();
                this.selectedStock = data;
            } catch (e: any) {
                this.error = e.message;
                console.error(`Error fetching stock details for ${id}:`, e);
            } finally {
                this.loading = false;
            }
        },
        formatCurrency(value: number | null | undefined): string {
            const numValue = (value === null || typeof value === 'undefined' || isNaN(value)) ? 0 : value;
            return new Intl.NumberFormat('en-US', {
                style: 'currency',
                currency: 'USD',
                minimumFractionDigits: 2,
                maximumFractionDigits: 2
            }).format(numValue);
        },
        formatPercentage(value: number | null | undefined): string {
            const numValue = (value === null || typeof value === 'undefined' || isNaN(value)) ? 0 : value;
            return new Intl.NumberFormat('en-US', {
                style: 'percent',
                minimumFractionDigits: 2,
                maximumFractionDigits: 2
            }).format(numValue);
        },
        formatMarketCap(value: number | null | undefined): string {
            const numValue = (value === null || typeof value === 'undefined' || isNaN(value)) ? 0 : value;
            if (numValue === 0) {
                return '$0.00M';
            }
            const absValue = Math.abs(numValue);
            if (absValue >= 1000) {
                return '$' + (numValue / 1000).toFixed(2) + 'B';
            } else {
                return '$' + numValue.toFixed(2) + 'M';
            }
        },
        formatDate(dateString: string | null | undefined): string {
            if (!dateString || dateString === "0001-01-01T00:00:00Z" || dateString.startsWith("0000-")) {
                return 'N/A';
            }
            try {
                const date = new Date(dateString);
                if (isNaN(date.getTime())) {
                    return 'N/A';
                }
                return new Intl.DateTimeFormat('en-US', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric'
                }).format(date);
            } catch (e) {
                console.error("Error formatting date:", e);
                return 'N/A';
            }
        },
    },
});
