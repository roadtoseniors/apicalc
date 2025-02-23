package config

import (
	"fmt"
	"os"
	"time"
)

const errMessageFmt = "The %s environment variable is not set or has an incorrect value."
type TimeConfig struct {
	Add time.Duration
	Sub time.Duration
	Mul time.Duration
	Div time.Duration
}

func NewConfigOrch() (*TimeConfig, error){

	at, err := time.ParseDuration(os.Getenv("TIME_ADDITION_MS") + "ms")
	if err != nil || at < 0{
		return nil, fmt.Errorf(errMessageFmt, "TIME_ADDITION_MS")
	}

	st, err := time.ParseDuration(os.Getenv("TIME_SUBTRACTION_MS") + "ms")
	if err != nil || st < 0{
		return nil, fmt.Errorf(errMessageFmt, "TIME_SUBTRACTION_MS")
	}

	mt, err := time.ParseDuration(os.Getenv("TIME_MULTIPLICATIONS_MS") + "ms")
	if err != nil || mt < 0 {
		return nil,  fmt.Errorf(errMessageFmt, "TIME_MULTIPLICATIONS_MS")
	}

	dt, err := time.ParseDuration(os.Getenv("TIME_DIVISIONS_MS") + "ms")
	if err != nil || dt < 0{
		return nil, fmt.Errorf(errMessageFmt, "TIME_DIVISIONS_MS")
	}

	orchcfg := TimeConfig{
		Add: at,
		Sub: st,
		Mul: mt,
		Div: dt, 
	}

	return &orchcfg, nil
}