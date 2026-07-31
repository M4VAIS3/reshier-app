package controllers

import (
	"fmt"
	"html/template"
	"reshier/models"
	"reshier/utils"
	"net/http"
	"strconv"
	"time"
)

func TampilkanTransaksi(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	var hasil []models.Transaksi

	if filter != "" {
		hasil = utils.FilterTransaksiByTime(filter)
	} else {
		hasil = models.DataTransaksi
	}

	// Balik urutan terbaru di atas
	reversed := make([]models.Transaksi, len(hasil))
	for i, v := range hasil {
		reversed[len(hasil)-1-i] = v
	}

	type PageData struct {
		Transaksi []models.Transaksi
		Filter    string
		Total     int
	}

	total := 0
	for _, t := range hasil {
		total += t.Total
	}

	tmpl, err := template.New("transaksi.html").Funcs(utils.TemplateFuncs).ParseFiles("views/transaksi.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, PageData{Transaksi: reversed, Filter: filter, Total: total})
}

func TambahTransaksi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var transaksi models.Transaksi
		now := time.Now()
		transaksi.ID = fmt.Sprintf("TRX-%s-%03d", now.Format("20060102"), len(models.DataTransaksi)+1)
		transaksi.Waktu = now

		bayar, _ := strconv.Atoi(r.FormValue("bayar"))

		for i := 0; i < 20; i++ {
			kode := r.FormValue("kode" + strconv.Itoa(i))
			if kode == "" {
				continue
			}
			idx := utils.CariBarangSequential(kode)
			if idx == -1 {
				continue
			}
			jumlah, _ := strconv.Atoi(r.FormValue("jumlah" + strconv.Itoa(i)))
			if jumlah <= 0 {
				continue
			}
			if jumlah > models.DataBarang[idx].Stok {
				continue
			}

			models.DataBarang[idx].Stok -= jumlah
			models.DataBarang[idx].StokKritis = models.DataBarang[idx].Stok <= 5

			item := models.ItemTransaksi{
				KodeBarang: models.DataBarang[idx].Kode,
				NamaBarang: models.DataBarang[idx].Nama,
				Harga:      models.DataBarang[idx].Harga,
				Jumlah:     jumlah,
				Subtotal:   models.DataBarang[idx].Harga * jumlah,
			}
			transaksi.Items = append(transaksi.Items, item)
			transaksi.Total += item.Subtotal
		}

		if len(transaksi.Items) > 0 {
			transaksi.Bayar = bayar
			transaksi.Kembalian = bayar - transaksi.Total
			models.DataTransaksi = append(models.DataTransaksi, transaksi)
			models.SaveData()
			http.Redirect(w, r, "/transaksi/detail?id="+transaksi.ID, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/transaksi/tambah", http.StatusSeeOther)
		}
	} else {
		tmpl, err := template.New("tambah_transaksi.html").Funcs(utils.TemplateFuncs).ParseFiles("views/tambah_transaksi.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, models.DataBarang)
	}
}

func DetailTransaksi(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var found *models.Transaksi

	for i := range models.DataTransaksi {
		if models.DataTransaksi[i].ID == id {
			found = &models.DataTransaksi[i]
			break
		}
	}

	if found == nil {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.New("detail_transaksi.html").Funcs(utils.TemplateFuncs).ParseFiles("views/detail_transaksi.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, found)
}
