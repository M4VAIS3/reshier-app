package utils

import "reshier/models"

// UrutkanKodeBarang - Selection Sort berdasarkan Kode Barang (in-place)
func UrutkanKodeBarang(data []models.Barang, ascending bool) {
	n := len(data)
	for i := 0; i < n-1; i++ {
		idx := i
		for j := i + 1; j < n; j++ {
			if ascending {
				if data[j].Kode < data[idx].Kode {
					idx = j
				}
			} else {
				if data[j].Kode > data[idx].Kode {
					idx = j
				}
			}
		}
		data[i], data[idx] = data[idx], data[i]
	}
}

// UrutkanHargaBarang - Insertion Sort berdasarkan Harga (in-place)
func UrutkanHargaBarang(data []models.Barang, ascending bool) {
	n := len(data)
	for i := 1; i < n; i++ {
		temp := data[i]
		j := i - 1
		if ascending {
			for j >= 0 && data[j].Harga > temp.Harga {
				data[j+1] = data[j]
				j--
			}
		} else {
			for j >= 0 && data[j].Harga < temp.Harga {
				data[j+1] = data[j]
				j--
			}
		}
		data[j+1] = temp
	}
}

// UrutkanStokBarang - Insertion Sort berdasarkan Stok (in-place)
func UrutkanStokBarang(data []models.Barang, ascending bool) {
	n := len(data)
	for i := 1; i < n; i++ {
		temp := data[i]
		j := i - 1
		if ascending {
			for j >= 0 && data[j].Stok > temp.Stok {
				data[j+1] = data[j]
				j--
			}
		} else {
			for j >= 0 && data[j].Stok < temp.Stok {
				data[j+1] = data[j]
				j--
			}
		}
		data[j+1] = temp
	}
}

// BarangTerlaris mengembalikan top N barang berdasarkan total jumlah terjual.
// Menggunakan Bubble Sort descending untuk showcase algoritma.
func BarangTerlaris(transaksi []models.Transaksi, topN int) []models.BarangLaris {
	tally := make(map[string]*models.BarangLaris)

	for _, trx := range transaksi {
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

	// Bubble Sort descending berdasarkan TotalTerjual
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
