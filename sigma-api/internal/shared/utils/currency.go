package utils

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// FormatRupiah converts a float64 to IDR currency string
func FormatRupiah(amount float64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("Rp %.2f", amount)
}
