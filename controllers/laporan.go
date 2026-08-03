package controllers

import (
	"html/template"
	"net/http"
	"reshier/models"
	"reshier/utils"
	"time"
)

// LaporanData adalah data yang dikirim ke template laporan.
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

// GrafikHarian adalah data per-hari untuk grafik omzet mingguan.
type GrafikHarian struct {
	Hari  string
	Total int
	Label string
}

func LaporanHarian(w http.ResponseWriter, r *http.Request) {
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
	thisMonth := now.Format("2006-01")

	data := LaporanData{
		TanggalHari: now.Format("02 January 2006"),
	}

	for _, trx := range transaksi {
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
		for _, trx := range transaksi {
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

	data.BarangTerlaris = utils.BarangTerlaris(transaksi, 5)

	for _, b := range barang {
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
