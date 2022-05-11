package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var Configuration *viper.Viper

func Setup(log zerolog.Logger) {
	log.Info().Msg("Configs are being initialized")

	// Set Configuration File Value
	configEnv := strings.ToLower(os.Getenv("SERVER_MODE"))
	if len(configEnv) == 0 {
		configEnv = "dev"
	}
	// Initialize Configurator
	Configuration = viper.New()
	// Set Configurator Configuration
	Configuration.SetConfigName(configEnv)
	Configuration.AddConfigPath("./configs")
	Configuration.AddConfigPath("../configs")

	Configuration.SetConfigType("yml")

	// Set Configurator Environment
	Configuration.SetEnvPrefix("GEO")

	// Set Configurator to Auto Bind Configuration Variables to
	// Environment Variables
	Configuration.AutomaticEnv()

	// Set Configurator to Load Configuration File
	err := configLoadFile()
	if err != nil {
		log.Fatal().Err(err)
	}

	// uncomment for debug purpose
	//Configuration.Debug()

}

// ConfigLoadFile Function to Load Configuration from File
func configLoadFile() error {
	// Load Configuration File
	return Configuration.ReadInConfig()
}
