package jwt

import (
	"testing"
	"to-do-list/app/internal/test"
)

const (
	Secret = "skDZeVzf7fsPo_GSlSJNuVbxbnWEOZEF"
	Token  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozMzAxMzkxMDgxMX0.ru_3WRffT2xuwBrnGGZ26kTaVyOoLbMLYDiJ_k4JKos"
)

func TestCreateJWT(t *testing.T) {
	j := NewJWT(Secret)
	token, errToken := j.CreateJWT(test.UserIDTest)
	if errToken != nil {
		t.Fatal(errToken)
	}
	t.Log(token)
}
func TestParseJWT(t *testing.T) {
	j := NewJWT(Secret)
	id, errParse := j.ParseJWt(Token)
	if errParse != nil {
		t.Fatal(errParse)
	}
	t.Log(uint(id))
}
