package controllers

import (
	"html/template"
	"net/http"
	"reshier/models"
	"reshier/utils"
	"time"
)

// DashboardData adalah data yang dikirim ke template dashboard.
type DashboardData struct {
	TotalBarang       int
	TotalTransaksi    int
	OmzetHarian       int
	StokKritis        int
	TransaksiTerakhir []models.Transaksi
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	if err := models.InitDB(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	barang, err := models.DBGetAllBarang()
	if err != nil {
		http.Error(w, "Gagal ambil data barang: "+err.Error(), http.StatusInternalServerError)
		return
	}

	transaksi, err := models.DBGetAllTransaksi()
	if err != nil {
		http.Error(w, "Gagal ambil data transaksi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	data := DashboardData{
		TotalBarang:    len(barang),
		TotalTransaksi: len(transaksi),
	}

	for _, trx := range transaksi {
		if trx.Waktu.Format("2006-01-02") == today {
			data.OmzetHarian += trx.Total
		}
	}

	for _, b := range barang {
		if b.Stok <= 5 {
			data.StokKritis++
		}
	}

	// 5 transaksi terakhir (transaksi sudah ORDER BY waktu DESC dari DB)
	n := len(transaksi)
	end := 5
	if end > n {
		end = n
	}
	data.TransaksiTerakhir = transaksi[:end]

	tmpl, err := template.New("index.html").Funcs(utils.TemplateFuncs).ParseFiles("views/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
