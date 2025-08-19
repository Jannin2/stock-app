// frontend/src/stores/stocks.ts
import { defineStore } from 'pinia';
import type { Stock } from '../types/stock';

interface StockState {
    stocks: Stock[];
    recommendedStocks: Stock[];
    selectedStock: Stock | null;
    loading: boolean;
    error: string | null;
}

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
                const response = await fetch('http://localhost:8081/api/v1/stocks');
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
                const response = await fetch('http://localhost:8081/api/v1/stocks/recommended');
                if (!response.ok) {
                    const errorText = await response.text(); // Capture potential error message from server
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
            this.selectedStock = null; // Clear previous selection
            try {
                const response = await fetch(`http://localhost:8081/api/v1/stocks/${id}`);
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
            // Check if value is a valid number, otherwise default to 0 for formatting
            const numValue = (value === null || typeof value === 'undefined' || isNaN(value)) ? 0 : value;
            return new Intl.NumberFormat('en-US', {
                style: 'currency',
                currency: 'USD',
                minimumFractionDigits: 2,
                maximumFractionDigits: 2
            }).format(numValue);
        },

        formatPercentage(value: number | null | undefined): string {
            // Check if value is a valid number, otherwise default to 0 for formatting
            const numValue = (value === null || typeof value === 'undefined' || isNaN(value)) ? 0 : value;
            return new Intl.NumberFormat('en-US', {
                style: 'percent',
                minimumFractionDigits: 2,
                maximumFractionDigits: 2
            }).format(numValue / 100); // Divide by 100 if your backend provides raw percentage (e.g., 5 for 5%)
            // If your backend provides it as a decimal (e.g., 0.05 for 5%), remove the / 100
        },

        formatMarketCap(value: number | null | undefined): string {
            // Check if value is a valid number, otherwise default to 0 for formatting
            const numValue = (value === null || typeof value === 'undefined' || isNaN(value)) ? 0 : value;

            if (numValue === 0) {
                return '$0.00M'; // Return 0.00M for zero value market cap
            }

            const absValue = Math.abs(numValue);

            // Assuming the market capitalization from the backend (e.g., 12477.23) is already in MILLIONS.
            if (absValue >= 1000) { // If 1000 Million or more (i.e., 1 Billion)
                return '$' + (numValue / 1000).toFixed(2) + 'B'; // Divide by 1000 to convert Millions to Billions
            } else { // If less than 1000 Million, display in Millions
                return '$' + numValue.toFixed(2) + 'M';
            }
        },

        formatDate(dateString: string | null | undefined): string {
            // Explicitly check for null, undefined, empty string, or the default Go "zero" time
            if (!dateString || dateString === "0001-01-01T00:00:00Z" || dateString.startsWith("0000-")) {
                return 'N/A'; // Return N/A for invalid or default dates
            }
            try {
                const date = new Date(dateString);
                if (isNaN(date.getTime())) { // Check for "Invalid Date" conversion
                    return 'N/A';
                }
                return new Intl.DateTimeFormat('en-US', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric'
                }).format(date);
            } catch (e) {
                console.error("Error formatting date:", e);
                return 'N/A'; // Catch any other parsing errors
            }
        },
    },
});