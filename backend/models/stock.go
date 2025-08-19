package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// NullFloat64 is a wrapper around sql.NullFloat64 that handles JSON marshaling/unmarshaling.
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON implements the json.Marshaler interface.
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil // Return JSON 'null' for invalid values
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	s := strings.TrimSpace(string(data))
	// ✅ FIX: Handle quoted "N/A" and other non-numeric strings that should be null
	if s == "null" || s == "" || s == `""` || strings.ToLower(s) == `"n/a"` {
		nf.Valid = false
		return nil
	}

	var f float64
	// Attempt to unmarshal directly. json.Unmarshal will handle "123.45" vs 123.45
	if err := json.Unmarshal(data, &f); err != nil {
		nf.Valid = false
		return fmt.Errorf("json: cannot unmarshal %s into Go value of type float64: %w", string(data), err)
	}
	nf.Float64 = f
	nf.Valid = true
	return nil
}

// NewNullFloat64 is a helper function to create a valid NullFloat64.
func NewNullFloat64(f float64) NullFloat64 {
	return NullFloat64{sql.NullFloat64{Float64: f, Valid: true}}
}

// NullTime is a wrapper around sql.NullTime that handles JSON marshaling/unmarshaling.
type NullTime struct {
	sql.NullTime
}

// MarshalJSON implements the json.Marshaler interface.
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil // Return JSON 'null' for invalid time
	}
	return json.Marshal(nt.Time.Format(time.RFC3339))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	s := strings.TrimSpace(string(data))
	if s == "null" || s == "" || s == `""` { // Handle null or empty string
		nt.Valid = false
		return nil
	}

	// Remove quotes if present, as time.Parse expects unquoted string
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	// Try parsing with RFC3339 first (standard for JSON dates)
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		nt.Time = t
		nt.Valid = true
		return nil
	}

	// If RFC3339 fails, try a common "YYYY-MM-DD" format as a fallback
	t, err = time.Parse("2006-01-02", s)
	if err == nil {
		nt.Time = t
		nt.Valid = true
		return nil
	}

	nt.Valid = false
	return fmt.Errorf("could not parse time %q, expected RFC3339 or 2006-01-02 format: %w", s, err)
}

// NewNullTime is a helper function to create a valid NullTime.
func NewNullTime(t time.Time) NullTime {
	return NullTime{sql.NullTime{Time: t, Valid: true}}
}

// Stock represents a stock entry with detailed financial metrics.
type Stock struct {
	ID                   uuid.UUID   `json:"id"`
	Ticker               string      `json:"ticker"`
	Company              string      `json:"company"`
	Brokerage            string      `json:"brokerage"`
	Action               string      `json:"action"`      // E.g., Buy, Sell, Hold
	RatingFrom           string      `json:"rating_from"` // Previous rating
	RatingTo             string      `json:"rating_to"`   // New rating
	TargetFrom           NullFloat64 `json:"target_from"` // Previous target price
	TargetTo             NullFloat64 `json:"target_to"`   // New target price
	CurrentPrice         float64     `json:"current_price"`
	PERatio              NullFloat64 `json:"pe_ratio"`
	DividendYield        NullFloat64 `json:"dividend_yield"`
	MarketCapitalization NullFloat64 `json:"market_capitalization"`
	Alpha                NullFloat64 `json:"alpha"`              // Alpha value
	LatestTradingDay     NullTime    `json:"latest_trading_day"` // Date of the latest trading data
	RecommendationScore  NullFloat64 `json:"recommendation_score"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}
