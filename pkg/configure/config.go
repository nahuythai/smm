package configure

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v10"
)

var cfg *Config

type Config struct {
	Debug                          bool          `env:"DEBUG" envDefault:"false"`
	Host                           string        `env:"HOST" envDefault:"0.0.0.0"`
	Port                           string        `env:"PORT" envDefault:"8000"`
	ServerDomain                   string        `env:"SERVER_DOMAIN" envDefault:"localhost:8000"`
	SecretKey                      string        `env:"SECRET_KEY" envDefault:"!change_me!"`
	MongoRequestTimeout            time.Duration `env:"MONGO_REQUEST_TIMEOUT" envDefault:"3m"`
	TransactionDuration            time.Duration `env:"TRANSACTION_DURATION" envDefault:"15m"`
	VerifyEmailTransactionDuration time.Duration `env:"VERIFY_EMAIL_TRANSACTION_DURATION" envDefault:"24h"`
	AccessTokenDuration            time.Duration `env:"ACCESS_TOKEN_DURATION" envDefault:"24h"`
	MongoDBUrl                     string        `env:"MONGODB_URL" envDefault:"mongodb://localhost:27017"`
	MongoDBName                    string        `env:"MONGODB_NAME" envDefault:"db_smm"`
}

func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
		if err := env.Parse(cfg); err != nil {
			log.Fatal("Fail to load env! ", err)
			return nil
		}
	}
	return cfg
}

func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}
