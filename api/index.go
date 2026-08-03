// Package handler adalah entry point untuk Vercel Go serverless function.
// Vercel akan memanggil fungsi Handler(w, r) untuk setiap HTTP request.
package handler

import (
	"net/http"
	"reshier/controllers"
)

// Handler adalah fungsi yang dipanggil Vercel untuk setiap request.
// Semua routing didefinisikan di sini — mirip main.go tapi tanpa ListenAndServe.
func Handler(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	// === Dashboard ===
	mux.HandleFunc("/", controllers.Dashboard)

	// === Barang ===
	mux.HandleFunc("/barang", controllers.TampilkanBarang)
	mux.HandleFunc("/barang/tambah", controllers.TambahBarang)
	mux.HandleFunc("/barang/edit", controllers.EditBarang)
	mux.HandleFunc("/barang/hapus", controllers.HapusBarang)
	mux.HandleFunc("/barang/cari", controllers.CariBarangJSON)

	// === Transaksi ===
	mux.HandleFunc("/transaksi", controllers.TampilkanTransaksi)
	mux.HandleFunc("/transaksi/tambah", controllers.TambahTransaksi)
	mux.HandleFunc("/transaksi/detail", controllers.DetailTransaksi)

	// === Laporan ===
	mux.HandleFunc("/laporan", controllers.LaporanHarian)

	// === Static Files (untuk local dev; di Vercel dilayani lewat vercel.json routes) ===
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.ServeHTTP(w, r)
}
