package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type CfgType int

const (
	OCR CfgType = iota
	SSO
)

type ocrCfg struct {
	Host         string `env:"OCR_HOST" env-default:"ocr"`
	Port         string `env:"OCR_PORT" env-default:"50051"`
	Timeout      string `env:"OCR_TIMEOUT" env-default:"1m"`
	RetriesCount string `env:"OCR_RETRIES" env-default:"10"`
}

type ssoCfg struct {
	Host         string `env:"SSO_HOST" env-default:"sso"`
	Port         string `env:"SSO_PORT" env-default:"9081"`
	Timeout      string `env:"SSO_TIMEOUT" env-default:"1m"`
	RetriesCount string `env:"SSO_RETIRES" env-default:"10"`
}

type Config struct {
	Addr         string
	Timeout      time.Duration
	RetriesCount int
}

type ConfigAttr struct {
	CfgType CfgType
}

func FetchConfig(attr ConfigAttr) (*Config, error) {
	const op = "config.FetchConfig"

	switch attr.CfgType {

	case OCR:
		var ocr ocrCfg
		if err := cleanenv.ReadEnv(&ocr); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		addr := fmt.Sprintf("%s:%s", ocr.Host, ocr.Port)
		timeout, err := time.ParseDuration(ocr.Timeout)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		retries, err := strconv.Atoi(ocr.RetriesCount)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		return &Config{
			Addr:         addr,
			Timeout:      timeout,
			RetriesCount: retries,
		}, nil

	case SSO:
		var sso ssoCfg
		if err := cleanenv.ReadEnv(&sso); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		addr := fmt.Sprintf("%s:%s", sso.Host, sso.Port)
		timeout, err := time.ParseDuration(sso.Timeout)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		retries, err := strconv.Atoi(sso.RetriesCount)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		return &Config{
			Addr:         addr,
			Timeout:      timeout,
			RetriesCount: retries,
		}, nil

	default:
		return nil, fmt.Errorf("%s:%s", op, "config attrs are empty")
	}
}
