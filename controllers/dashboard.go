package controllers

import (
	"html/template"
	"reshier/models"
	"reshier/utils"
	"net/http"
	"time"
)

type DashboardData struct {
	TotalBarang       int
	TotalTransaksi    int
	OmzetHarian       int
	StokKritis        int
	TransaksiTerakhir []models.Transaksi
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	today := now.Format("2006-01-02")

	data := DashboardData{
		TotalBarang:    len(models.DataBarang),
		TotalTransaksi: len(models.DataTransaksi),
	}

	for _, trx := range models.DataTransaksi {
		if trx.Waktu.Format("2006-01-02") == today {
			data.OmzetHarian += trx.Total
		}
	}

	for _, b := range models.DataBarang {
		if b.Stok <= 5 {
			data.StokKritis++
		}
	}

	// 5 transaksi terakhir (terbaru di atas)
	all := models.DataTransaksi
	n := len(all)
	start := n - 5
	if start < 0 {
		start = 0
	}
	recent := make([]models.Transaksi, 0)
	for i := n - 1; i >= start; i-- {
		recent = append(recent, all[i])
	}
	data.TransaksiTerakhir = recent

	tmpl, err := template.New("index.html").Funcs(utils.TemplateFuncs).ParseFiles("views/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
