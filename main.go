//go:build ignore
// File ini hanya untuk local development. Jalankan dengan: go run main.go
// Vercel mengabaikan file ini dan menggunakan api/index.go sebagai serverless handler.

package main

import (
	"fmt"
	"log"
	"net/http"
	"reshier/controllers"
	"reshier/models"
)

func main() {
	// Inisialisasi koneksi database
	if err := models.InitDB(); err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}
	fmt.Println("✅ Database terhubung")

	// === Dashboard ===
	http.HandleFunc("/", controllers.Dashboard)

	// === Barang ===
	http.HandleFunc("/barang", controllers.TampilkanBarang)
	http.HandleFunc("/barang/tambah", controllers.TambahBarang)
	http.HandleFunc("/barang/edit", controllers.EditBarang)
	http.HandleFunc("/barang/hapus", controllers.HapusBarang)
	http.HandleFunc("/barang/cari", controllers.CariBarangJSON)

	// === Transaksi ===
	http.HandleFunc("/transaksi", controllers.TampilkanTransaksi)
	http.HandleFunc("/transaksi/tambah", controllers.TambahTransaksi)
	http.HandleFunc("/transaksi/detail", controllers.DetailTransaksi)

	// === Laporan ===
	http.HandleFunc("/laporan", controllers.LaporanHarian)

	// === Static Files ===
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("🚀 Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
