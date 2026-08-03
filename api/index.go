// Package handler adalah entry point untuk Vercel Go serverless function.
package handler

import (
	"log"
	"net/http"
	"reshier/controllers"
	"runtime/debug"
)

// Handler adalah fungsi yang dipanggil Vercel untuk setiap HTTP request.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Recover dari panic agar error terlihat di logs, bukan FUNCTION_INVOCATION_FAILED
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("[PANIC] %v\n%s", rec, debug.Stack())
			http.Error(w, "Internal server error (panic recovered)", http.StatusInternalServerError)
		}
	}()

	log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)

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

	// === Static Files ===
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.ServeHTTP(w, r)
}
