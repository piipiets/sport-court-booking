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

	viper.BindEnv("DATABASE_URL", "DATABASE_URL")

	viper.BindEnv("SUPABASE_SECRET_KEY", "SUPABASE_SECRET_KEY")
	viper.BindEnv("SUPABASE_STORAGE_BUCKET", "SUPABASE_STORAGE_BUCKET")
	viper.BindEnv("SUPABASE_STORAGE_URL", "SUPABASE_STORAGE_URL")
	viper.BindEnv("DB_ENGINE", "DB_ENGINE")
	viper.BindEnv("jwt_secret_key", "JWT_SECRET_KEY", "jwt_secret_key")

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
		"SUPABASE_STORAGE_URL",
		"SUPABASE_SECRET_KEY",
		"SUPABASE_STORAGE_BUCKET",
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

func SupabaseStorageURL() string {
	return viper.GetString("SUPABASE_STORAGE_URL")
}

func SupabaseSecretKey() string {
	return viper.GetString("SUPABASE_SECRET_KEY")
}

func SupabaseStorageBucket() string {
	return viper.GetString("SUPABASE_STORAGE_BUCKET")
}
