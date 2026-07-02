# Database Schema Design

## Sistem Point of Sales (POS) Coffee Shop Skala Menengah

Dokumen ini menjelaskan desain database MySQL untuk sistem POS coffee shop yang didasarkan pada dokumen `PRD.md`. Arsitektur skema dirancang untuk memberikan konsistensi data yang tinggi, audit trail yang jelas, serta performa maksimal dengan mengimplementasikan standar teknis yang ketat.

---

### 1. Daftar Tabel dan Deskripsi Singkat

| Nama Tabel        | Deskripsi Singkat                                                                                 |
| :---------------- | :------------------------------------------------------------------------------------------------ |
| `users`           | Menyimpan informasi kredensial, profil, dan peran pengguna (Owner dan Cashier).                   |
| `categories`      | Menyimpan data pengelompokan produk (kategori menu).                                              |
| `products`        | Menyimpan katalog menu utama coffee shop beserta status dan harganya.                             |
| `tables`          | Menyimpan daftar meja beserta pelacakan status ketersediaannya.                                   |
| `promos`          | Menyimpan aturan, parameter potongan harga, dan validitas kode promo.                             |
| `shifts`          | Mencatat siklus kerja kasir, pelacakan modal awal, kas akhir, dan selisih kas.                    |
| `orders`          | Menyimpan header transaksi penjualan, status pesanan, meja yang digunakan, serta rekap harga.     |
| `order_items`     | Menyimpan detail item produk yang dibeli di dalam setiap transaksi (tabel persimpangan/junction). |
| `payments`        | Menyimpan riwayat pembayaran terintegrasi dengan payment gateway Midtrans Snap.                   |
| `stock_movements` | Mencatat log mutasi/riwayat pergerakan keluar masuknya stok produk untuk audit trail.             |

---

### 2. Struktur Kolom per Tabel

#### 2.1. Tabel: `users` 0

Menyimpan data pengguna sistem (Owner & Cashier).

- **nama_kolom** | **tipe_data** | **constraint** | **keterangan**
- `id` | VARCHAR(36) | PRIMARY KEY | Unique ID berbasis UUID v4.
- `name` | VARCHAR(100) | NOT NULL | Nama lengkap pengguna.
- `email` | VARCHAR(100) | NOT NULL, UNIQUE | Email untuk login.
- `password` | VARCHAR(255) | NOT NULL | Password terenkripsi menggunakan Bcrypt.
- `role` | ENUM('OWNER', 'CASHIER') | NOT NULL | Peran pengguna dalam sistem.
- `is_active` | BOOLEAN | NOT NULL, DEFAULT TRUE | Status keaktifan akun kasir/owner.
- `created_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu pencatatan akun dibuat.
- `updated_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu pembaruan data terakhir.
- `deleted_at` | TIMESTAMP | NULLABLE | Timestamp untuk mendukung mekanisme Soft Delete.

#### 2.2. Tabel: `categories`

Menyimpan kategori produk/menu.

- **nama_kolom** | **tipe_data** | **constraint** | **keterangan**
- `id` | VARCHAR(36) | PRIMARY KEY | Unique ID berbasis UUID v4.
- `name` | VARCHAR(50) | NOT NULL, UNIQUE | Nama kategori (contoh: 'Coffee', 'Pastry').
- `created_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu pencatatan dibuat.
- `updated_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu pembaruan data terakhir.

#### 2.3. Tabel: `products` 0

Menyimpan katalog menu masakan/minuman coffee shop.

- **nama_kolom** | **tipe_data** | **constraint** | **keterangan**
- `id` | VARCHAR(36) | PRIMARY KEY | Unique ID berbasis UUID v4.
- `category_id` | VARCHAR(36) | NOT NULL | FOREIGN KEY | Relasi ke tabel `categories`.
- `sku` | VARCHAR(50) | NOT NULL, UNIQUE | Stock Keeping Unit sebagai kode unik produk.
- `name` | VARCHAR(100) | NOT NULL | Nama menu produk.
- `description` | TEXT | NULLABLE | Deskripsi detail isi/komposisi produk.
- `price` | BIGINT | NOT NULL, DEFAULT 0 | Harga jual dasar produk dalam satuan sen (Rupiah x 100).
- `current_stock` | INT | NOT NULL, DEFAULT 0 | Kuantitas stok siap jual saat ini.
- `image_url` | VARCHAR(255) | NULLABLE | Path atau URL lokasi penyimpanan foto produk.
- `is_available` | BOOLEAN | DEFAULT TRUE | Ketersediaan barang
- `created_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu pencatatan produk dibuat.
- `updated_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu pembaruan data terakhir.
- `deleted_at` | TIMESTAMP | NULLABLE | Timestamp untuk mendukung mekanisme Soft Delete.

#### 2.4. Tabel: `tables`

Menyimpan nomor atau penamaan meja di coffee shop.

- **nama_kolom** | **tipe_data** | **constraint** | **keterangan**
- `id` | VARCHAR(36) | PRIMARY KEY | Unique ID berbasis UUID v4.
- `table_number` | VARCHAR(20) | NOT NULL, UNIQUE | Penomoran atau nama meja (contoh: 'Meja 01', 'Area VIP-A').
- `status` | ENUM('AVAILABLE', 'OCCUPIED') | NOT NULL, DEFAULT 'AVAILABLE' | Menunjukkan keterisian meja saat ini.
- `created_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu pencatatan dibuat.
- `updated_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu pembaruan data terakhir.

#### 2.5. Tabel: `promos`

Menyimpan data master promo diskon.

- **nama_kolom** | **tipe_data** | **constraint** | **keterangan**
- `id` | VARCHAR(36) | PRIMARY KEY | Unique ID berbasis UUID v4.
- `code` | VARCHAR(50) | NOT NULL, UNIQUE | Kode promo unik yang diinput kasir (contoh: 'JUMATBERKAH').
- `name` | VARCHAR(100) | NOT NULL | Nama program promosi.
- `promo_type` | ENUM('PERCENTAGE', 'NOMINAL') | NOT NULL | Jenis tipe kalkulasi pemotongan harga.
- `value` | BIGINT | NOT NULL | Nilai diskon (jika PERCENTAGE disimpan 10 untuk 10%, jika NOMINAL dalam sen, misal: 500000 untuk Rp5.000).
- `min_purchase` | BIGINT | NOT NULL, DEFAULT 0 | Syarat minimal belanja sebelum diskon dalam satuan sen.
- `max_discount` | BIGINT | NOT NULL, DEFAULT 0 | Batas maksimal potongan diskon (terutama untuk tipe persentase) dalam sen.
- `start_date` | TIMESTAMP | NOT NULL | Batas awal periode berlakunya promo.
- `end_date` | TIMESTAMP | NOT NULL | Batas akhir periode berlakunya promo.
- `is_active` | BOOLEAN | NOT NULL, DEFAULT TRUE | Status keaktifan master diskon oleh owner.
- `created_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Waktu data dibuat.
- `updated_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu pembaruan data terakhir.

#### 2.6. Tabel: `shifts` 0

Menyimpan pencatatan shift kasir harian.

- **nama_kolom** | **tipe_data** | **constraint** | **keterangan**
- `id` | VARCHAR(36) | PRIMARY KEY | Unique ID berbasis UUID v4.
- `cashier_id` | VARCHAR(36) | NOT NULL | FOREIGN KEY | ID Kasir pelaksana shift, relasi ke tabel `users`.
- `opening_cash` | BIGINT | NOT NULL |
- `closing_cash` | BIGINT | NULLABLE |
- `total_sales` | BIGINT | DEFAULT 0 |
- `status` | ENUM('open', 'closed') | NOT NULL, DEFAULT 'open' | Status status aktivitas operasional shift.
- `opened_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
- `closed_at` | TIMESTAMP | NULLABLE |

#### 2.7. Tabel: `orders` 0

Menyimpan ringkasan utama (header) dari transaksi penjualan.

- **nama_kolom** | **tipe_data** | **constraint** | **keterangan**
- `id` | VARCHAR(36) | PRIMARY KEY | Unique ID berbasis UUID v4.
- `shift_id` | VARCHAR(36) | NOT NULL | FOREIGN KEY | ID Shift transaksi berjalan, relasi ke tabel `shifts`.
- `table_id` | VARCHAR(36) | FOREIGN KEY, NULLABLE | ID Meja terikat, bernilai NULL jika tipe pesanan 'TAKEAWAY'.
- `cashier_id` | VARCHAR(36) | NOT NULL | FOREIGN KEY | ID Kasir pembuat pesanan, relasi ke tabel `users`.
- `promo_id` | VARCHAR(36) | FOREIGN KEY, NULLABLE | ID Promo yang diaplikasikan, bernilai NULL jika tanpa promo.
- `status` | ENUM('draft', 'pending_payment', 'paid', 'cancelled') | NOT NULL, DEFAULT 'PENDING' | Status alur hidup transaksi penjualan.
- `subtotal`| BIGINT | NOT NULL, DEFAULT 0 | Total harga seluruh item pesanan sebelum diskon/pajak (sen).
- `discount_amount`| BIGINT | NOT NULL, DEFAULT 0 | Total potongan harga dari promo yang diaplikasikan (sen).
- `total` | BIGINT | NOT NULL | Nilai bersih akhir yang wajib dibayar pelanggan (sen).
- `payment_method`| VARCHAR(50) | NULLABLE | Metode Pembayaran.
- `midtrans_order_id` | VARCHAR(100) | NULLABLE, UNIQUE | ID unik referensi pesanan yang dikirim ke sistem Midtrans. Digunakan untuk pencocokan data saat menerima webhook. Bernilai NULL jika pembayaran tunai/belum checkout.
- `created_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Tanggal/Waktu pesanan dibuat.
- `updated_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Tanggal/Waktu pembaruan status transaksi terakhir.

#### 2.8. Tabel: `order_items` 0

Menyimpan detail baris item yang dibeli di dalam transaksi (Hubungan Many-to-Many antara `orders` dan `products`).

- **nama_kolom** | **tipe_data** | **constraint** | **keterangan**
- `id` | VARCHAR(36) | PRIMARY KEY | Unique ID berbasis UUID v4.
- `order_id` | VARCHAR(36) | NOT NULL | FOREIGN KEY | Relasi ke tabel master header transaksi `orders`.
- `product_id` | VARCHAR(36) | FOREIGN KEY | Relasi ke tabel katalog `products`.
- `quantity` | INT | NOT NULL | Jumlah kuantitas barang yang dibeli.
- `price` | BIGINT | NOT NULL | Harga jual produk per unit saat transaksi terjadi (snapshot price, sen).
- `subtotal` | BIGINT | NOT NULL | Kalkulasi otomatis harga unit x kuantitas (sen).
- `notes` | VARCHAR(255) | NULLABLE | Catatan kustom item (contoh: 'Less sugar, extra ice').
- `created_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Log waktu pembuatan.
- `updated_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Log pembaruan data terakhir.

#### 2.9. Tabel: `payments` 0

Menyimpan data integrasi pembayaran digital via payment gateway Midtrans Snap.

- **nama_kolom** | **tipe_data** | **constraint** | **keterangan**
- `id` | VARCHAR(36) | PRIMARY KEY | Unique ID berbasis UUID v4.
- `order_id` | VARCHAR(36) | NOT NULL | FOREIGN KEY | Relasi ke tabel master transaksi `orders`.
- `midtrans_order_id` | VARCHAR (100) | NOT NULL |
- `midtrans_transaction_id` | VARCHAR (100) | NULLABLE |
- `status` | ENUM('pending', 'paid', 'failed', 'expired') | NOT NULL, DEFAULT 'PENDING' | Status pembayaran sinkronisasi webhook Midtrans.
- `amount` | BIGINT | NOT NULL | Total nilai nominal pembayaran yang ditagihkan (sen).
- `payment_type` | VARCHAR(50) | NULLABLE | Jenis pembayaran dari webhook (contoh: 'qris', 'gopay', 'bank_transfer').
- `raw_notification` | JSON | NULLABLE |
- `paid_at` | TIMESTAMP | NULLABLE | Log penanda waktu dari webhook saat pembayaran sukses/settlement.
- `created_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Log waktu order checkout diproses.
- `updated_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Log sinkronisasi webhook terupdate.

#### 2.10. Tabel: `stock_movements` 0

Menyimpan rekaman mutasi/riwayat pergerakan keluar masuknya stok produk (Audit Trail Log).

- **nama_kolom** | **tipe_data** | **constraint** | **keterangan**
- `id` | VARCHAR(36) | PRIMARY KEY | Unique ID berbasis UUID v4.
- `product_id` | VARCHAR(36) | | NOT NULL | FOREIGN KEY | Target produk yang mengalami mutasi stok, relasi ke `products`.
- `type` | ENUM('purchase', 'sale', 'adjustment', 'return') | NOT NULL | Jenis pemicu mutasi stok (Stok Manual Masuk, Rusak/Koreksi, atau Penjualan Kasir).
- `quantity` | INT | NOT NULL | Jumlah perubahan kuantitas (bernilai positif jika penambahan, negatif jika pengurangan).
- `notes` | TEXT | NULLABLE | Alasan tertulis penyesuaian (contoh: 'Bahan baku masuk', 'Susu kedaluwarsa').
- `created_by` | VARCHAR(36) | NOT NULL | FOREIGN KEY | ID User (Owner/Cashier) pemicu mutasi, relasi ke `users`.
- `created_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Log penanda mutasi stok terjadi.
- `updated_at` | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Log pembaruan data internal.

---

### 3. Relasi Antar Tabel (Foreign Keys)

Berikut adalah pemetaan foreign key relationships di database untuk menjamin integritas data referensial:

- `[users.id]` $\leftarrow$ `[orders.cashier_id]` | Menghubungkan kasir pelaksana transaksi dengan pesanan penjualan terkait.
- `[users.id]` $\leftarrow$ `[shifts.cashier_id]` | Menghubungkan kasir yang masuk/piket dengan log siklus pertanggungjawaban modal shift.
- `[users.id]` $\leftarrow$ `[stock_movements.created_by]` | Melacak aktor/staf (Owner/Cashier) bertanggung jawab di balik penyesuaian stok.
- `[categories.id]` $\leftarrow$ `[products.category_id]` | Mengelompokkan setiap produk menu ke dalam satu kategori master yang valid.
- `[products.id]` $\leftarrow$ `[order_items.product_id]` | Mengidentifikasi menu produk yang dibeli di dalam baris item transaksi penjualan.
- `[products.id]` $\leftarrow$ `[stock_movements.product_id]` | Mengaitkan riwayat catatan mutasi log keluar-masuk dengan produk spesifik.
- `[orders.id]` $\leftarrow$ `[order_items.order_id]` | Mengikat relasi baris-baris produk belanjaan ke dalam satu payung invoice pemesanan.
- `[orders.id]` $\leftarrow$ `[payments.order_id]` | Menghubungkan log data transaksi payment gateway Midtrans dengan tagihan pesanan lokal.
- `[shifts.id]` $\leftarrow$ `[orders.shift_id]` | Memasukkan transaksi penjualan kasir ke dalam akumulasi rekap shift kerja kasir yang aktif.
- `[tables.id]` $\leftarrow$ `[orders.table_id]` | Mengaitkan pesanan pelanggan dine-in dengan meja yang sedang mereka gunakan (nullable).
- `[promos.id]` $\leftarrow$ `[orders.promo_id]` | Mencatat penggunaan kode promosi khusus yang digunakan untuk memotong total invoice (nullable).

---

### 4. Keputusan Desain Penting (Design Decisions)

- **Penggunaan UUID v4 (VARCHAR 36) dibanding AUTO_INCREMENT:**
  1.  _Keamanan:_ ID berurutan (`1, 2, 3...`) sangat rentan terhadap serangan kerentanan _Insecure Direct Object Reference (IDOR)_, di mana kompetitor/penyerang dapat menebak ID transaksi atau data user dengan mudah hanya lewat manipulasi URL parameter.
  2.  _Skalabilitas:_ Penggunaan UUID di level kode aplikasi (Go) memungkinkan pembuatan ID yang unik secara universal bahkan sebelum data dimasukkan ke dalam MySQL. Ini sangat berguna jika ke depan coffee shop melakukan ekspansi menggunakan strategi replikasi database terdistribusi (_multi-master/sharding_) tanpa risiko bentrok ID (_collision_).
- **Penyimpanan Uang Menggunakan BIGINT Sen (Rupiah Terkecil):**
  Penggunaan tipe data `FLOAT` atau `DOUBLE` dilarang keras untuk data finansial karena karakteristik _floating-point arithmetic_ pada arsitektur komputer yang tidak presisi (menghasilkan pembulatan desimal yang salah seperti `0.1 + 0.2 = 0.30000000000000004`). Penggunaan `DECIMAL` sebenarnya aman, namun operasi matematika di level bahasa backend Go jauh lebih cepat dan efisien jika memproses data bertipe integer (`BIGINT` di database, `int64` di Go). Nilai uang disimpan dengan mengalikannya dengan 100 (contoh: harga kopi Rp25.000 disimpan sebagai `2500000` sen).
- **Mekanisme Soft Delete (`deleted_at`):**
  Tabel vital seperti `users` dan `products` menggunakan pendekatan _Soft Delete_. Jika Owner menghapus menu "Espresso" dari katalog, sistem tidak boleh menghapusnya secara fisik dengan `DELETE FROM`. Karena jika dihapus fisik, seluruh data transaksi historis di masa lalu yang pernah menjual "Espresso" akan mengalami kegagalan _Foreign Key Constraint_ atau menghasilkan pelaporan yang korup. Kolom `deleted_at` yang terisi menandakan item disembunyikan dari kasir, namun tetap valid untuk laporan keuangan historis.
- **Penerapan Snapshot Harga (`unit_price` di `order_items`):**
  Harga jual produk di tabel `products` bersifat dinamis dan dapat diubah Owner sewaktu-waktu akibat inflasi bahan baku. Oleh karena itu, sistem wajib menyalin (mengambil _snapshot_) harga produk terkini pada saat transaksi sukses dilakukan ke dalam kolom `unit_price` di tabel `order_items`. Jika di kemudian hari harga "Americano" naik dari Rp20.000 menjadi Rp22.000, laporan pendapatan di bulan-bulan sebelumnya tidak akan ikut berubah/rusak.
- **Penggunaan ENUM untuk Kolom Status Dinamis:**
  Status seperti `role`, `status` produk, status meja, status shift, status order, dan status payment diatur menggunakan tipe data `ENUM`. Ini bertindak sebagai _data-integrity guard_ di level MySQL, mencegah masuknya teks sampah atau status tidak valid (invalid state) yang berada di luar cakupan aturan bisnis (Business Rules) aplikasi backend Go.
