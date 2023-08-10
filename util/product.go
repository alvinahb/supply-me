package util

// Constants for all supported product categories
const (
	meat      = "Viande"
	vegetable = "Légume"
	fruit     = "Fruit"
)

// IsSupportedProductCategory retruns true is the category is supported
func IsSupportedProductCategory(category string) bool {
	switch category {
	case meat, vegetable, fruit:
		return true
	}

	return false
}
