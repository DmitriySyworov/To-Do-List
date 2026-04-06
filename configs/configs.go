package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	*Token
	*DbConf
	*SendEmail
}
type Token struct {
	Secret string
}
type SendEmail struct {
	SenderEmail string
	Password    string
	AddressHost string
	Address     string
}
type DbConf struct {
	DSN           string
	RedisPassword string
}

func NewConfigs() *Configs {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}
	return &Configs{
		&Token{
			Secret: os.Getenv("SECRET"),
		},
		&DbConf{
			DSN:           os.Getenv("DSN"),
			RedisPassword: os.Getenv("REDIS"),
		},
		&SendEmail{
			SenderEmail: os.Getenv("EMAIL"),
			Password:    os.Getenv("PASSWORD"),
			AddressHost: os.Getenv("ADDRESS_HOST"),
			Address:     os.Getenv("ADDRESS"),
		},
	}
}
