# Reshier — Simple Restaurant Cashier App

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supabase-336791?logo=postgresql&logoColor=white)](https://supabase.com/)
[![Deploy](https://img.shields.io/badge/Deployed%20on-Vercel-000000?logo=vercel&logoColor=white)](https://reshier-app.vercel.app/)

> A web-based restaurant cashier application built with **Go (net/http)** and **PostgreSQL (Supabase)**. Deployed as a serverless function on **Vercel**.

🔗 **Live Demo:** [reshier-app.vercel.app](https://reshier-app.vercel.app/)

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Application Architecture](#application-architecture)
- [Database Schema](#database-schema)
- [Data Models](#data-models)
- [Routing & Endpoints](#routing--endpoints)
- [Algorithms Used](#algorithms-used)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Deployment](#deployment)

---

## Overview

**Reshier** is a lightweight point-of-sale application designed for small restaurants. Built using Go's standard library (`net/http`) and rendered via Go HTML Templates, Reshier stores all data persistently in a **PostgreSQL** database hosted on [Supabase](https://supabase.com/).

The application is deployed on **Vercel** as a Go serverless function, with static assets and HTML templates embedded into the binary using Go's `embed` package.

This project also demonstrates the manual implementation of **sorting and searching algorithms** (without using `sort.Slice` from the standard library) as a learning exercise.

---

## Features

| Module | Description |
|---|---|
| **Dashboard** | Summary of total items, total transactions, daily revenue, and low-stock count |
| **Item Management** | Add, edit, and delete items; filter by category; sort by code, price, or stock |
| **Transactions** | Record multi-item sales transactions with automatic change calculation |
| **Live Search** | Real-time item search via a JSON API endpoint |
| **Reports** | Daily, weekly, and monthly reports; 7-day sales chart; best-selling items; low-stock alerts |
| **Low Stock Alert** | Automatic flag for items with stock of 5 units or fewer |

---

## Technology Stack

| Category | Technology |
|---|---|
| **Backend** | Go 1.23+ (standard library: `net/http`, `encoding/json`, `html/template`, `database/sql`, `embed`) |
| **Database** | PostgreSQL via [Supabase](https://supabase.com/) |
| **Database Driver** | [`github.com/lib/pq`](https://github.com/lib/pq) v1.12.3 |
| **Frontend** | HTML5, CSS3 (Vanilla CSS), JavaScript |
| **Charts** | [Chart.js](https://www.chartjs.org/) (bundled locally) |
| **Icons** | [Lucide Icons](https://lucide.dev/) (bundled locally) |
| **Deployment** | [Vercel](https://vercel.com/) (Go Serverless Function) |

---

## Project Structure

```
reshier-app/
├── main.go                     # Local dev entry point & routing
├── go.mod                      # Go module definition (reshier)
├── go.sum                      # Dependency checksum (github.com/lib/pq)
├── vercel.json                 # Vercel build & routing configuration
├── .vercelignore               # Files excluded from Vercel deployment
│
├── api/                        # Vercel serverless function entry point
│   ├── index.go                # Handler() — single entry point for all requests
│   ├── views/                  # Embedded HTML templates (go:embed)
│   │   ├── index.html
│   │   ├── barang.html
│   │   ├── tambah_barang.html
│   │   ├── edit_barang.html
│   │   ├── transaksi.html
│   │   ├── tambah_transaksi.html
│   │   ├── detail_transaksi.html
│   │   └── laporan.html
│   └── static/                 # Embedded static assets (go:embed)
│       ├── css/
│       ├── js/
│       └── images/
│
├── controllers/                # HTTP handlers — business logic per page
│   ├── dashboard.go            # Dashboard handler (/)
│   ├── barang.go               # Item CRUD handler (/barang/*)
│   ├── transaksi.go            # Transaction handler (/transaksi/*)
│   ├── laporan.go              # Report handler (/laporan)
│   └── templates.go            # Template parser (embed.FS / OS filesystem)
│
├── models/                     # Data struct definitions & database access
│   ├── barang.go               # Barang struct & DB CRUD functions
│   ├── transaksi.go            # Transaksi, ItemTransaksi structs & DB functions
│   ├── laporan.go              # BarangLaris struct (for best-seller reports)
│   └── storage.go              # InitDB(), GetDB(), createTables() — PostgreSQL
│
├── utils/                      # Helper functions
│   ├── search.go               # Search algorithms (Sequential & Binary Search)
│   ├── sort.go                 # Sorting algorithms (Selection, Insertion, Bubble Sort)
│   └── template.go             # Template FuncMap (formatRupiah, inc, shortDate)
│
├── views/                      # HTML templates for local development
│   ├── index.html
│   ├── barang.html
│   ├── tambah_barang.html
│   ├── edit_barang.html
│   ├── transaksi.html
│   ├── tambah_transaksi.html
│   ├── detail_transaksi.html
│   └── laporan.html
│
└── static/                     # Static assets for local development
    ├── css/
    │   └── styles.css
    ├── js/
    │   ├── chart.umd.min.js    # Chart.js library (bundled locally)
    │   ├── lucide.min.js       # Lucide Icons library (bundled locally)
    │   └── sidebar.js          # Sidebar navigation toggle logic
    └── images/
        └── logo.png
```

> **Note:** The `views/` and `static/` directories exist in two places:
> - Root level — used during **local development** (read from OS filesystem)
> - Inside `api/` — used in **Vercel deployment** (embedded into binary via `go:embed`)

---

## Application Architecture

```
Browser (HTTP Request)
        │
        ▼
 ┌────────────────────────────────────────────┐
 │  Entry Point                               │
 │  ├── main.go          (local dev)          │
 │  └── api/index.go     (Vercel serverless)  │
 └────────────────────────────────────────────┘
        │
        ▼
  controllers/
  ├── dashboard.go  →  views/index.html
  ├── barang.go     →  views/barang.html
  │                    views/tambah_barang.html
  │                    views/edit_barang.html
  ├── transaksi.go  →  views/transaksi.html
  │                    views/tambah_transaksi.html
  │                    views/detail_transaksi.html
  └── laporan.go    →  views/laporan.html
        │
        ▼
     models/                          utils/
  ├── storage.go  (InitDB, GetDB)  ├── CariBarangSequential()  (Sequential Search)
  ├── barang.go   (DB CRUD)        ├── CariBarangBinary()      (Binary Search)
  ├── transaksi.go (DB CRUD)       ├── CariBarangByNama()      (case-insensitive filter)
  └── laporan.go   (BarangLaris)   ├── FilterTransaksiByTime() (filter by time)
         │                         ├── UrutkanKodeBarang()     (Selection Sort)
         ▼                         ├── UrutkanHargaBarang()    (Insertion Sort)
   PostgreSQL (Supabase)           ├── UrutkanStokBarang()     (Insertion Sort)
   ├── barang                      └── BarangTerlaris()        (Bubble Sort + tally)
   ├── transaksi
   └── item_transaksi
```

---

## Database Schema

The application uses **PostgreSQL** hosted on [Supabase](https://supabase.com/). Tables are created automatically on first connection via `models.InitDB()`.

### `barang`

| Column | Type | Constraint |
|---|---|---|
| `kode` | `TEXT` | `PRIMARY KEY` |
| `nama` | `TEXT` | `NOT NULL` |
| `kategori` | `TEXT` | `NOT NULL DEFAULT ''` |
| `harga` | `INTEGER` | `NOT NULL DEFAULT 0` |
| `stok` | `INTEGER` | `NOT NULL DEFAULT 0` |

### `transaksi`

| Column | Type | Constraint |
|---|---|---|
| `id` | `TEXT` | `PRIMARY KEY` |
| `waktu` | `TIMESTAMPTZ` | `NOT NULL` |
| `total` | `INTEGER` | `NOT NULL DEFAULT 0` |
| `bayar` | `INTEGER` | `NOT NULL DEFAULT 0` |
| `kembalian` | `INTEGER` | `NOT NULL DEFAULT 0` |

### `item_transaksi`

| Column | Type | Constraint |
|---|---|---|
| `id` | `SERIAL` | `PRIMARY KEY` |
| `transaksi_id` | `TEXT` | `NOT NULL REFERENCES transaksi(id) ON DELETE CASCADE` |
| `kode_barang` | `TEXT` | `NOT NULL` |
| `nama_barang` | `TEXT` | `NOT NULL` |
| `harga` | `INTEGER` | `NOT NULL DEFAULT 0` |
| `jumlah` | `INTEGER` | `NOT NULL DEFAULT 0` |
| `subtotal` | `INTEGER` | `NOT NULL DEFAULT 0` |

---

## Data Models

### `models.Barang`

```go
type Barang struct {
    Kode       string  // Unique item code, e.g. "NSG001"
    Nama       string  // Item name
    Kategori   string  // Category, e.g. "Makanan" (Food), "Minuman" (Drink)
    Harga      int     // Unit price (in Indonesian Rupiah)
    Stok       int     // Available stock quantity
    StokKritis bool    // true if Stok <= 5 (computed, not stored in DB)
}
```

### `models.ItemTransaksi`

```go
type ItemTransaksi struct {
    KodeBarang string  // Item code at time of purchase
    NamaBarang string  // Item name at time of purchase
    Harga      int     // Unit price at time of purchase
    Jumlah     int     // Number of units purchased
    Subtotal   int     // Harga * Jumlah
}
```

### `models.Transaksi`

```go
type Transaksi struct {
    ID        string           // Format: "TRX-YYYYMMDD-NNN"
    Waktu     time.Time        // Transaction timestamp
    Items     []ItemTransaksi  // List of purchased items
    Total     int              // Total price of all items
    Bayar     int              // Amount tendered by the customer
    Kembalian int              // Change returned (Bayar - Total)
}
```

### `models.BarangLaris`

```go
type BarangLaris struct {
    KodeBarang      string  // Item code
    NamaBarang      string  // Item name
    TotalTerjual    int     // Cumulative units sold
    TotalPendapatan int     // Cumulative revenue from this item
}
```

---

## Routing & Endpoints

| Method | Path | Handler | Description |
|---|---|---|---|
| `GET` | `/` | `Dashboard` | Main dashboard & summary page |
| `GET` | `/barang` | `TampilkanBarang` | Item list (supports `?sort=`, `?q=`, `?kategori=`) |
| `GET/POST` | `/barang/tambah` | `TambahBarang` | Add new item form & action |
| `GET/POST` | `/barang/edit` | `EditBarang` | Edit item form & action (`?kode=`) |
| `GET` | `/barang/hapus` | `HapusBarang` | Delete item by code (`?kode=`) |
| `GET` | `/barang/cari` | `CariBarangJSON` | **JSON API** — live item search (`?q=`) |
| `GET` | `/transaksi` | `TampilkanTransaksi` | Transaction list (supports `?q=`) |
| `GET/POST` | `/transaksi/tambah` | `TambahTransaksi` | New transaction form & action |
| `GET` | `/transaksi/detail` | `DetailTransaksi` | Transaction receipt / detail (`?id=`) |
| `GET` | `/laporan` | `LaporanHarian` | Full reports page |
| `GET` | `/static/` | `FileServer` | Serves static files (CSS, JS, images) |

### Query Parameters for `/barang`

| Parameter | Value | Description |
|---|---|---|
| `sort` | `kode-asc`, `kode-desc` | Sort by item code |
| `sort` | `harga-asc`, `harga-desc` | Sort by price |
| `sort` | `stok-asc`, `stok-desc` | Sort by stock level |
| `q` | `<keyword>` | Search items (name, code, or category) |
| `kategori` | `<category name>` | Filter by category |

### Query Parameters for `/transaksi`

| Parameter | Format | Example |
|---|---|---|
| `q` | `DD-MM-YYYY` | `?q=21-07-2026` |
| `q` | `MM-YYYY` | `?q=07-2026` |
| `q` | `YYYY` | `?q=2026` |
| `q` | `HH:MM` | `?q=20:54` |

---

## Algorithms Used

This project implements sorting and searching algorithms manually as a learning demonstration:

### Searching (`utils/search.go`)

| Function | Algorithm | Purpose |
|---|---|---|
| `CariBarangSequential(kode)` | **Sequential Search** | Finds an item by exact code match; used in transaction item lookup |
| `CariBarangBinary(kode)` | **Binary Search** | Alternative code lookup (requires data to be sorted in ascending order) |
| `CariBarangByNama(query)` | Linear filter | Case-insensitive search across item name, code, and category |
| `FilterTransaksiByTime(query)` | Linear filter | Filters transactions by time precision (minute, hour, day, month, or year) |

### Sorting (`utils/sort.go`)

| Function | Algorithm | Complexity |
|---|---|---|
| `UrutkanKodeBarang(asc)` | **Selection Sort** | O(n²) |
| `UrutkanHargaBarang(asc)` | **Insertion Sort** | O(n²) worst-case, O(n) best-case |
| `UrutkanStokBarang(asc)` | **Insertion Sort** | O(n²) worst-case, O(n) best-case |
| `BarangTerlaris(topN)` | **Bubble Sort** (descending) | O(n²) — applied after building the tally map |

> **Note:** Sorting is applied in-memory after data is fetched from the database. The database default sort order is `ORDER BY kode ASC` for items and `ORDER BY waktu DESC` for transactions.

---

## Getting Started

### Prerequisites

- **Go 1.23** or later — run `go version` to verify
- **PostgreSQL** database (or a free [Supabase](https://supabase.com/) project)

```bash
go version
# go version go1.23.x ...
```

### Steps

**1. Clone the repository**

```bash
git clone https://github.com/M4VAIS3/reshier-app.git
cd reshier-app
```

**2. Set the database connection string**

```bash
# Linux / macOS
export DATABASE_URL="postgresql://user:password@host:port/dbname?sslmode=require"

# Windows (PowerShell)
$env:DATABASE_URL = "postgresql://user:password@host:port/dbname?sslmode=require"
```

**3. Install dependencies**

```bash
go mod download
```

**4. Start the server**

```bash
go run main.go
```

**5. Open in your browser**

```
http://localhost:8080
```

Expected terminal output:

```
✅ Database terhubung
🚀 Server berjalan di http://localhost:8080
```

> **Note:** Tables (`barang`, `transaksi`, `item_transaksi`) are created automatically on first connection if they do not already exist.

---

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `DATABASE_URL` | **Yes** | PostgreSQL connection string (e.g., Supabase connection pooler URL) |
| `PORT` | No | Server port for local dev (default: `8080`) |

### Supabase Connection String

You can find your `DATABASE_URL` in the Supabase dashboard under **Project Settings → Database → Connection string → URI**.

Example format:
```
postgresql://postgres.[project-ref]:[password]@aws-0-[region].pooler.supabase.com:6543/postgres?sslmode=require
```

---

## Deployment

### Vercel (Production)

The application is deployed on [Vercel](https://vercel.com/) as a Go serverless function. The configuration is defined in [`vercel.json`](vercel.json):

```json
{
  "builds": [
    {
      "src": "api/index.go",
      "use": "@vercel/go"
    }
  ],
  "routes": [
    { "src": "/(.*)", "dest": "/api/index.go" }
  ]
}
```

**How it works:**
1. All HTTP requests are routed to `api/index.go` via `Handler()` function
2. HTML templates and static assets are embedded into the binary using `go:embed` directives
3. The `controllers.ViewsFS` variable is set to the embedded filesystem at init time
4. The `DATABASE_URL` environment variable must be configured in Vercel project settings

**Deploy steps:**
1. Push to your GitHub repository
2. Import the project in [Vercel](https://vercel.com/)
3. Add `DATABASE_URL` to the Vercel project's Environment Variables
4. Vercel automatically builds and deploys on each push

### Local Development

For local development, templates and static files are read directly from the OS filesystem (no embedding required). Set `DATABASE_URL` as an environment variable and run `go run main.go`.

---

## Template Helpers

Custom functions available in all HTML templates (`utils/template.go`):

| Function | Description | Example |
|---|---|---|
| `formatRupiah` | Formats an integer as Indonesian Rupiah with dot separators | `15000` → `15.000` |
| `inc` | Increments an integer by one (used for row numbering in tables) | `inc 0` → `1` |
| `shortDate` | Trims whitespace from a date string | `" 21 July 2026 "` → `"21 July 2026"` |

---

## Limitations

- **No authentication** — anyone with access to the URL can use the application
- **Single-user** — there is no concurrent write safety mechanism beyond database-level transactions (not ideal for multiple cashiers simultaneously)
- **No automatic backups** — relies on Supabase's built-in backup features

---

## Licence

This project was created for personal portfolio purposes. However, you are free to use it as long as you credit me as the original author.

---
