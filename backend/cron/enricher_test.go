package enricher

import (
	"database/sql"
	"testing"

	"github.com/jannin2/stock-app/backend/models"
)

func TestCalculateRecommendationScore(t *testing.T) {
	tests := []struct {
		name          string
		stock         models.Stock
		expectedScore float64
	}{
		{
			name: "Buy action, target met",
			stock: models.Stock{
				Action:       "Buy",
				CurrentPrice: 100.0,
				TargetTo:     models.NullFloat64{sql.NullFloat64{Float64: 120.0, Valid: true}}, // 120 is > 100 * 1.1 (110)
			},
			expectedScore: 8.0, // 5 (Buy) + 3 (Target)
		},
		{
			name: "Strong Buy action, target met",
			stock: models.Stock{
				Action:       "Strong Buy",
				CurrentPrice: 50.0,
				TargetTo:     models.NullFloat64{sql.NullFloat64{Float64: 60.0, Valid: true}}, // 60 is > 50 * 1.1 (55)
			},
			expectedScore: 8.0, // 5 (Strong Buy) + 3 (Target)
		},
		{
			name: "Hold action, target met",
			stock: models.Stock{
				Action:       "Hold",
				CurrentPrice: 100.0,
				TargetTo:     models.NullFloat64{sql.NullFloat64{Float64: 120.0, Valid: true}},
			},
			expectedScore: 3.0, // 0 (Hold) + 3 (Target)
		},
		{
			name: "Buy action, target not met (too low)",
			stock: models.Stock{
				Action:       "Buy",
				CurrentPrice: 100.0,
				TargetTo:     models.NullFloat64{sql.NullFloat64{Float64: 105.0, Valid: true}}, // 105 is NOT > 100 * 1.1 (110)
			},
			expectedScore: 5.0, // 5 (Buy) + 0 (Target)
		},
		{
			name: "Buy action, target invalid (null)",
			stock: models.Stock{
				Action:       "Buy",
				CurrentPrice: 100.0,
				TargetTo:     models.NullFloat64{sql.NullFloat64{Valid: false}},
			},
			expectedScore: 5.0, // 5 (Buy) + 0 (Target)
		},
		{
			name: "Buy action, target 0",
			stock: models.Stock{
				Action:       "Buy",
				CurrentPrice: 100.0,
				TargetTo:     models.NullFloat64{sql.NullFloat64{Float64: 0.0, Valid: true}},
			},
			expectedScore: 5.0, // 5 (Buy) + 0 (Target)
		},
		{
			name: "Neutral action, target not met",
			stock: models.Stock{
				Action:       "Neutral",
				CurrentPrice: 100.0,
				TargetTo:     models.NullFloat64{sql.NullFloat64{Float64: 105.0, Valid: true}},
			},
			expectedScore: 0.0, // 0 (Neutral) + 0 (Target)
		},
		{
			name: "Sell action, target met (shouldn't add points for Sell action)",
			stock: models.Stock{
				Action:       "Sell",
				CurrentPrice: 100.0,
				TargetTo:     models.NullFloat64{sql.NullFloat64{Float64: 120.0, Valid: true}},
			},
			expectedScore: 3.0, // 0 (Sell) + 3 (Target)
		},
		{
			name: "All conditions not met (example from your data)",
			stock: models.Stock{
				Action:       "target lowered by", // Example's action
				CurrentPrice: 122.06,
				TargetTo:     models.NullFloat64{sql.NullFloat64{Float64: 0.0, Valid: true}}, // Example's target to
			},
			expectedScore: 0.0,
		},
		{
			name: "Current price is zero (target logic skipped)",
			stock: models.Stock{
				Action:       "Buy",
				CurrentPrice: 0.0, // CurrentPrice is 0, so target condition is skipped
				TargetTo:     models.NullFloat64{sql.NullFloat64{Float64: 10.0, Valid: true}},
			},
			expectedScore: 5.0, // Only Buy action contributes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualScore := CalculateRecommendationScore(tt.stock)
			if actualScore != tt.expectedScore {
				t.Errorf("For test '%s': Expected score %.2f, got %.2f", tt.name, tt.expectedScore, actualScore)
			}
		})
	}
}
