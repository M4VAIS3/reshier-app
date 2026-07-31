package main

import (
	"fmt"
	"reshier/controllers"
	"reshier/models"
	"log"
	"net/http"
)

func main() {
	// Load data saat startup
	if err := models.LoadData(); err != nil {
		log.Println("Warning: gagal load data:", err)
	}
	fmt.Println("Data berhasil dimuat:", len(models.DataBarang), "barang,", len(models.DataTransaksi), "transaksi")

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

	fmt.Println("Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
