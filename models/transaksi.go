package models

import (
	"database/sql"
	"fmt"
	"time"
)

// ItemTransaksi merepresentasikan satu baris item dalam transaksi.
type ItemTransaksi struct {
	KodeBarang string
	NamaBarang string
	Harga      int
	Jumlah     int
	Subtotal   int
}

// Transaksi merepresentasikan satu sesi transaksi penjualan.
type Transaksi struct {
	ID        string
	Waktu     time.Time
	Items     []ItemTransaksi
	Total     int
	Bayar     int
	Kembalian int
}

// DBGetAllTransaksi mengambil semua transaksi beserta item-nya dari database.
func DBGetAllTransaksi() ([]Transaksi, error) {
	rows, err := db.Query(
		"SELECT id, waktu, total, bayar, kembalian FROM transaksi ORDER BY waktu DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Transaksi
	for rows.Next() {
		var t Transaksi
		if err := rows.Scan(&t.ID, &t.Waktu, &t.Total, &t.Bayar, &t.Kembalian); err != nil {
			return nil, err
		}
		items, err := dbGetItemsByTransaksiID(t.ID)
		if err != nil {
			return nil, err
		}
		t.Items = items
		result = append(result, t)
	}
	if result == nil {
		result = []Transaksi{}
	}
	return result, nil
}

// DBGetTransaksiByID mengambil satu transaksi berdasarkan ID.
func DBGetTransaksiByID(id string) (*Transaksi, error) {
	var t Transaksi
	err := db.QueryRow(
		"SELECT id, waktu, total, bayar, kembalian FROM transaksi WHERE id=$1", id,
	).Scan(&t.ID, &t.Waktu, &t.Total, &t.Bayar, &t.Kembalian)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	items, err := dbGetItemsByTransaksiID(t.ID)
	if err != nil {
		return nil, err
	}
	t.Items = items
	return &t, nil
}

// DBCountTransaksi mengembalikan jumlah total transaksi (untuk generate ID).
func DBCountTransaksi() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM transaksi").Scan(&count)
	return count, err
}

// DBSaveTransaksi menyimpan transaksi baru dan mengurangi stok barang secara atomik.
func DBSaveTransaksi(t Transaksi) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("gagal mulai transaksi DB: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO transaksi (id, waktu, total, bayar, kembalian) VALUES ($1,$2,$3,$4,$5)",
		t.ID, t.Waktu, t.Total, t.Bayar, t.Kembalian,
	)
	if err != nil {
		return fmt.Errorf("gagal insert transaksi: %w", err)
	}

	for _, item := range t.Items {
		_, err = tx.Exec(
			`INSERT INTO item_transaksi
				(transaksi_id, kode_barang, nama_barang, harga, jumlah, subtotal)
			VALUES ($1,$2,$3,$4,$5,$6)`,
			t.ID, item.KodeBarang, item.NamaBarang, item.Harga, item.Jumlah, item.Subtotal,
		)
		if err != nil {
			return fmt.Errorf("gagal insert item transaksi: %w", err)
		}

		// Kurangi stok secara atomik di database
		res, err := tx.Exec(
			"UPDATE barang SET stok = stok - $1 WHERE kode = $2 AND stok >= $1",
			item.Jumlah, item.KodeBarang,
		)
		if err != nil {
			return fmt.Errorf("gagal update stok: %w", err)
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return fmt.Errorf("stok barang %s tidak mencukupi", item.KodeBarang)
		}
	}

	return tx.Commit()
}

// dbGetItemsByTransaksiID adalah helper internal untuk ambil items berdasarkan transaksi ID.
func dbGetItemsByTransaksiID(transaksiID string) ([]ItemTransaksi, error) {
	rows, err := db.Query(
		`SELECT kode_barang, nama_barang, harga, jumlah, subtotal
		FROM item_transaksi WHERE transaksi_id=$1`,
		transaksiID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ItemTransaksi
	for rows.Next() {
		var item ItemTransaksi
		if err := rows.Scan(&item.KodeBarang, &item.NamaBarang, &item.Harga, &item.Jumlah, &item.Subtotal); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
