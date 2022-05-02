package referral_code

import (
	"crypto/rand"
	"fmt"
)

func RandomString() string {
	n := 3
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return (s)
}
