
export interface Stock {
  id: string;
  ticker: string;
  company: string;
  brokerage: string;
  action: string;
  rating_from: string | null; 
  rating_to: string | null;   
  target_from: number | null; 
  target_to: number | null;   
  current_price: number;      
  pe_ratio: number | null;    
  market_capitalization: number; 
  alpha: number | null;       
  latest_trading_day: string; 
  recommendation_score: number | null; 
  created_at: string;         
  updated_at: string;         
}