package controllers

import (
	"encoding/json"
	"html/template"
	"reshier/models"
	"reshier/utils"
	"net/http"
	"strconv"
	"strings"
)

func TampilkanBarang(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("q")

	switch sort {
	case "kode-asc":
		utils.UrutkanKodeBarang(true)
	case "kode-desc":
		utils.UrutkanKodeBarang(false)
	case "harga-asc":
		utils.UrutkanHargaBarang(true)
	case "harga-desc":
		utils.UrutkanHargaBarang(false)
	case "stok-asc":
		utils.UrutkanStokBarang(true)
	case "stok-desc":
		utils.UrutkanStokBarang(false)
	}

	data := models.DataBarang
	if search != "" {
		data = utils.CariBarangByNama(search)
	}

	// Hitung kategori unik
	kategoriMap := make(map[string]bool)
	for _, b := range models.DataBarang {
		if b.Kategori != "" {
			kategoriMap[b.Kategori] = true
		}
	}
	var kategoriList []string
	for k := range kategoriMap {
		kategoriList = append(kategoriList, k)
	}

	filterKategori := r.URL.Query().Get("kategori")
	if filterKategori != "" && search == "" {
		var filtered []models.Barang
		for _, b := range data {
			if b.Kategori == filterKategori {
				filtered = append(filtered, b)
			}
		}
		data = filtered
	}

	type PageData struct {
		Barang       []models.Barang
		KategoriList []string
		Sort         string
		Search       string
		Kategori     string
		StokKritis   int
	}

	stokKritis := 0
	for _, b := range models.DataBarang {
		if b.Stok <= 5 {
			stokKritis++
		}
	}

	tmpl, err := template.New("barang.html").Funcs(utils.TemplateFuncs).ParseFiles("views/barang.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, PageData{
		Barang:       data,
		KategoriList: kategoriList,
		Sort:         sort,
		Search:       search,
		Kategori:     filterKategori,
		StokKritis:   stokKritis,
	})
}

func TambahBarang(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		kode := strings.TrimSpace(r.FormValue("kode"))
		nama := strings.TrimSpace(r.FormValue("nama"))
		kategori := strings.TrimSpace(r.FormValue("kategori"))
		harga, _ := strconv.Atoi(r.FormValue("harga"))
		stok, _ := strconv.Atoi(r.FormValue("stok"))

		// Validasi
		if kode == "" || nama == "" {
			http.Error(w, "Kode dan Nama tidak boleh kosong", http.StatusBadRequest)
			return
		}
		if utils.KodeBarangExists(kode) {
			http.Error(w, "Kode barang sudah ada: "+kode, http.StatusConflict)
			return
		}
		if harga < 0 || stok < 0 {
			http.Error(w, "Harga dan stok tidak boleh negatif", http.StatusBadRequest)
			return
		}

		barang := models.Barang{
			Kode:       kode,
			Nama:       nama,
			Kategori:   kategori,
			Harga:      harga,
			Stok:       stok,
			StokKritis: stok <= 5,
		}
		models.DataBarang = append(models.DataBarang, barang)
		models.SaveData()

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
	if r.Method == "POST" {
		kode := r.FormValue("kode")
		idx := utils.CariBarangSequential(kode)
		if idx != -1 {
			nama := strings.TrimSpace(r.FormValue("nama"))
			harga, _ := strconv.Atoi(r.FormValue("harga"))
			stok, _ := strconv.Atoi(r.FormValue("stok"))
			kategori := strings.TrimSpace(r.FormValue("kategori"))

			models.DataBarang[idx].Nama = nama
			models.DataBarang[idx].Harga = harga
			models.DataBarang[idx].Stok = stok
			models.DataBarang[idx].Kategori = kategori
			models.DataBarang[idx].StokKritis = stok <= 5
			models.SaveData()
		}
		http.Redirect(w, r, "/barang", http.StatusSeeOther)
	} else {
		kode := r.URL.Query().Get("kode")
		idx := utils.CariBarangSequential(kode)
		if idx != -1 {
			tmpl, err := template.New("edit_barang.html").Funcs(utils.TemplateFuncs).ParseFiles("views/edit_barang.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, models.DataBarang[idx])
		} else {
			http.NotFound(w, r)
		}
	}
}

func HapusBarang(w http.ResponseWriter, r *http.Request) {
	kode := r.URL.Query().Get("kode")
	idx := utils.CariBarangSequential(kode)
	if idx != -1 {
		models.DataBarang = append(models.DataBarang[:idx], models.DataBarang[idx+1:]...)
		models.SaveData()
	}
	http.Redirect(w, r, "/barang", http.StatusSeeOther)
}

// CariBarangJSON handler untuk live search (API JSON)
func CariBarangJSON(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	var hasil []models.Barang
	if q == "" {
		hasil = models.DataBarang
	} else {
		hasil = utils.CariBarangByNama(q)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hasil)
}
