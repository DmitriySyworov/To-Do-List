package generateRand_test

import (
	"log"
	"testing"
	"to-do-list/app/pkg/generateRand"
)

const (
	lengthUserId = 11
	lengthSessId = 6
)

func TestGenerateNumbers(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Println(generateRand.GenerateNumbers(lengthUserId))
	}
}
func TestGenerateStr(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Println(generateRand.GenerateStr(lengthSessId))
	}
}
