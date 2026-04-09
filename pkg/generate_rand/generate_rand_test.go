package generate_rand_test

import (
	"log"
	"testing"
	"to-do-list/app/pkg/generate_rand"
)

const (
	lengthUserId = 11
	lengthSessId = 6
)

func TestGenerateNumbers(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Println(generate_rand.GenerateNumbers(lengthUserId))
	}
}
func TestGenerateStr(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Println(generate_rand.GenerateStr(lengthSessId))
	}
}
