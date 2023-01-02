package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PiAwareConfig struct {
		AircraftUrl string `yaml:"aircraft-url" env:"FA_AC_URL" env-default:"http://piaware:80/skyaware/data/aircraft.json"`
		ReceiverUrl string `yaml:"receiver-url" env:"FA_RCV_URL" env-defualt:"http://piaware:80/skyaware/data/receiver.json"`
	} `yaml:"piaware"`
	Display struct {
		Delay              int    `yaml:"delay" env:"FA_DISPLAY_DELAY" env-default:"5"`
		TransitionAltitude int    `yaml:"transAlt" env:"FA_TRANSALT" env-default:"19000"`
		UseFlightLevel     bool   `yaml:"useFL" env:"FA_USEFL" env-default:"true"`
		AltitudeUnits      string `yaml:"altUnits" env:"FA_ALTUOM" env-default:"FT"`
		SpeedUnits         string `yaml:"speedUnits" env:"FA_SPEEDUOM" env-default:"KN"`
		DistanceUnits      string `yaml:"distanceUnits" env:"FA_DISTUOM" env-default:"NM"`
	} `yaml:"display"`
}

var ServiceError = errors.New("Could not invoke receiver service on piaware")

func GetConfig(configPath string, config *Config) error {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil
	}

	return nil
}
