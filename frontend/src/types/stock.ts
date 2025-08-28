
// Define y exporta la interfaz 'Stock'. Esta interfaz sirve como un contrato o plantilla
// para asegurar que cualquier objeto que represente una acción en tu aplicación
// tenga una estructura de datos consistente.
export interface Stock {
  // Un identificador único para la acción.
  id: string;
  // El símbolo de cotización de la acción (por ejemplo, "AAPL" para Apple).
  ticker: string;
  // El nombre completo de la compañía.
  company: string;
  // La firma de corretaje que emitió la recomendación.
  brokerage: string;
  // La recomendación de la acción (por ejemplo, "Buy", "Hold", "Sell").
  action: string;
  // La calificación inicial, puede ser una cadena de texto o null si no está disponible.
  rating_from: string | null; 
  // La calificación final, también puede ser null.
  rating_to: string | null;   
  // El precio objetivo inicial, puede ser un número o null.
  target_from: number | null; 
  // El precio objetivo final, también puede ser null.
  target_to: number | null;   
  // El precio actual de la acción.
  current_price: number;      
  // El ratio Precio/Ganancias, puede ser un número o null.
  pe_ratio: number | null;    
  // La capitalización de mercado de la compañía.
  market_capitalization: number; 
  // El valor Alpha, un indicador de rendimiento, puede ser un número o null.
  alpha: number | null;       
  // La fecha del último día de negociación en formato de cadena de texto.
  latest_trading_day: string; 
  // La puntuación numérica de la recomendación, puede ser un número o null.
  recommendation_score: number | null; 
  // La fecha de creación del registro.
  created_at: string;         
  // La fecha de la última actualización del registro.
  updated_at: string;  
  // El rendimiento por dividendo, una propiedad recién añadida.
  dividend_yield: number;
}
