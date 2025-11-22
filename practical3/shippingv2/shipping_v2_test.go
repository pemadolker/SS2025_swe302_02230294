package shipping

import (
	"testing"
	"math"
)

func TestCalculateShippingFee_V2(t *testing.T) {

	testCases := []struct {
		name        string
		weight      float64
		zone        string
		insured     bool
		expectedFee float64
		expectError bool
	}{
		// -------------------------
		// WEIGHT PARTITIONS
		// -------------------------

		// P1: Invalid weight ≤ 0
		{"Weight <= 0 (invalid)", 0, "Domestic", false, 0, true},
		{"Weight negative (invalid)", -5, "Express", true, 0, true},

		// P2: Standard package (0 < weight ≤ 10)
		{"Standard weight", 5, "Domestic", false, 5, false},
		{"Standard weight insured", 10, "Domestic", true, 5 + (5*0.015), false},

		// P3: Heavy package (10 < weight ≤ 50)
		{"Heavy weight", 20, "International", false, 20 + 7.50, false},
		{"Heavy weight insured", 20, "International", true, (20 + 7.50) * 1.015, false},

		// P4: Weight > 50 (invalid)
		{"Weight > 50", 51, "Express", false, 0, true},

		// -------------------------
		// ZONE PARTITIONS
		// -------------------------

		// P5: Valid zones
		{"Domestic valid", 5, "Domestic", false, 5, false},
		{"International valid", 5, "International", false, 20, false},
		{"Express valid", 5, "Express", false, 30, false},

		// P6: Invalid zone
		{"Invalid zone string", 5, "Local", false, 0, true},

		// -------------------------
		// INSURANCE (P7 & P8)
		// -------------------------

		{"Uninsured standard", 5, "Domestic", false, 5, false},
		{"Insured standard", 5, "Domestic", true, 5 * 1.015, false},

		{"Uninsured heavy", 20, "Express", false, 30 + 7.50, false},
		{"Insured heavy", 20, "Express", true, (30 + 7.50) * 1.015, false},

		// -------------------------
		// BOUNDARY VALUES
		// -------------------------

		// Lower boundary: around 0
		{"Boundary: weight 0 → invalid", 0, "Domestic", false, 0, true},
		{"Boundary: weight 0.1 → valid standard", 0.1, "Domestic", false, 5, false},

		// Mid boundary: around 10
		{"Boundary: weight 10 → standard", 10, "Domestic", false, 5, false},
		{"Boundary: weight 10.1 → heavy", 10.1, "Domestic", false, 5 + 7.50, false},

		// Upper boundary: around 50
		{"Boundary: weight 50 → heavy valid", 50, "Domestic", false, 5 + 7.50, false},
		{"Boundary: weight 50.1 → invalid", 50.1, "Domestic", false, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			fee, err := CalculateShippingFee(tc.weight, tc.zone, tc.insured)

			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Use a float tolerance for small floating-point calculation differences
			if math.Abs(fee-tc.expectedFee) > 0.0001 {
				t.Errorf("Expected fee %.4f, got %.4f", tc.expectedFee, fee)
			}
		})
	}
}
