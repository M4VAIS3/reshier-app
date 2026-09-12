// Package handler adalah entry point untuk Vercel Go serverless function.
// Vercel memanggil Handler(w, r) untuk setiap HTTP request.
package handler

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"reshier/controllers"
	"runtime/debug"
	"sync"
)

//go:embed all:views
var viewsFS embed.FS

//go:embed all:static
var staticFS embed.FS

var initOnce sync.Once

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

	// Set embedded FS untuk controllers (sekali saja)
	initOnce.Do(func() {
		controllers.ViewsFS = viewsFS
	})

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

	// === Static Files (dari embed) ===
	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		http.Error(w, "Static FS error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	mux.ServeHTTP(w, r)
}
