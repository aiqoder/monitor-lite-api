package jwt

import (
	"fmt"
	"testing"
)

func TestGenPassword(t *testing.T) {
	fmt.Println(GeneratePassword("Thai1999@@", "TSGMM"))
}
