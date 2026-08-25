package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func Initiator() {

	// =========================
	// 1. System environment
	// =========================
	viper.AutomaticEnv()

	if err := viper.BindEnv("DATABASE_URL", "DATABASE_URL", "POSTGRES_URL"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("DB_ENGINE", "DB_ENGINE"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("jwt_secret_key", "JWT_SECRET_KEY", "jwt_secret_key"); err != nil {
		panic(err)
	}
	viper.SetDefault("DB_ENGINE", "postgres")

	// =========================
	// 2. Try read .env
	// =========================
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println(".env not found, using system environment")
		} else {
			panic(err)
		}
	} else {
		fmt.Println("Using .env:", viper.ConfigFileUsed())
	}

	// =========================
	// 3. Validate
	// =========================
	requiredEnv := []string{
		"DATABASE_URL",
		"DB_ENGINE",
		"jwt_secret_key",
	}

	for _, env := range requiredEnv {
		if viper.GetString(env) == "" {
			panic(fmt.Sprintf(
				"%s environment variable is required",
				env,
			))
		}
	}

	fmt.Println("Successfully loaded configuration")
}
