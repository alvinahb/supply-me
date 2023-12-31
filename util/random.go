package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomCompanyType generates a random string among companies types
func RandomCompanyType() string {
	companyTypes := []string{"Restaurant", "Fournisseur"}
	n := len(companyTypes)

	return companyTypes[rand.Intn(n)]
}

// RandomCompanyType generates a random role among roles
func RandomRole() string {
	roles := []string{"Owner", "Employee"}
	n := len(roles)

	return roles[rand.Intn(n)]
}
