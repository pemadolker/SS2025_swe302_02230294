// shipping.go
package shipping

import (
	"errors"
	"fmt"
)

// CalculateShippingFee calculates the fee based on weight and zone.
func CalculateShippingFee(weight float64, zone string) (float64, error) {
	// This block directly implements Rule #1 and #4
	if weight <= 0 || weight > 50 {
		return 0, errors.New("invalid weight")
	}

	// This switch statement implements Rule #2, #3, and #5
	switch zone {
	case "Domestic":
		return 5.0 + (weight * 1.0), nil
	case "International":
		return 20.0 + (weight * 2.5), nil
	case "Express":
		return 30.0 + (weight * 5.0), nil
	default:
		// This handles any zone not explicitly listed above
		return 0, fmt.Errorf("invalid zone: %s", zone)
	}
}