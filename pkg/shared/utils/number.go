package utils

import (
	"fmt"
	"math"
	"strings"
)

func ToFixed(amount float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return float64(int(amount*multiplier)) / multiplier
}

func FormatIndonesianNumber(amount float64) string {
	amountStr := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(amountStr, ".")

	intPart := parts[0]
	decPart := parts[1]

	var result strings.Builder
	intLen := len(intPart)

	for i, digit := range intPart {
		if i > 0 && (intLen-i)%3 == 0 {
			result.WriteString(".")
		}
		result.WriteRune(digit)
	}

	return result.String() + "," + decPart
}
