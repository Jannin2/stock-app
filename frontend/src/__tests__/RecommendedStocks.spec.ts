// src/__tests__/RecommendedStocks.spec.ts
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia, type TestingPinia } from '@pinia/testing'; // Import TestingPinia
import { setActivePinia } from 'pinia'; // Import setActivePinia
import RecommendedStocks from '../components/RecommendedStocks.vue';
import { useStockStore } from '../stores/stocks';

const mockRecommendedStocks = [
  // ... (your mock data remains the same) ...
  {
    id: '3ddc91f7-2010-4d99-bb93-93c9cff0bc5e',
    ticker: 'LAMR',
    company: 'Lamar Advertising',
    brokerage: 'Wells Fargo & Company',
    action: 'target lowered by',
    rating_from: null,
    rating_to: null,
    target_from: null,
    target_to: null,
    current_price: 122.48,
    pe_ratio: 28.4,
    dividend_yield: 0.00,
    market_capitalization: 12477.23,
    alpha: null,
    latest_trading_day: '2025-08-18T16:25:13Z',
    recommendation_score: 0,
    created_at: '2025-08-18T15:29:12.132415Z',
    updated_at: '2025-08-18T16:27:35.640371Z',
  },
  {
    id: '44812b78-6558-4b64-b20a-b842ac3a3921',
    ticker: 'VYGR',
    company: 'Voyager Therapeutics',
    brokerage: 'Wedbush',
    action: 'buy',
    rating_from: 'Neutral',
    rating_to: 'Buy',
    target_from: 10,
    target_to: 15,
    current_price: 3.86,
    pe_ratio: 0,
    dividend_yield: 0.05,
    market_capitalization: 178.18,
    alpha: 0.1,
    latest_trading_day: '2025-08-18T16:23:49Z',
    recommendation_score: 7.5,
    created_at: '2025-08-18T15:29:12.132415Z',
    updated_at: '2025-08-18T16:27:35.640371Z',
  },
];

describe('RecommendedStocks.vue', () => {
  let stockStore: ReturnType<typeof useStockStore>;
  let pinia: TestingPinia; // Declare pinia here

  beforeEach(() => {
    // 1. Reset all Vitest mocks/spies before each test
    vi.clearAllMocks();

    // 2. Create and set an active Pinia instance for the current test
    // This allows useStockStore() to be called immediately afterwards
    pinia = createTestingPinia({
      createSpy: vi.fn,
      // stubActions: true is the default and preferred for isolated component testing.
    });
    setActivePinia(pinia); // Make this Pinia instance the active one for the current test

    // 3. Get the stockStore instance associated with this active Pinia
    stockStore = useStockStore();

    // 4. Now, mock the formatting functions on this specific stockStore instance
    (stockStore.formatCurrency as ReturnType<typeof vi.fn>).mockImplementation((val) => {
      const numValue = (val === null || typeof val === 'undefined' || isNaN(val as number)) ? 0 : val as number;
      return `$${numValue.toFixed(2)}`;
    });
    (stockStore.formatPercentage as ReturnType<typeof vi.fn>).mockImplementation((val) => {
      const numValue = (val === null || typeof val === 'undefined' || isNaN(val as number)) ? 0 : val as number;
      return `${(numValue * 100).toFixed(2)}%`;
    });
    (stockStore.formatMarketCap as ReturnType<typeof vi.fn>).mockImplementation((val) => {
      const numValue = (val === null || typeof val === 'undefined' || isNaN(val as number)) ? 0 : val as number;
      if (numValue === 0) return '$0.00M';
      return numValue >= 1000 ? `$${(numValue / 1000).toFixed(2)}B` : `$${numValue.toFixed(2)}M`;
    });
    (stockStore.formatDate as ReturnType<typeof vi.fn>).mockImplementation((dateString) => {
      if (!dateString || dateString === "0001-01-01T00:00:00Z" || dateString.startsWith("0000-")) {
        return 'N/A';
      }
      return 'Mocked Date';
    });
  });

  // Helper function to mount the component with specific Pinia store state
  const mountComponentWithStoreState = (initialState: any) => {
    // We create the testing pinia inside the global plugin for the component
    // but the store instance obtained above (stockStore) is already tied to the active pinia.
    const wrapper = mount(RecommendedStocks, {
      global: {
        plugins: [
          pinia // Use the 'pinia' instance already set up in beforeEach
        ],
      },
    });

    // Update the store's state *after* it's been initialized by the wrapper mount
    // This is the most reliable way to set reactive state for the component
    Object.assign(stockStore, initialState);

    return wrapper;
  };

  // Test 1: renders loading message initially
  it('renders loading message initially', async () => {
    const wrapper = mountComponentWithStoreState({
      loading: true,
      error: null,
      recommendedStocks: [],
    });

    await wrapper.vm.$nextTick();
    expect(wrapper.text()).toContain('Loading recommended stocks... ⏳');
    // Assuming your component has a div/span with class="loading-message" for this text
    expect(wrapper.find('.loading-message').exists()).toBe(true);
  });

  // Test 2: renders error message if fetching fails
  it('renders error message if fetching fails', async () => {
    const wrapper = mountComponentWithStoreState({
      loading: false,
      error: 'Failed to fetch data',
      recommendedStocks: [],
    });

    await wrapper.vm.$nextTick();
    expect(wrapper.text()).toContain('Error: Failed to fetch data 🔴');
    // Assuming your component has a div/span with class="error-message" for this text
    expect(wrapper.find('.error-message').exists()).toBe(true);
  });

  // Test 3: renders no data message if no recommended stocks are found
  it('renders no data message if no recommended stocks are found', async () => {
    const wrapper = mountComponentWithStoreState({
      loading: false,
      error: null,
      recommendedStocks: [], // Empty array for no data
    });

    await wrapper.vm.$nextTick();
    expect(wrapper.text()).toContain('No recommended stocks found. 🤷‍♂️');
    // Assuming your component has a div/span with class="no-data-message" for this text
    expect(wrapper.find('.no-data-message').exists()).toBe(true);
  });

  // Test 4: renders recommended stock cards when data is available
  it('renders recommended stock cards when data is available', async () => {
    const wrapper = mountComponentWithStoreState({
      loading: false,
      error: null,
      recommendedStocks: mockRecommendedStocks,
    });
    await wrapper.vm.$nextTick(); // Wait for Vue to process data and re-render

    expect(wrapper.find('h2').text()).toBe('Recommended Stocks');

    const stockCards = wrapper.findAll('.stock-card');
    expect(stockCards.length).toBe(mockRecommendedStocks.length);

    // Verify content of the first card (LAMR)
    expect(stockCards[0].text()).toContain('LAMR (Lamar Advertising)');
    expect(stockCards[0].text()).toContain('Brokerage: Wells Fargo & Company');
    expect(stockCards[0].text()).toContain('Action: target lowered by');
    expect(stockCards[0].text()).toContain('Rating: N/A to N/A');
    expect(stockCards[0].text()).toContain('Target Price: $0.00 - $0.00');
    expect(stockCards[0].text()).toContain('Current Price: $122.48');
    expect(stockCards[0].text()).toContain('PE Ratio: 28.40');
    expect(stockCards[0].text()).toContain('Dividend Yield: 0.00%');
    expect(stockCards[0].text()).toContain('Market Cap: $12.48B');
    expect(stockCards[0].text()).toContain('Alpha: 0.00%');
    expect(stockCards[0].text()).toContain('Latest Trading Day: Mocked Date');
    expect(stockCards[0].text()).toContain('Recommendation Score: 0.00');

    // Verify content of the second card (VYGR)
    expect(stockCards[1].text()).toContain('VYGR (Voyager Therapeutics)');
    expect(stockCards[1].text()).toContain('Action: buy');
    expect(stockCards[1].text()).toContain('Rating: Neutral to Buy');
    expect(stockCards[1].text()).toContain('Target Price: $10.00 - $15.00');
    expect(stockCards[1].text()).toContain('Current Price: $3.86');
    expect(stockCards[1].text()).toContain('PE Ratio: 0.00');
    expect(stockCards[1].text()).toContain('Dividend Yield: 5.00%');
    expect(stockCards[1].text()).toContain('Market Cap: $178.18M');
    expect(stockCards[1].text()).toContain('Alpha: 10.00%');
    expect(stockCards[1].text()).toContain('Latest Trading Day: Mocked Date');
    expect(stockCards[1].text()).toContain('Recommendation Score: 7.50');

    // Verify card classes
    expect(stockCards[0].classes()).toContain('card-neutral');
    expect(stockCards[1].classes()).toContain('card-buy');
  });

  // Test 5: calls fetchRecommendedStocks on mount
  it('calls fetchRecommendedStocks on mount', async () => {
    // Initial state where fetchRecommendedStocks *should* be called by onMounted
    const wrapper = mountComponentWithStoreState({
      loading: false,
      error: null,
      recommendedStocks: [], // Start empty to trigger fetch
    });

    // Mock the action *after* getting the store instance, but before it's called by the component
    (stockStore.fetchRecommendedStocks as ReturnType<typeof vi.fn>).mockResolvedValue(undefined);

    // Wait for the onMounted hook to run and trigger the action
    await wrapper.vm.$nextTick();

    // Use vi.waitFor to ensure the async mock function has been called
    await vi.waitFor(() => {
      expect(stockStore.fetchRecommendedStocks).toHaveBeenCalledTimes(1);
    }, { timeout: 2000 }); // Increased timeout for robustness
  });
});