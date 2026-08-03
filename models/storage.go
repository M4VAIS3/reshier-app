package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db      *sql.DB
	dbErr   error // simpan error permanen dari InitDB
	once    sync.Once
)

// InitDB menginisialisasi koneksi PostgreSQL (singleton).
func InitDB() error {
	once.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			dbErr = fmt.Errorf("DATABASE_URL environment variable tidak di-set")
			log.Println("[DB ERROR]", dbErr)
			return
		}
		log.Println("[DB] Membuka koneksi ke PostgreSQL...")

		var err error
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			dbErr = fmt.Errorf("sql.Open gagal: %w", err)
			log.Println("[DB ERROR]", dbErr)
			return
		}

		db.SetMaxOpenConns(5)
		db.SetMaxIdleConns(2)

		if err = db.Ping(); err != nil {
			dbErr = fmt.Errorf("db.Ping() gagal: %w", err)
			log.Println("[DB ERROR]", dbErr)
			db = nil
			return
		}

		log.Println("[DB] Koneksi berhasil! Membuat tabel...")
		if err = createTables(); err != nil {
			dbErr = fmt.Errorf("createTables() gagal: %w", err)
			log.Println("[DB ERROR]", dbErr)
			return
		}
		log.Println("[DB] Siap!")
	})
	return dbErr
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
