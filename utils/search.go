package utils

import (
	"reshier/models"
	"strings"
)

// CariBarangSequential - Sequential search berdasarkan kode (exact match)
func CariBarangSequential(data []models.Barang, kode string) int {
	for i, b := range data {
		if b.Kode == kode {
			return i
		}
	}
	return -1
}

// CariBarangBinary - Binary Search berdasarkan kode (data harus terurut ascending)
func CariBarangBinary(data []models.Barang, kode string) int {
	low, high := 0, len(data)-1
	for low <= high {
		mid := (low + high) / 2
		if data[mid].Kode == kode {
			return mid
		} else if data[mid].Kode < kode {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}

// CariBarangByNama mencari barang yang mengandung kata kunci (case-insensitive)
func CariBarangByNama(data []models.Barang, query string) []models.Barang {
	var hasil []models.Barang
	q := strings.ToLower(query)
	for _, b := range data {
		if strings.Contains(strings.ToLower(b.Nama), q) ||
			strings.Contains(strings.ToLower(b.Kode), q) ||
			strings.Contains(strings.ToLower(b.Kategori), q) {
			hasil = append(hasil, b)
		}
	}
	return hasil
}

// FilterTransaksiByTime memfilter transaksi berdasarkan query waktu
func FilterTransaksiByTime(data []models.Transaksi, query string) []models.Transaksi {
	var hasil []models.Transaksi
	for _, trx := range data {
		t := trx.Waktu
		match := false
		switch len(query) {
		case 16:
			match = t.Format("02-01-2006 15:04") == query
		case 13:
			match = t.Format("02-01-2006 15") == query
		case 10:
			match = t.Format("02-01-2006") == query
		case 7:
			match = t.Format("01-2006") == query
		case 4:
			match = t.Format("2006") == query
		case 5:
			match = t.Format("15:04") == query
		case 2:
			match = t.Format("15") == query
		}
		if match {
			hasil = append(hasil, trx)
		}
	}
	return hasil
}
