package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reshier/controllers"
	"reshier/models"
)

func main() {
	// Inisialisasi koneksi database
	if err := models.InitDB(); err != nil {
		// Log error tapi jangan langsung Fatal — biarkan server tetap start
		// agar Vercel tidak menganggap proses crash
		log.Println("⚠️  Gagal konek ke database:", err)
	} else {
		fmt.Println("✅ Database terhubung")
	}

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

	// Vercel menyediakan PORT env var — fallback ke 8080 untuk local dev
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server berjalan di http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
