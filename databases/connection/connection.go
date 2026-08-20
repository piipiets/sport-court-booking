package connection

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

var (
	DBConnections *sql.DB
	err           error
)

func Initiator() {
	dbEngine := viper.GetString("DB_ENGINE")
	dsn := viper.GetString("DATABASE_URL")

	DBConnections, err = sql.Open(dbEngine, dsn)
	if err != nil {
		panic(err)
	}

	// --- Connection pooling ---
	DBConnections.SetMaxOpenConns(20)                  // batas maksimum koneksi aktif ke DB
	DBConnections.SetMaxIdleConns(10)                  // koneksi menganggur yang tetap disimpan (siap pakai)
	DBConnections.SetConnMaxLifetime(30 * time.Minute) // paksa recycle koneksi tua
	DBConnections.SetConnMaxIdleTime(5 * time.Minute)  // tutup koneksi yg nganggur kelamaan

	// check connection
	err = DBConnections.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database")
}
