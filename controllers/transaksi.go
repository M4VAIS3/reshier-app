package controllers

import (
	"fmt"
	"net/http"
	"reshier/models"
	"reshier/utils"
	"strconv"
	"time"
)

func TampilkanTransaksi(w http.ResponseWriter, r *http.Request) {
	if err := models.InitDB(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	semua, err := models.DBGetAllTransaksi()
	if err != nil {
		http.Error(w, "Gagal ambil transaksi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filter := r.URL.Query().Get("q")
	hasil := semua
	if filter != "" {
		hasil = utils.FilterTransaksiByTime(semua, filter)
	}

	// Balik urutan agar terbaru di atas (DB sudah ORDER BY waktu DESC, balik lagi agar search tetap konsisten)
	reversed := make([]models.Transaksi, len(hasil))
	for i, t := range hasil {
		reversed[len(hasil)-1-i] = t
	}
	// Sudah DESC dari DB, tidak perlu balik lagi
	reversed = hasil

	total := 0
	for _, t := range hasil {
		total += t.Total
	}

	type PageData struct {
		Transaksi []models.Transaksi
		Filter    string
		Total     int
	}

	tmpl, err := parseTemplate("transaksi.html", "views/transaksi.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, PageData{Transaksi: reversed, Filter: filter, Total: total})
}

func TambahTransaksi(w http.ResponseWriter, r *http.Request) {
	if err := models.InitDB(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		bayar, _ := strconv.Atoi(r.FormValue("bayar"))

		// Ambil semua barang dari DB untuk lookup
		semuaBarang, err := models.DBGetAllBarang()
		if err != nil {
			http.Error(w, "Gagal ambil data barang: "+err.Error(), http.StatusInternalServerError)
			return
		}

		now := time.Now()
		count, err := models.DBCountTransaksi()
		if err != nil {
			http.Error(w, "Gagal generate ID transaksi: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var transaksi models.Transaksi
		transaksi.ID = fmt.Sprintf("TRX-%s-%03d", now.Format("20060102"), count+1)
		transaksi.Waktu = now

		for i := 0; i < 20; i++ {
			kode := r.FormValue("kode" + strconv.Itoa(i))
			if kode == "" {
				continue
			}
			idx := utils.CariBarangSequential(semuaBarang, kode)
			if idx == -1 {
				continue
			}
			jumlah, _ := strconv.Atoi(r.FormValue("jumlah" + strconv.Itoa(i)))
			if jumlah <= 0 {
				continue
			}
			if jumlah > semuaBarang[idx].Stok {
				continue
			}

			item := models.ItemTransaksi{
				KodeBarang: semuaBarang[idx].Kode,
				NamaBarang: semuaBarang[idx].Nama,
				Harga:      semuaBarang[idx].Harga,
				Jumlah:     jumlah,
				Subtotal:   semuaBarang[idx].Harga * jumlah,
			}
			transaksi.Items = append(transaksi.Items, item)
			transaksi.Total += item.Subtotal
		}

		if len(transaksi.Items) > 0 {
			transaksi.Bayar = bayar
			transaksi.Kembalian = bayar - transaksi.Total

			if err := models.DBSaveTransaksi(transaksi); err != nil {
				http.Error(w, "Gagal simpan transaksi: "+err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/transaksi/detail?id="+transaksi.ID, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/transaksi/tambah", http.StatusSeeOther)
		}
	} else {
		// GET: tampilkan form kasir dengan daftar barang
		semuaBarang, err := models.DBGetAllBarang()
		if err != nil {
			http.Error(w, "Gagal ambil data barang: "+err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := parseTemplate("tambah_transaksi.html", "views/tambah_transaksi.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, semuaBarang)
	}
}

func DetailTransaksi(w http.ResponseWriter, r *http.Request) {
	if err := models.InitDB(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id := r.URL.Query().Get("id")
	found, err := models.DBGetTransaksiByID(id)
	if err != nil {
		http.Error(w, "Gagal ambil detail transaksi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if found == nil {
		http.NotFound(w, r)
		return
	}

	tmpl, err := parseTemplate("detail_transaksi.html", "views/detail_transaksi.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, found)
}
