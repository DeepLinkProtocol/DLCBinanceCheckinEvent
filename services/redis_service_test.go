package services

import (
	"fmt"
	"testing"
)

func TestSignedToday(t *testing.T) {
	r := GetSignInDateKey("0x1234567890")
	fmt.Println(r)
}
