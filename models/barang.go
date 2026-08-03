package models

// Barang merepresentasikan produk/menu di kasir.
type Barang struct {
	Kode       string
	Nama       string
	Kategori   string
	Harga      int
	Stok       int
	StokKritis bool // true jika stok <= 5
}

// DBGetAllBarang mengambil semua data barang dari database.
func DBGetAllBarang() ([]Barang, error) {
	rows, err := db.Query(
		"SELECT kode, nama, kategori, harga, stok FROM barang ORDER BY kode ASC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Barang
	for rows.Next() {
		var b Barang
		if err := rows.Scan(&b.Kode, &b.Nama, &b.Kategori, &b.Harga, &b.Stok); err != nil {
			return nil, err
		}
		b.StokKritis = b.Stok <= 5
		result = append(result, b)
	}
	if result == nil {
		result = []Barang{}
	}
	return result, nil
}

// DBGetBarangByKode mengambil satu barang berdasarkan kode.
func DBGetBarangByKode(kode string) (*Barang, error) {
	var b Barang
	err := db.QueryRow(
		"SELECT kode, nama, kategori, harga, stok FROM barang WHERE kode=$1", kode,
	).Scan(&b.Kode, &b.Nama, &b.Kategori, &b.Harga, &b.Stok)
	if err != nil {
		return nil, err
	}
	b.StokKritis = b.Stok <= 5
	return &b, nil
}

// DBSaveBarang menyimpan barang baru ke database.
func DBSaveBarang(b Barang) error {
	_, err := db.Exec(
		"INSERT INTO barang (kode, nama, kategori, harga, stok) VALUES ($1,$2,$3,$4,$5)",
		b.Kode, b.Nama, b.Kategori, b.Harga, b.Stok,
	)
	return err
}

// DBUpdateBarang memperbarui data barang di database.
func DBUpdateBarang(b Barang) error {
	_, err := db.Exec(
		"UPDATE barang SET nama=$1, kategori=$2, harga=$3, stok=$4 WHERE kode=$5",
		b.Nama, b.Kategori, b.Harga, b.Stok, b.Kode,
	)
	return err
}

// DBDeleteBarang menghapus barang dari database berdasarkan kode.
func DBDeleteBarang(kode string) error {
	_, err := db.Exec("DELETE FROM barang WHERE kode=$1", kode)
	return err
}

// DBKodeExists mengecek apakah kode barang sudah terdaftar.
func DBKodeExists(kode string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM barang WHERE kode=$1", kode).Scan(&count)
	return count > 0, err
}

// TandaiStokKritis memperbarui flag StokKritis pada slice barang.
func TandaiStokKritis(barang []Barang) {
	for i := range barang {
		barang[i].StokKritis = barang[i].Stok <= 5
	}
}
