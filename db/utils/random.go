package utils

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Create a new random source and a new random generator



var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt(min, max int64) int64 {
	return min + seededRand.Int63n(max-min+1)
}

func RandomString(length int) string {
	bytes := make([]byte, length/2)
	_, err := seededRand.Read(bytes)
	if err != nil {
		log.Printf("can not read bytes: %s", err)
		return ""
	}
	return hex.EncodeToString(bytes)
}

func RandomOwner() string {
	return RandomString(5)
}

func RandomMoney() int64 {
	return RandomInt(11, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[seededRand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(10))
}
