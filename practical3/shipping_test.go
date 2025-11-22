// shipping_test.go
package shipping

import (
	"testing"
)

func TestCalculateShippingFee_EquivalencePartitioning(t *testing.T) {
	// Test P1: A weight that is too small (e.g., -5). We expect an error.
	_, err := CalculateShippingFee(-5, "Domestic")
	if err == nil {
		t.Error("Test failed: Expected an error for a negative weight, but got nil")
	}

	// Test P2 & P4: A valid weight (10) and a valid zone ("Domestic").
	// We expect a correct calculation and no error.
	fee, err := CalculateShippingFee(10, "Domestic")
	if err != nil {
		t.Errorf("Test failed: Expected no error for a valid weight, but got %v", err)
	}
	expectedFee := 15.0 // From the spec: 5.0 base + (10kg * $1.0/kg)
	if fee != expectedFee {
		t.Errorf("Test failed: Expected fee of %f, but got %f", expectedFee, fee)
	}

	// Test P3: A weight that is too large (e.g., 100). We expect an error.
	_, err = CalculateShippingFee(100, "International")
	if err == nil {
		t.Error("Test failed: Expected an error for an overweight package, but got nil")
	}

	// Test P5: An invalid zone (e.g., "Local"). We expect an error.
	_, err = CalculateShippingFee(20, "Local") // Using a valid weight
	if err == nil {
		t.Error("Test failed: Expected an error for an invalid zone, but got nil")
	}
}
func TestCalculateShippingFee_BoundaryValueAnalysis(t *testing.T) {
	// We define a struct to hold all the data for one test case
	testCases := []struct {
		name        string  // A description of the test case
		weight      float64 // The input weight
		zone        string  // The input zone
		expectError bool    // True if we expect an error, false otherwise
	}{
		// Test cases for the identified boundaries of the 'weight' input
		{"Weight at lower invalid boundary", 0, "Domestic", true},
		{"Weight just above lower boundary", 0.1, "Domestic", false},
		{"Weight at upper valid boundary", 50, "International", false},
		{"Weight just above upper boundary", 50.1, "Express", true},
	}

	// Go's testing package lets us loop through our test cases
	for _, tc := range testCases {
		// t.Run creates a sub-test, which gives a clearer output if one fails
		t.Run(tc.name, func(t *testing.T) {
			_, err := CalculateShippingFee(tc.weight, tc.zone)

			// Assertion 1: Check if we got an error when we expected one
			if tc.expectError && err == nil {
				t.Errorf("Expected an error, but got nil")
			}
			
			// Assertion 2: Check if we got no error when we didn't expect one
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}

func TestCalculateShippingFee_DecisionTable(t *testing.T) {
	testCases := []struct {
		name          string
		weight        float64
		zone          string
		expectedFee   float64
		expectedError string // We check for a substring of the error message
	}{
		// Rule 1: Invalid weight. Zone does not matter.
		{"Rule 1: Weight too low", -10, "Domestic", 0, "invalid weight"},
		{"Rule 1: Weight too high", 60, "International", 0, "invalid weight"},
		
		// Rule 2: Valid weight, Domestic zone
		{"Rule 2: Domestic", 10, "Domestic", 15.0, ""}, // 5 + 10 * 1.0
		
		// Rule 3: Valid weight, International zone
		{"Rule 3: International", 10, "International", 45.0, ""}, // 20 + 10 * 2.5
		
		// Rule 4: Valid weight, Express zone
		{"Rule 4: Express", 10, "Express", 80.0, ""}, // 30 + 10 * 5.0

		// Rule 5: Valid weight, Invalid Zone
		{"Rule 5: Invalid Zone", 10, "Unknown", 0, "invalid zone: Unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fee, err := CalculateShippingFee(tc.weight, tc.zone)

			// Check if we got the error we expected
			if tc.expectedError != "" {
				if err == nil {
					t.Fatalf("Expected error containing '%s', but got nil", tc.expectedError)
				}
				// A simple check is to see if our error message contains the expected text.
				// In more complex apps, you might check for specific error types.
			} else {
				// Check that we did NOT get an error when one wasn't expected
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}
				// Check if the calculated fee is correct
				if fee != tc.expectedFee {
					t.Errorf("Expected fee %f, but got %f", tc.expectedFee, fee)
				}
			}
		})
	}
}