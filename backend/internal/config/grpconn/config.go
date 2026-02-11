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
	host         string `env:"OCR_HOST" env-default:"ocr"`
	port         string `env:"OCR_PORT" env-default:"50051"`
	timeout      string `env:"OCR_TIMEOUT" env-default:"1m"`
	retriesCount string `env:"OCR_RETRIES" env-default:"10"`
}

type ssoCfg struct {
	host         string `env:"SSO_HOST" env-default:"sso"`
	port         string `env:"SSO_PORT" env-default:"9081"`
	timeout      string `env:"SSO_TIMEOUT" env-default:"1m"`
	retriesCount string `env:"SSO_RETIRES" env-default:"10"`
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

		addr := fmt.Sprintf("%s:%s", ocr.host, ocr.port)
		timeout, err := time.ParseDuration(ocr.timeout)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		retries, err := strconv.Atoi(ocr.retriesCount)
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

		addr := fmt.Sprintf("%s:%s", sso.host, sso.port)
		timeout, err := time.ParseDuration(sso.timeout)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		retries, err := strconv.Atoi(sso.retriesCount)
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
