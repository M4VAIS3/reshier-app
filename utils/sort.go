package utils

import "reshier/models"

// Selection Sort berdasarkan Kode Barang
func UrutkanKodeBarang(ascending bool) {
	n := len(models.DataBarang)
	for i := 0; i < n-1; i++ {
		idx := i
		for j := i + 1; j < n; j++ {
			if ascending {
				if models.DataBarang[j].Kode < models.DataBarang[idx].Kode {
					idx = j
				}
			} else {
				if models.DataBarang[j].Kode > models.DataBarang[idx].Kode {
					idx = j
				}
			}
		}
		models.DataBarang[i], models.DataBarang[idx] = models.DataBarang[idx], models.DataBarang[i]
	}
}

// Insertion Sort berdasarkan Harga
func UrutkanHargaBarang(ascending bool) {
	n := len(models.DataBarang)
	for i := 1; i < n; i++ {
		temp := models.DataBarang[i]
		j := i - 1
		if ascending {
			for j >= 0 && models.DataBarang[j].Harga > temp.Harga {
				models.DataBarang[j+1] = models.DataBarang[j]
				j--
			}
		} else {
			for j >= 0 && models.DataBarang[j].Harga < temp.Harga {
				models.DataBarang[j+1] = models.DataBarang[j]
				j--
			}
		}
		models.DataBarang[j+1] = temp
	}
}

// Insertion Sort berdasarkan Stok
func UrutkanStokBarang(ascending bool) {
	n := len(models.DataBarang)
	for i := 1; i < n; i++ {
		temp := models.DataBarang[i]
		j := i - 1
		if ascending {
			for j >= 0 && models.DataBarang[j].Stok > temp.Stok {
				models.DataBarang[j+1] = models.DataBarang[j]
				j--
			}
		} else {
			for j >= 0 && models.DataBarang[j].Stok < temp.Stok {
				models.DataBarang[j+1] = models.DataBarang[j]
				j--
			}
		}
		models.DataBarang[j+1] = temp
	}
}

// BarangTerlaris mengembalikan top N barang berdasarkan total jumlah terjual
func BarangTerlaris(topN int) []models.BarangLaris {
	tally := make(map[string]*models.BarangLaris)

	for _, trx := range models.DataTransaksi {
		for _, item := range trx.Items {
			if _, ok := tally[item.KodeBarang]; !ok {
				tally[item.KodeBarang] = &models.BarangLaris{
					KodeBarang: item.KodeBarang,
					NamaBarang: item.NamaBarang,
				}
			}
			tally[item.KodeBarang].TotalTerjual += item.Jumlah
			tally[item.KodeBarang].TotalPendapatan += item.Subtotal
		}
	}

	var list []models.BarangLaris
	for _, v := range tally {
		list = append(list, *v)
	}

	// Bubble sort descending berdasarkan TotalTerjual (untuk showcase algoritma)
	for i := 0; i < len(list)-1; i++ {
		for j := 0; j < len(list)-i-1; j++ {
			if list[j].TotalTerjual < list[j+1].TotalTerjual {
				list[j], list[j+1] = list[j+1], list[j]
			}
		}
	}

	if topN > len(list) {
		topN = len(list)
	}
	return list[:topN]
}
