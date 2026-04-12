package generate_rand

import (
	"math/rand/v2"
	"strconv"
)

func GenerateNumbers(length int) uint {
	randomer := ""
	for {
		symbol := rand.IntN(58)
		if len(randomer) != length {
			if (len(randomer) == 0 && symbol > 48) || (len(randomer) != 0 && symbol > 47) {
				randomer += string(byte(symbol))
			}
		} else {
			break
		}
	}
	resNumber, _ := strconv.Atoi(randomer)
	return uint(resNumber)
}
func GenerateStr(length int) string {
	randomer := ""
	for {
		symbol := rand.IntN(123)
		if len(randomer) != length {
			if (symbol > 64 && symbol < 91) || symbol > 96 {
				randomer += string(byte(symbol))
			}
		} else {
			return randomer
		}
	}
}
