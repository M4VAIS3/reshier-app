// Package handler adalah entry point untuk Vercel Go serverless function.
// Vercel memanggil Handler(w, r) untuk setiap HTTP request.
package handler

import (
	"fmt"
	"log"
	"net/http"
	"reshier/controllers"
	"runtime/debug"
)

// Handler adalah fungsi serverless yang dipanggil Vercel.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Tangkap panic agar muncul sebagai HTTP error, bukan FUNCTION_INVOCATION_FAILED
	defer func() {
		if rec := recover(); rec != nil {
			stack := debug.Stack()
			log.Printf("[PANIC] %v\n%s", rec, stack)
			http.Error(w, fmt.Sprintf("Server error: %v", rec), http.StatusInternalServerError)
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
