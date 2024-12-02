package config

import (
	"errors"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mkadit/go-toybox/internal/models"
	"github.com/spf13/viper"
)

var (
	_, b, _, _        = runtime.Caller(0)
	ProjectRootPath   = filepath.Join(filepath.Dir(b), "../..")
	TimeLoc, _        = time.LoadLocation(AreaLocation)
	AppSecret         []byte
	JwtExpireDuration = time.Duration(2) * time.Hour
	Routes            []string
	AppConfig         models.Configuration
)

const (
	AppName      = "Toybox"
	TimeFormat   = "2006-01-02"
	AreaLocation = "Asia/Jakarta"
)

var (
	TLV_FORMAT_TIMESTAMP         = "20060102150405"
	FORMAT_TRANSMISSION_DATETIME = "0102150405"
	FORMAT_DATE                  = "0102"
	FORMAT_TIME                  = "150405"
	FORMAT_ISO8601               = "2006-01-02T15:04:05-07:00"
	LOCATION                     = time.FixedZone("UTC+07:00", 7*60*60)
	TIMESTAMP_LOCATION_FORMAT    = "2006-01-02 15:04:05.999 +0700"
	FORMAT_TIME_FOR_LOG          = "2006-01-02 15:04:05.000"
	ErrLoadEnv                   = errors.New("cannot load env file")
)

// LoadConfig reads configuration from file
func LoadConfig(path string) (config models.Configuration, err error) {
	viper.SetConfigType("json")
	viper.AddConfigPath(path)
	viper.SetConfigName("config")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	ProjectRootPath = filepath.Join(filepath.Dir(b), "../..")

	AppConfig = config
	return
}
