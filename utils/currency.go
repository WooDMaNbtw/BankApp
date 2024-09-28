package utils

// All supported currencies
const (
	USD = "USD"
	EUR = "UER"
	CAD = "CAD"
	JPY = "JPY"
)

// IsSupportedCurrency returns true if currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, JPY:
		return true
	}
	return false
}
