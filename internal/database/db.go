package database

// import "database/sql"

// Koneksi ke database
import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	// Pastikan Anda mengganti konfigurasi koneksi database sesuai dengan pengaturan Anda
	// Misalnya, Anda dapat menggunakan environment variable atau file konfigurasi
	// Untuk contoh ini, saya akan menggunakan koneksi langsung dengan MySQL
	var err error
	DB, err = sql.Open("mysql", "root@tcp(localhost:3306)/majoris_sidu")
	if err != nil {
		log.Fatal(err)
	}

	// Uji koneksi ke database
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected successfully")
}

// Fungsi ini dapat digunakan untuk mendapatkan koneksi database yang dibuka
func GetDB() *sql.DB {
	return DB
}
