# Product Requirement Document (PRD)
## Sistem Point of Sales (POS) Coffee Shop Skala Menengah

### 1. Deskripsi Produk
Sistem Point of Sales (POS) ini adalah platform manajemen penjualan dan operasional harian yang dirancang khusus untuk bisnis coffee shop skala menengah. Sistem ini bertujuan untuk mendigitalisasi proses pemesanan, pengelolaan stok, integrasi pembayaran cashless, serta menyediakan pelaporan keuangan yang akurat untuk pemilik bisnis. 

Sistem ini dibangun menggunakan backend berbasis **Go (Golang)** untuk performa tinggi dan konkurensi yang andal, serta **MySQL** sebagai sistem manajemen database relasional untuk menjaga integritas data transaksi.

**Target Pengguna:**
* **Owner (Pemilik Coffee Shop):** Pengguna yang membutuhkan visibilitas penuh terhadap performa bisnis, laporan keuangan, kontrol inventaris, dan manajemen promosi.
* **Cashier (Kasir/Barista):** Pengguna di lini depan yang membutuhkan antarmuka cepat, responsif, dan andal untuk memproses pesanan pelanggan dan mengelola uang kas (shift).

---

### 2. User Roles
Sistem ini membagi hak akses ke dalam dua peran utama dengan batasan otoritas yang ketat:

#### A. Owner (Pemilik)
* Memiliki akses penuh (Full CRUD) ke semua modul data master: Produk, Kategori, Meja, dan User Cashier.
* Dapat melihat dashboard performa bisnis dan mengunduh laporan keuangan/penjualan.
* Mengontrol manajemen stok dan melihat riwayat pergerakan stok secara menyeluruh.
* Mengonfigurasi promo dan diskon yang berlaku di outlet.
* *Tidak memiliki akses untuk melakukan transaksi kasir harian kecuali membuat akun kasir untuk dirinya sendiri.*

#### B. Cashier (Kasir)
* Memiliki akses operasional terbatas untuk modul transaksi dan kasir.
* Wajib melakukan manajemen shift (buka dan tutup shift) untuk mencatat akuntabilitas kas.
* Dapat membuat pesanan baru, memilih meja, menerapkan promo, dan memproses pembayaran via payment gateway.
* Dapat melihat riwayat transaksi yang terjadi terbatas pada shift yang sedang berjalan.
* *Tidak memiliki akses ke dashboard laporan pemilik, pengaturan produk, modifikasi stok manual, atau konfigurasi promo.*

---

### 3. Fitur Owner

#### 3.1. Manajemen Produk
Fitur untuk mengelola seluruh katalog produk/menu yang dijual di coffee shop.
* **Fungsionalitas:**
    * Membuat produk baru (Nama, SKU, Deskripsi, Harga, Kategori, Foto Produk, Status Aktif/Nonaktif).
    * Membaca/Menampilkan daftar produk dengan fitur pencarian dan filter per kategori.
    * Mengubah data produk yang sudah ada.
    * Menghapus produk menggunakan mekanisme *Soft Delete* (data tetap ada di database untuk menjaga integritas data transaksi masa lalu, namun tidak muncul di aplikasi kasir).
* **Acceptance Criteria:**
    * GIVEN Owner berada di halaman manajemen produk, WHEN Owner menambahkan produk baru dengan data valid dan mengunggah foto, THEN sistem menyimpan data ke database dan menampilkan pesan sukses.
    * GIVEN Owner memilih untuk menghapus produk, WHEN tombol hapus diklik, THEN sistem mengubah kolom `deleted_at` di MySQL dan menyembunyikannya dari katalog kasir tanpa merusak data transaksi historis yang mengandung produk tersebut.
    * GIVEN Owner mengubah status produk menjadi 'Nonaktif', WHEN kasir membuka menu, THEN produk tersebut tidak dapat dipilih untuk transaksi baru.

#### 3.2. Manajemen Kategori Produk
Fitur untuk mengelompokkan produk guna mempermudah pencarian menu (misal: Coffee, Non-Coffee, Pastry, Heavy Meal).
* **Fungsionalitas:** CRUD (Create, Read, Update, Delete) kategori produk.
* **Acceptance Criteria:**
    * GIVEN Owner membuat kategori baru, WHEN nama kategori unik dimasukkan, THEN kategori berhasil disimpan.
    * GIVEN Owner menghapus kategori yang masih memiliki produk aktif di dalamnya, THEN sistem menolak penghapusan dan menampilkan pesan error agar Owner memindahkan produk tersebut terlebih dahulu.

#### 3.3. Manajemen Stok
Fitur untuk memantau kuantitas stok produk siap jual dan mencatat setiap perubahan yang terjadi.
* **Fungsionalitas:**
    * Menampilkan jumlah stok terkini untuk setiap produk.
    * Menambah stok manual (Stock In) untuk pasokan baru atau mengurangi stok (Stock Out/Adjustment) jika ada produk rusak/kedaluwarsa.
    * Mencatat otomatis riwayat pergerakan stok mencakup: Waktu, Jumlah, Tipe Pergerakan (Stock In, Stock Out, Sales), Keterangan, dan User ID pelaksana.
* **Acceptance Criteria:**
    * GIVEN Owner melakukan penyesuaian stok manual, WHEN jumlah baru dimasukkan beserta alasan penyesuaian, THEN stok produk terupdate di MySQL dan entri baru terbentuk di tabel riwayat pergerakan stok.
    * GIVEN Transaksi kasir berhasil dikonfirmasi, THEN sistem secara otomatis mengurangi stok produk terkait dan mencatat tipe pergerakan 'Sales' di tabel riwayat stok.

#### 3.4. Manajemen Meja
Fitur untuk mendaftarkan nomor atau nama meja yang tersedia di coffee shop guna pencatatan pesanan *dine-in*.
* **Fungsionalitas:** CRUD meja dan pelacakan status meja (Tersedia / Terisi).
* **Acceptance Criteria:**
    * GIVEN Owner menambahkan meja baru, WHEN nomor meja diinput, THEN sistem memvalidasi agar tidak ada nomor meja ganda.
    * GIVEN Kasir memilih meja untuk transaksi baru, THEN status meja otomatis berubah menjadi 'Terisi' hingga transaksi diselesaikan atau dibayarkan.

#### 3.5. Manajemen User Cashier
Fitur bagi pemilik untuk membuat dan mengelola akun masuk bagi staf kasir.
* **Fungsionalitas:** CRUD user kasir (Nama, Email/Username, Password terenkripsi, Status Aktif/Nonaktif).
* **Acceptance Criteria:**
    * GIVEN Owner membuat akun kasir baru, WHEN password diinput, THEN sistem wajib melakukan hashing password menggunakan Bcrypt sebelum disimpan ke database MySQL.
    * GIVEN Akun kasir dinonaktifkan oleh Owner, WHEN kasir tersebut mencoba login, THEN sistem menolak akses login dengan pesan error yang sesuai.

#### 3.6. Dashboard Laporan
Halaman utama bagi Owner untuk melihat visualisasi ringkasan performa penjualan coffee shop.
* **Fungsionalitas:**
    * Menampilkan total revenue dalam rentang waktu harian, mingguan, dan bulanan.
    * Menampilkan daftar produk terlaris (*Top Selling Products*).
    * Menampilkan ringkasan total transaksi dan performa penjualan per kasir.
* **Acceptance Criteria:**
    * GIVEN Owner membuka dashboard, WHEN memilih filter waktu tertentu, THEN sistem melakukan agregasi data transaksi dari database MySQL secara real-time dan menampilkan metrik revenue dengan benar.

#### 3.7. Manajemen Promo
Fitur untuk merencanakan dan mengaktifkan program potongan harga guna menarik pelanggan.
* **Fungsionalitas:**
    * Membuat promo dengan tipe: Persentase (misal: 10%) atau Nominal (misal: Rp 15.000).
    * Mengatur parameter promo: Kode Promo, Nama Promo, Minimum Pembelian, Maksimum Potongan, Tanggal Mulai, Tanggal Berakhir, dan Status Aktif.
* **Acceptance Criteria:**
    * GIVEN Owner membuat promo baru, WHEN parameter diisi lengkap dan disimpan, THEN kode promo tersebut tersedia untuk digunakan oleh kasir sesuai periode tanggal berlakunya.

#### 3.8. Export CSV
Fitur untuk mengunduh data laporan penjualan mentah untuk keperluan analisis akuntansi eksternal.
* **Fungsionalitas:** Mengonversi data transaksi berdasarkan filter tanggal menjadi file berformat `.csv` yang dapat diunduh langsung.
* **Acceptance Criteria:**
    * GIVEN Owner mengklik tombol "Export CSV" pada modul laporan, WHEN rentang tanggal valid dipilih, THEN backend Go menghasilkan stream data atau file CSV yang berisi detail transaksi lengkap dan mengunduhnya ke perangkat pengguna.

---

### 4. Fitur Cashier

#### 4.1. Shift Management
Fitur untuk memastikan akuntabilitas arus kas masuk dan keluar di laci kasir (*cash drawer*) pada setiap pergantian staf.
* **Fungsionalitas:**
    * **Buka Shift:** Kasir wajib menginput jumlah modal kas awal (uang tunai di laci) sebelum mulai melayani transaksi.
    * **Tutup Shift:** Kasir memasukkan jumlah uang kas akhir aktual di laci saat shift selesai. Sistem akan menghitung rekapitulasi otomatis (Modal Awal + Total Penjualan Tunai vs Uang Kas Aktual) dan mencatat selisih (*variance*) jika ada.
* **Acceptance Criteria:**
    * GIVEN Kasir baru saja login ke sistem, WHEN belum melakukan Buka Shift, THEN seluruh menu pembuatan transaksi terkunci dan mengarahkan kasir ke halaman Buka Shift.
    * GIVEN Kasir melakukan Tutup Shift, WHEN memasukkan jumlah uang tunai akhir, THEN sistem mengunci shift tersebut, mencatat waktu penutupan, dan menghasilkan ringkasan rekap shift yang tidak dapat diubah kembali.

#### 4.2. Buat Transaksi
Fitur utama kasir untuk mencatat pesanan pelanggan di coffee shop.
* **Fungsionalitas:**
    * Memilih nomor meja (untuk *dine-in*) atau memilih opsi *takeaway*.
    * Memilih produk dari katalog, mengatur kuantitas (*quantity*), serta menambahkan catatan kustom opsional per item.
    * Menghitung otomatis subtotal, pajak, diskon, dan total akhir pesanan.
* **Acceptance Criteria:**
    * GIVEN Kasir memilih beberapa produk, WHEN kuantitas diubah, THEN sistem memperbarui subtotal secara real-time di layar kasir.
    * GIVEN Kasir mengklik simpan pesanan, WHEN data valid, THEN sistem membuat record transaksi baru di MySQL dengan status awal 'Pending'.

#### 4.3. Apply Promo
Fitur untuk menerapkan potongan harga ke dalam transaksi yang sedang aktif sebelum proses checkout.
* **Fungsionalitas:** Memasukkan atau memilih kode promo aktif yang memenuhi syarat minimum pembelian.
* **Acceptance Criteria:**
    * GIVEN Kasir menerapkan kode promo, WHEN total transaksi memenuhi minimum pembelian dan promo masih aktif, THEN sistem memotong total biaya sesuai aturan diskon (nominal/persentase) dan memperbarui total akhir pembayaran.

#### 4.4. Checkout Midtrans
Integrasi dengan payment gateway Midtrans Snap untuk memproses pembayaran digital yang aman secara real-time.
* **Fungsionalitas:**
    * Mengirim data transaksi ke API Midtrans dari backend Go untuk mendapatkan `snap_token`.
    * Menampilkan pop-up Midtrans Snap di frontend bagi pelanggan untuk membayar via QRIS, e-Wallet, Virtual Account, atau Kartu Kredit.
    * Menangani status pembayaran via HTTP Webhook dari Midtrans.
* **Acceptance Criteria:**
    * GIVEN Kasir memilih metode pembayaran cashless dan mengklik "Bayar", WHEN backend Go berhasil berkomunikasi dengan Midtrans, THEN QRIS/Pop-up Midtrans Snap muncul di layar untuk di-scan pelanggan.
    * GIVEN Pembayaran diselesaikan oleh pelanggan, WHEN Midtrans mengirimkan webhook 'settlement' ke backend Go, THEN status transaksi otomatis berubah menjadi 'Success/Paid' dan struk dapat dicetak.

#### 4.5. Riwayat Transaksi
Fitur untuk melihat kembali transaksi yang telah diproses oleh kasir bersangkutan guna pengecekan ulang.
* **Fungsionalitas:** Menampilkan daftar transaksi yang dibuat khusus selama shift yang sedang aktif pada hari itu.
* **Acceptance Criteria:**
    * GIVEN Kasir membuka modul riwayat, WHEN shift masih berjalan, THEN kasir dapat melihat daftar transaksi lengkap dengan status pembayaran masing-masing, namun tidak diizinkan mengubah item transaksi yang sudah sukses.

---

### 5. Business Rules
Sistem wajib menegakkan aturan-aturan bisnis berikut di tingkat aplikasi (Go) dan database (MySQL):
1.  **Stok hanya berkurang setelah webhook Midtrans confirmed diterima:** Stok produk di gudang/katalog hanya akan berkurang secara permanen setelah sistem menerima webhook resmi dengan status `settlement` atau `capture` dari Midtrans. Selama status masih `pending`, stok tidak berkurang melainkan hanya ditandai sebagai ter-booking sementara guna menghindari *overselling*.
2.  **Cashier harus buka shift sebelum bisa membuat transaksi baru:** Kasir mutlak harus melakukan proses "Buka Shift" dan menginput saldo modal kas awal terlebih dahulu sebelum diizinkan oleh sistem untuk mengakses modul pembuatan transaksi baru atau memproses checkout.
3.  **Satu transaksi hanya bisa menggunakan satu promo:** Di dalam satu transaksi penjualan, sistem hanya memperbolehkan penerapan maksimal 1 (satu) kode promo atau diskon. Penumpukan promo (*promo stacking*) dilarang oleh sistem.
4.  **Order yang sudah checkout tidak bisa diubah itemnya:** Order/transaksi yang telah melalui proses checkout dan telah dinyatakan sukses/lunas oleh sistem (status: `Success/Paid`) bersifat imutabel atau tidak dapat diubah, ditambah, atau dikurangi lagi item pesanan di dalamnya demi menjaga validitas laporan keuangan.

---

### 6. Out of Scope
Fitur-fitur berikut ini secara eksplisit tidak akan dikembangkan pada fase versi ini, namun dapat dipertimbangkan untuk peta jalan (*roadmap*) masa depan:
* **Manajemen Bahan Baku & Resep (Inventory COGS):** Sistem tidak melacak stok biji kopi mentah, susu, atau sirup secara granular (tidak ada fitur pengurangan stok bahan baku berdasarkan resep menu). Sistem hanya melacak stok produk jadi siap jual.
* **Manajemen Multi-Outlet (Multi-Tenancy):** Sistem ini dirancang eksklusif hanya untuk satu outlet coffee shop tunggal (*single outlet*). Tidak mendukung pengelolaan cabang terpusat.
* **Sistem Loyalitas Pelanggan (Customer Loyalty System):** Tidak ada pencatatan data pelanggan, pengumpulan poin reward, tiering member, atau penggunaan voucher loyalitas khusus pelanggan.
* **Manajemen Absensi & Penggajian Karyawan:** Pelacakan jam masuk/keluar kerja staf kasir atau barista di luar manajemen shift kasir, serta kalkulasi gaji bulanan berada di luar cakupan aplikasi POS ini.
* **Integrasi Kurir Pengiriman Pihak Ketiga:** Sistem hanya melayani transaksi langsung di tempat (*dine-in*) dan dibawa pulang (*takeaway*), tidak terintegrasi dengan API logistik online seperti GoFood, GrabFood, atau ShopeeFood.
