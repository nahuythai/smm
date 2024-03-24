package configure

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

var cfg *Config

type Config struct {
	Debug                      bool          `env:"DEBUG" envDefault:"false"`
	Host                       string        `env:"HOST" envDefault:"0.0.0.0"`
	Port                       string        `env:"PORT" envDefault:"8000"`
	ServerDomain               string        `env:"SERVER_DOMAIN" envDefault:"localhost:8000"`
	SecretKey                  string        `env:"SECRET_KEY" envDefault:"!change_me!"`
	MailEmail                  string        `env:"MAIL_EMAIL" envDefault:"smmnoreply@localhost"`
	MailPassword               string        `env:"MAIL_PASSWORD" envDefault:""`
	MailHost                   string        `env:"MAIL_HOST" envDefault:"localhost"`
	MailPort                   int           `env:"MAIL_PORT" envDefault:"1025"`
	MongoRequestTimeout        time.Duration `env:"MONGO_REQUEST_TIMEOUT" envDefault:"3m"`
	SessionDuration            time.Duration `env:"TRANSACTION_DURATION" envDefault:"15m"`
	VerifyEmailSessionDuration time.Duration `env:"VERIFY_EMAIL_TRANSACTION_DURATION" envDefault:"24h"`
	AccessTokenDuration        time.Duration `env:"ACCESS_TOKEN_DURATION" envDefault:"24h"`
	BackgroundTaskDuration     time.Duration `env:"BACKGROUND_TASK_DURATION" envDefault:"1m"`
	MongoDBUrl                 string        `env:"MONGODB_URL" envDefault:"mongodb://localhost:27017"`
	MongoDBName                string        `env:"MONGODB_NAME" envDefault:"db_smm"`
	Web2MAccessToken           string        `env:"WEB2M_ACCESS_TOKEN" envDefault:"!change_me!"`
}

func GetConfig() *Config {
	if cfg == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Fail to load env!", err)
		}
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
