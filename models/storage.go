package models

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB menginisialisasi koneksi PostgreSQL (singleton).
// Dipanggil sekali saat aplikasi start atau per-request di Vercel.
func InitDB() error {
	var initErr error
	once.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			initErr = fmt.Errorf("DATABASE_URL environment variable tidak di-set")
			return
		}
		db, initErr = sql.Open("postgres", dsn)
		if initErr != nil {
			return
		}
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		if initErr = db.Ping(); initErr != nil {
			initErr = fmt.Errorf("gagal konek ke database: %w", initErr)
			return
		}
		initErr = createTables()
	})
	return initErr
}

// GetDB mengembalikan instance DB yang sudah diinisialisasi.
func GetDB() *sql.DB {
	return db
}

// createTables membuat tabel jika belum ada.
func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS barang (
			kode      TEXT PRIMARY KEY,
			nama      TEXT NOT NULL,
			kategori  TEXT NOT NULL DEFAULT '',
			harga     INTEGER NOT NULL DEFAULT 0,
			stok      INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS transaksi (
			id        TEXT PRIMARY KEY,
			waktu     TIMESTAMPTZ NOT NULL,
			total     INTEGER NOT NULL DEFAULT 0,
			bayar     INTEGER NOT NULL DEFAULT 0,
			kembalian INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS item_transaksi (
			id             SERIAL PRIMARY KEY,
			transaksi_id   TEXT NOT NULL REFERENCES transaksi(id) ON DELETE CASCADE,
			kode_barang    TEXT NOT NULL,
			nama_barang    TEXT NOT NULL,
			harga          INTEGER NOT NULL DEFAULT 0,
			jumlah         INTEGER NOT NULL DEFAULT 0,
			subtotal       INTEGER NOT NULL DEFAULT 0
		)`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("gagal membuat tabel: %w", err)
		}
	}
	return nil
}
