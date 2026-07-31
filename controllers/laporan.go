package controllers

import (
	"html/template"
	"net/http"
	"reshier/models"
	"reshier/utils"
	"time"
)

type LaporanData struct {
	TanggalHari     string
	OmzetHarian     int
	OmzetMingguan   int
	OmzetBulanan    int
	JmlTransHarian  int
	JmlTransBulanan int
	BarangTerlaris  []models.BarangLaris
	GrafikMingguan  []GrafikHarian
	StokKritis      []models.Barang
}

type GrafikHarian struct {
	Hari  string
	Total int
	Label string
}

func LaporanHarian(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	today := now.Format("2006-01-02")
	thisMonth := now.Format("2006-01")

	data := LaporanData{
		TanggalHari: now.Format("02 January 2006"),
	}

	for _, trx := range models.DataTransaksi {
		tDate := trx.Waktu.Format("2006-01-02")
		tMonth := trx.Waktu.Format("2006-01")

		if tDate == today {
			data.OmzetHarian += trx.Total
			data.JmlTransHarian++
		}
		if tMonth == thisMonth {
			data.OmzetBulanan += trx.Total
			data.JmlTransBulanan++
		}
		if trx.Waktu.After(now.AddDate(0, 0, -7)) {
			data.OmzetMingguan += trx.Total
		}
	}

	// Grafik 7 hari terakhir
	for i := 6; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		dayStr := day.Format("2006-01-02")
		dayLabel := day.Format("02/01")
		total := 0
		for _, trx := range models.DataTransaksi {
			if trx.Waktu.Format("2006-01-02") == dayStr {
				total += trx.Total
			}
		}
		data.GrafikMingguan = append(data.GrafikMingguan, GrafikHarian{
			Hari:  dayStr,
			Total: total,
			Label: dayLabel,
		})
	}

	data.BarangTerlaris = utils.BarangTerlaris(5)

	for _, b := range models.DataBarang {
		if b.Stok <= 5 {
			data.StokKritis = append(data.StokKritis, b)
		}
	}

	tmpl, err := template.New("laporan.html").Funcs(utils.TemplateFuncs).ParseFiles("views/laporan.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
