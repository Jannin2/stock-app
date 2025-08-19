// src/__tests__/stocks.spec.ts
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useStockStore } from '../stores/stocks';

// Datos de mock que simulan la respuesta de tu backend
const mockStocksData = [
  {
    id: '1',
    ticker: 'TESTA',
    company: 'Test Company A',
    brokerage: 'Broker A',
    action: 'buy',
    rating_from: 'Hold',
    rating_to: 'Buy',
    target_from: 100,
    target_to: 120,
    current_price: 105,
    pe_ratio: 25.5,
    dividend_yield: 0.02, // 2%
    market_capitalization: 1500.25, // en Millones
    alpha: 0.005, // 0.5%
    latest_trading_day: '2025-08-18T16:00:00Z',
    recommendation_score: 8.5,
    created_at: '2025-08-18T15:00:00Z',
    updated_at: '2025-08-18T15:00:00Z',
  },
  {
    id: '2',
    ticker: 'TESTB',
    company: 'Test Company B',
    brokerage: 'Broker B',
    action: 'hold',
    rating_from: null,
    rating_to: null,
    target_from: null,
    target_to: null,
    current_price: 50,
    pe_ratio: null, // Probar PE Ratio nulo
    dividend_yield: 0.0, // 0%
    market_capitalization: 50.75,
    alpha: null, // Probar Alpha nulo
    latest_trading_day: '2025-08-17T16:00:00Z',
    recommendation_score: null, // Probar Recommendation Score nulo
    created_at: '2025-08-17T15:00:00Z',
    updated_at: '2025-08-17T15:00:00Z',
  },
];

describe('Stock Store', () => {
  beforeEach(() => {
    // 1. Resetea Pinia y activa una nueva instancia antes de cada prueba
    setActivePinia(createPinia());
    // 2. Mockea la función global `fetch` para simular respuestas de la API
    vi.stubGlobal('fetch', vi.fn((url) => {
      if (url.includes('/stocks/recommended')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve(mockStocksData),
        });
      }
      if (url.includes('/stocks')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve(mockStocksData), // Usar los mismos datos mock para simplicidad
        });
      }
      return Promise.reject(new Error('URL de API no mockeada'));
    }));
  });

  // --- Pruebas de Acciones Asíncronas (fetch) ---
  it('fetches recommended stocks correctly and updates store', async () => {
    const store = useStockStore();
    expect(store.loading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.recommendedStocks).toEqual([]);

    await store.fetchRecommendedStocks();

    expect(store.loading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.recommendedStocks).toEqual(mockStocksData);
  });

  it('sets error on fetch recommended stocks failure', async () => {
    // Mockear fetch para simular un error de servidor (status: 500)
    vi.stubGlobal('fetch', vi.fn(() => Promise.resolve({ ok: false, status: 500, text: () => Promise.resolve('Server Error Mock') })));
    const store = useStockStore();

    await store.fetchRecommendedStocks();

    expect(store.loading).toBe(false);
    expect(store.error).toBe('HTTP error! status: 500, message: Server Error Mock');
    expect(store.recommendedStocks).toEqual([]);
  });

  it('fetches all stocks correctly and updates store', async () => {
    const store = useStockStore();
    expect(store.loading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.stocks).toEqual([]);

    await store.fetchStocks();

    expect(store.loading).toBe(false);
    expect(store.error).toBeNull();
    expect(store.stocks).toEqual(mockStocksData);
  });

  // --- Pruebas de Funciones de Formato ---
  it('formats currency correctly', () => {
    const store = useStockStore();
    expect(store.formatCurrency(123.456)).toBe('$123.46');
    expect(store.formatCurrency(0)).toBe('$0.00');
    expect(store.formatCurrency(null)).toBe('$0.00');
    expect(store.formatCurrency(undefined)).toBe('$0.00');
    expect(store.formatCurrency(NaN)).toBe('$0.00');
  });

  it('formats percentage correctly', () => {
    const store = useStockStore();
    // NOTA: Si tu función formatPercentage divide por 100, el mock de datos debe ser el valor decimal (e.g., 0.02 para 2%)
    // O bien, el test debe pasar el valor entero y esperar la conversión (e.g., 2 para 2.00%)
    expect(store.formatPercentage(0.02)).toBe('2.00%'); // Si la función espera decimal
    expect(store.formatPercentage(0)).toBe('0.00%');
    expect(store.formatPercentage(null)).toBe('0.00%');
    expect(store.formatPercentage(undefined)).toBe('0.00%');
    expect(store.formatPercentage(NaN)).toBe('0.00%');
  });

  it('formats market cap correctly', () => {
    const store = useStockStore();
    // 12477.23 millones = 12.47723 billones
    expect(store.formatMarketCap(12477.23)).toBe('$12.48B');
    expect(store.formatMarketCap(150.5)).toBe('$150.50M');
    expect(store.formatMarketCap(0)).toBe('$0.00M');
    expect(store.formatMarketCap(null)).toBe('$0.00M');
    expect(store.formatMarketCap(undefined)).toBe('$0.00M');
    expect(store.formatMarketCap(NaN)).toBe('$0.00M');
  });

  it('formats date correctly', () => {
    const store = useStockStore();
    expect(store.formatDate('2025-08-18T16:25:13Z')).toBe('Aug 18, 2025');
    expect(store.formatDate(null)).toBe('N/A');
    expect(store.formatDate(undefined)).toBe('N/A');
    expect(store.formatDate('')).toBe('N/A');
    expect(store.formatDate('0001-01-01T00:00:00Z')).toBe('N/A');
    expect(store.formatDate('invalid date string')).toBe('N/A');
  });
});