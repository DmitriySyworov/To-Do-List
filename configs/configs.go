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
	RedisHost     string
}

func NewConfigs() *Configs {
	if os.Getenv("GO_TEST") == "test" {
		errEnv := godotenv.Load(".env.test")
		if errEnv != nil {
			panic(errEnv)
		}
		return &Configs{
			&Token{
				Secret: os.Getenv("SECRET_TEST"),
			},
			&DbConf{
				DSN:           os.Getenv("DSN_TEST"),
				RedisPassword: os.Getenv("REDIS_TEST"),
				RedisHost:     os.Getenv("REDIS_PORT_TEST"),
			},
			&SendEmail{
				SenderEmail: os.Getenv("EMAIL_TEST"),
				Password:    os.Getenv("PASSWORD_TEST"),
				AddressHost: os.Getenv("ADDRESS_HOST_TEST"),
				Address:     os.Getenv("ADDRESS_TEST"),
			},
		}
	} else {
		errEnv := godotenv.Load(".env")
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
				RedisHost:     os.Getenv("REDIS_PORT"),
			},
			&SendEmail{
				SenderEmail: os.Getenv("EMAIL"),
				Password:    os.Getenv("PASSWORD"),
				AddressHost: os.Getenv("ADDRESS_HOST"),
				Address:     os.Getenv("ADDRESS"),
			},
		}
	}
}
