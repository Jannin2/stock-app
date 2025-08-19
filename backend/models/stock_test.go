package models

import (
	"encoding/json"
	"testing"
)

// ... (TestNullFloat64_MarshalJSON remains the same) ...

// Test for NullFloat64 UnmarshalJSON
func TestNullFloat64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected NullFloat64
		wantErr  bool
	}{
		{
			name:     "Valid float string",
			input:    []byte("123.45"),
			expected: NewNullFloat64(123.45),
			wantErr:  false,
		},
		{
			name:     "Null literal",
			input:    []byte("null"),
			expected: NullFloat64{}, // Valid: false, Float64: 0.0
			wantErr:  false,
		},
		{
			name:     "Empty string",
			input:    []byte(`""`), // JSON empty string
			expected: NullFloat64{},
			wantErr:  false,
		},
		{
			name:     "String NA",     // Should be treated as null
			input:    []byte(`"N/A"`), // This is a JSON string "N/A"
			expected: NullFloat64{},
			wantErr:  false, // ✅ FIX: Should no longer cause an error, but result in null.
		},
		{
			name:     "Invalid format",
			input:    []byte(`"abc"`),
			expected: NullFloat64{},
			wantErr:  true, // Expect an error for non-numeric string that's not "null", "", or "N/A"
		},
		{
			name:     "Zero float string",
			input:    []byte("0"),
			expected: NewNullFloat64(0.0),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nf NullFloat64
			err := json.Unmarshal(tt.input, &nf)

			// The logic here needs to be careful: if err is nil and tt.wantErr is true, that's a problem.
			// If err is not nil and tt.wantErr is false, that's also a problem.
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalJSON() error = %v, wantErr %v for input %s", err, tt.wantErr, string(tt.input))
			}
			if err != nil && tt.wantErr {
				// If error is expected, we don't check values.
				return
			}

			if nf.Valid != tt.expected.Valid {
				t.Errorf("Expected Valid to be %t, got %t for input %s", tt.expected.Valid, nf.Valid, string(tt.input))
			}
			if nf.Valid && nf.Float64 != tt.expected.Float64 {
				t.Errorf("Expected Float64 %f, got %f for input %s", tt.expected.Float64, nf.Float64, string(tt.input))
			}
		})
	}
}
