package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"reshier/models"
	"reshier/utils"
	"strconv"
	"strings"
)

func TampilkanBarang(w http.ResponseWriter, r *http.Request) {
	if err := models.InitDB(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := models.DBGetAllBarang()
	if err != nil {
		http.Error(w, "Gagal ambil data barang: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Sorting (in-memory setelah ambil dari DB)
	sort := r.URL.Query().Get("sort")
	switch sort {
	case "kode-asc":
		utils.UrutkanKodeBarang(data, true)
	case "kode-desc":
		utils.UrutkanKodeBarang(data, false)
	case "harga-asc":
		utils.UrutkanHargaBarang(data, true)
	case "harga-desc":
		utils.UrutkanHargaBarang(data, false)
	case "stok-asc":
		utils.UrutkanStokBarang(data, true)
	case "stok-desc":
		utils.UrutkanStokBarang(data, false)
	}

	// Filter pencarian
	search := r.URL.Query().Get("q")
	filtered := data
	if search != "" {
		filtered = utils.CariBarangByNama(data, search)
	}

	// Filter kategori
	filterKategori := r.URL.Query().Get("kategori")
	if filterKategori != "" && search == "" {
		var byKategori []models.Barang
		for _, b := range filtered {
			if b.Kategori == filterKategori {
				byKategori = append(byKategori, b)
			}
		}
		filtered = byKategori
	}

	// Hitung kategori unik dari semua data (bukan filtered)
	kategoriMap := make(map[string]bool)
	for _, b := range data {
		if b.Kategori != "" {
			kategoriMap[b.Kategori] = true
		}
	}
	var kategoriList []string
	for k := range kategoriMap {
		kategoriList = append(kategoriList, k)
	}

	stokKritis := 0
	for _, b := range data {
		if b.Stok <= 5 {
			stokKritis++
		}
	}

	type PageData struct {
		Barang       []models.Barang
		KategoriList []string
		Sort         string
		Search       string
		Kategori     string
		StokKritis   int
	}

	tmpl, err := template.New("barang.html").Funcs(utils.TemplateFuncs).ParseFiles("views/barang.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, PageData{
		Barang:       filtered,
		KategoriList: kategoriList,
		Sort:         sort,
		Search:       search,
		Kategori:     filterKategori,
		StokKritis:   stokKritis,
	})
}

func TambahBarang(w http.ResponseWriter, r *http.Request) {
	if err := models.InitDB(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		kode := strings.TrimSpace(r.FormValue("kode"))
		nama := strings.TrimSpace(r.FormValue("nama"))
		kategori := strings.TrimSpace(r.FormValue("kategori"))
		harga, _ := strconv.Atoi(r.FormValue("harga"))
		stok, _ := strconv.Atoi(r.FormValue("stok"))

		if kode == "" || nama == "" {
			http.Error(w, "Kode dan Nama tidak boleh kosong", http.StatusBadRequest)
			return
		}
		if harga < 0 || stok < 0 {
			http.Error(w, "Harga dan stok tidak boleh negatif", http.StatusBadRequest)
			return
		}

		exists, err := models.DBKodeExists(kode)
		if err != nil {
			http.Error(w, "Gagal cek kode: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Kode barang sudah ada: "+kode, http.StatusConflict)
			return
		}

		barang := models.Barang{
			Kode:     kode,
			Nama:     nama,
			Kategori: kategori,
			Harga:    harga,
			Stok:     stok,
		}
		if err := models.DBSaveBarang(barang); err != nil {
			http.Error(w, "Gagal simpan barang: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/barang", http.StatusSeeOther)
	} else {
		tmpl, err := template.New("tambah_barang.html").Funcs(utils.TemplateFuncs).ParseFiles("views/tambah_barang.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

func EditBarang(w http.ResponseWriter, r *http.Request) {
	if err := models.InitDB(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		kode := r.FormValue("kode")
		nama := strings.TrimSpace(r.FormValue("nama"))
		harga, _ := strconv.Atoi(r.FormValue("harga"))
		stok, _ := strconv.Atoi(r.FormValue("stok"))
		kategori := strings.TrimSpace(r.FormValue("kategori"))

		barang := models.Barang{
			Kode:     kode,
			Nama:     nama,
			Kategori: kategori,
			Harga:    harga,
			Stok:     stok,
		}
		if err := models.DBUpdateBarang(barang); err != nil {
			http.Error(w, "Gagal update barang: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/barang", http.StatusSeeOther)
	} else {
		kode := r.URL.Query().Get("kode")
		barang, err := models.DBGetBarangByKode(kode)
		if err != nil || barang == nil {
			http.NotFound(w, r)
			return
		}
		tmpl, err := template.New("edit_barang.html").Funcs(utils.TemplateFuncs).ParseFiles("views/edit_barang.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, *barang)
	}
}

func HapusBarang(w http.ResponseWriter, r *http.Request) {
	if err := models.InitDB(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	kode := r.URL.Query().Get("kode")
	if err := models.DBDeleteBarang(kode); err != nil {
		http.Error(w, "Gagal hapus barang: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/barang", http.StatusSeeOther)
}

// CariBarangJSON handler untuk live search (API JSON)
func CariBarangJSON(w http.ResponseWriter, r *http.Request) {
	if err := models.InitDB(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	q := r.URL.Query().Get("q")
	data, err := models.DBGetAllBarang()
	if err != nil {
		http.Error(w, "Gagal ambil data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var hasil []models.Barang
	if q == "" {
		hasil = data
	} else {
		hasil = utils.CariBarangByNama(data, q)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hasil)
}
