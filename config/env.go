package config

import "github.com/rickferrdev/dotenv"

type Envr struct {
	GongoPort      string `env:"GONGO_PORT"`
	GongoMode      string `env:"GONGO_MODE"`
	GongoMongoURI  string `env:"GONGO_MONGO_URI"`
	GongoJWTSecret string `env:"GONGO_JWT_SECRET"`
}

func NewEnvr() (*Envr, error) {
	var envr Envr
	dotenv.Collect()

	if err := dotenv.Unmarshal(&envr); err != nil {
		return nil, err
	}

	return &envr, nil
}
