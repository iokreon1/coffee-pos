# API Contract
## Sistem Point of Sales (POS) Coffee Shop

### 1. Base URL & Versioning
- **Base URL:** `https://api.yourdomain.com` (atau `http://localhost:8080` untuk local development)
- **Versioning:** `/api/v1`

### 2. Authentication
Semua endpoint yang dilindungi (protected) wajib menyertakan token JWT pada header HTTP.

```http
Authorization: Bearer <token_jwt>
```

### 3. HTTP Status Codes
Berikut adalah daftar HTTP status code yang digunakan secara konsisten dalam API ini:
- **200 OK:** Sukses untuk request GET dan PUT.
- **201 Created:** Sukses untuk request POST (pembuatan data baru).
- **400 Bad Request:** Request tidak valid (parameter salah atau logika tidak memenuhi syarat).
- **401 Unauthorized:** Tidak ada token JWT, token salah, atau token expired.
- **403 Forbidden:** Token valid tapi role user tidak punya hak akses ke endpoint tersebut.
- **404 Not Found:** Data yang diminta tidak ditemukan.
- **422 Unprocessable Entity:** Validasi input gagal.
- **429 Too Many Requests:** Terkena rate limit.
- **500 Internal Server Error:** Terjadi kesalahan di sisi server.

### 4. Standard Response Format
Sistem ini menggunakan format respons standar untuk memudahkan parsing di sisi client.

#### 4.1. Success Response
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": { ... }
}
```

#### 4.2. Success Response with List (Pagination)
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [ ... ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_pages": 5,
    "total_data": 45
  }
}
```

#### 4.3. Validation Error Response (422)
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "field_name": "Error description message"
  }
}
```

#### 4.4. General Error Response (4xx, 5xx)
```json
{
  "success": false,
  "message": "Error description message"
}
```

---

### 5. Endpoints

#### 5.1. Auth (Public)

**1. Login User**
- **Method & Path:** `POST /api/v1/auth/login`
- **Request Body:**
  ```json
  {
    "email": "owner@coffeeshop.com",
    "password": "secretpassword"
  }
  ```
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Login successful",
    "data": {
      "token": "eyJhb...",
      "user": { "id": "uuid", "name": "Owner", "role": "OWNER" }
    }
  }
  ```

**2. Refresh Token**
- **Method & Path:** `POST /api/v1/auth/refresh`
- **Request Body:**
  ```json
  {
    "refresh_token": "token_refresh_..."
  }
  ```
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Token refreshed",
    "data": {
      "token": "eyJhb_new..."
    }
  }
  ```

**3. Logout User**
- **Method & Path:** `POST /api/v1/auth/logout`
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Logout successful"
  }
  ```

---

#### 5.2. Categories (Owner)

**1. Get All Categories**
- **Method & Path:** `GET /api/v1/owner/categories`
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Categories retrieved",
    "data": [ { "id": "uuid", "name": "Coffee" } ]
  }
  ```

**2. Create Category**
- **Method & Path:** `POST /api/v1/owner/categories`
- **Request Body:** `{ "name": "Pastry" }`
- **Response Sukses (201 Created):**
  ```json
  {
    "success": true,
    "message": "Category created",
    "data": { "id": "uuid", "name": "Pastry" }
  }
  ```

**3. Update Category**
- **Method & Path:** `PUT /api/v1/owner/categories/:id`
- **Request Body:** `{ "name": "Desserts" }`
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Category updated",
    "data": { "id": "uuid", "name": "Desserts" }
  }
  ```

**4. Delete Category**
- **Method & Path:** `DELETE /api/v1/owner/categories/:id`
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Category deleted successfully"
  }
  ```

---

#### 5.3. Products (Owner)

**1. Get All Products**
- **Method & Path:** `GET /api/v1/owner/products`
- **Query Params:** `?category_id=uuid&search=espresso&page=1&limit=10`
- **Response Sukses (200 OK):** (Menyertakan blok `meta` untuk pagination)

**2. Get Product Details**
- **Method & Path:** `GET /api/v1/owner/products/:id`
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Product retrieved",
    "data": {
      "id": "uuid",
      "sku": "CFF-ESP-01",
      "name": "Espresso",
      "price": 2000000
    }
  }
  ```

**3. Create Product**
- **Method & Path:** `POST /api/v1/owner/products`
- **Request Body:** JSON / Form-Data `{ "name": "Latte", "price": 2500000, "category_id": "uuid" }`
- **Response Sukses (201 Created):**
  ```json
  {
    "success": true,
    "message": "Product created",
    "data": { "id": "uuid", "name": "Latte" }
  }
  ```

**4. Update Product**
- **Method & Path:** `PUT /api/v1/owner/products/:id`
- **Request Body:** Sama seperti Create.
- **Response Sukses (200 OK):** Pesan update sukses beserta data terbaru.

**5. Delete Product (Soft Delete)**
- **Method & Path:** `DELETE /api/v1/owner/products/:id`
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Product deleted"
  }
  ```

---

#### 5.4. Stock (Owner)

**1. Get Current Stock**
- **Method & Path:** `GET /api/v1/owner/products/:id/stock`
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Stock retrieved",
    "data": { "product_id": "uuid", "current_stock": 50 }
  }
  ```

**2. Adjust Stock**
- **Method & Path:** `POST /api/v1/owner/products/:id/stock/adjustment`
- **Request Body:**
  ```json
  {
    "type": "adjustment",
    "quantity": 10,
    "notes": "Restock dari supplier"
  }
  ```
- **Response Sukses (201 Created):**
  ```json
  {
    "success": true,
    "message": "Stock adjusted",
    "data": { "current_stock": 60 }
  }
  ```

**3. Get Stock Movements**
- **Method & Path:** `GET /api/v1/owner/products/:id/stock/movements`
- **Response Sukses (200 OK):** Mengembalikan array riwayat perubahan stok untuk produk tersebut.

---

#### 5.5. Tables (Owner)

**1. Get All Tables**
- **Method & Path:** `GET /api/v1/owner/tables`
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Tables retrieved",
    "data": [ { "id": "uuid", "table_number": "Meja 01", "status": "AVAILABLE" } ]
  }
  ```

**2. Create Table**
- **Method & Path:** `POST /api/v1/owner/tables`
- **Request Body:** `{ "table_number": "Meja 01" }`
- **Response Sukses (201 Created):** Data meja.

**3. Update Table**
- **Method & Path:** `PUT /api/v1/owner/tables/:id`
- **Request Body:** `{ "table_number": "Meja VIP" }`
- **Response Sukses (200 OK):** Data meja terupdate.

**4. Delete Table**
- **Method & Path:** `DELETE /api/v1/owner/tables/:id`
- **Response Sukses (200 OK):** Meja terhapus.

---

#### 5.6. Cashier Management (Owner)

**1. Get All Cashiers**
- **Method & Path:** `GET /api/v1/owner/cashiers`
- **Response Sukses (200 OK):** Array user dengan role CASHIER.

**2. Create Cashier**
- **Method & Path:** `POST /api/v1/owner/cashiers`
- **Request Body:** `{ "name": "Kasir", "email": "kasir@co.com", "password": "pass" }`
- **Response Sukses (201 Created):** Data kasir baru.

**3. Update Cashier**
- **Method & Path:** `PUT /api/v1/owner/cashiers/:id`
- **Request Body:** `{ "name": "Kasir Baru" }`
- **Response Sukses (200 OK):** Kasir terupdate.

**4. Toggle Cashier Status**
- **Method & Path:** `PATCH /api/v1/owner/cashiers/:id/toggle-status`
- **Request Body:** `{ "is_active": false }`
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Cashier status updated",
    "data": { "id": "uuid", "is_active": false }
  }
  ```

---

#### 5.7. Promos (Owner)

**1. Get All Promos**
- **Method & Path:** `GET /api/v1/owner/promos`
- **Response Sukses (200 OK):** List semua promo.

**2. Create Promo**
- **Method & Path:** `POST /api/v1/owner/promos`
- **Request Body:**
  ```json
  {
    "code": "JUMATBERKAH",
    "name": "Diskon Jumat Berkah 10%",
    "promo_type": "PERCENTAGE",
    "value": 10,
    "min_purchase": 5000000,
    "start_date": "2026-07-01T00:00:00Z",
    "end_date": "2026-07-31T23:59:59Z"
  }
  ```
- **Response Sukses (201 Created):** Data promo.

**3. Update Promo**
- **Method & Path:** `PUT /api/v1/owner/promos/:id`
- **Request Body:** Sama seperti Create.
- **Response Sukses (200 OK):** Promo terupdate.

**4. Delete Promo**
- **Method & Path:** `DELETE /api/v1/owner/promos/:id`
- **Response Sukses (200 OK):** Promo terhapus.

---

#### 5.8. Reports (Owner)

**1. Get Revenue Report**
- **Method & Path:** `GET /api/v1/owner/reports/revenue`
- **Query Params:** `?start_date=2026-06-01&end_date=2026-06-30`
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Revenue retrieved",
    "data": { "total_revenue": 5000000000, "total_transactions": 250 }
  }
  ```

**2. Get Top Products**
- **Method & Path:** `GET /api/v1/owner/reports/top-products`
- **Response Sukses (200 OK):** Array produk terlaris.

**3. Get Cashier Summary**
- **Method & Path:** `GET /api/v1/owner/reports/cashier-summary`
- **Response Sukses (200 OK):** Performa per kasir.

**4. Export Report**
- **Method & Path:** `GET /api/v1/owner/reports/export`
- **Response:** Unduhan file CSV.

---

#### 5.9. Shifts (Cashier)

**1. Open Shift**
- **Method & Path:** `POST /api/v1/cashier/shifts/open`
- **Request Body:** `{ "opening_cash": 50000000 }`
- **Response Sukses (201 Created):**
  ```json
  {
    "success": true,
    "message": "Shift opened",
    "data": { "id": "uuid", "status": "open", "opening_cash": 50000000 }
  }
  ```

**2. Close Shift**
- **Method & Path:** `POST /api/v1/cashier/shifts/close`
- **Request Body:** `{ "closing_cash": 150000000 }`
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Shift closed"
  }
  ```

**3. Get Current Active Shift**
- **Method & Path:** `GET /api/v1/cashier/shifts/current`
- **Response Sukses (200 OK):** Menampilkan shift kasir yang sedang `open`.

---

#### 5.10. Orders (Cashier)

**1. Create Initial Order (Draft)**
- **Method & Path:** `POST /api/v1/cashier/orders`
- **Request Body:** `{ "table_id": "uuid" }` // Opsional jika dine-in
- **Response Sukses (201 Created):**
  ```json
  {
    "success": true,
    "message": "Order created",
    "data": { "id": "order-uuid", "status": "draft" }
  }
  ```

**2. Get Order Details**
- **Method & Path:** `GET /api/v1/cashier/orders/:id`
- **Response Sukses (200 OK):** Detail pesanan beserta item di dalamnya.

**3. Add Items to Order**
- **Method & Path:** `POST /api/v1/cashier/orders/:id/items`
- **Request Body:**
  ```json
  {
    "product_id": "uuid",
    "quantity": 2,
    "notes": "Less ice"
  }
  ```
- **Response Sukses (201 Created):**
  ```json
  {
    "success": true,
    "message": "Item added",
    "data": { "item_id": "uuid", "subtotal": 4000000 }
  }
  ```

**4. Update Item in Order**
- **Method & Path:** `PUT /api/v1/cashier/orders/:id/items/:item_id`
- **Request Body:** `{ "quantity": 3 }`
- **Response Sukses (200 OK):** Kuantitas dan subtotal terupdate.

**5. Delete Item from Order**
- **Method & Path:** `DELETE /api/v1/cashier/orders/:id/items/:item_id`
- **Response Sukses (200 OK):** Item dihapus dari pesanan.

**6. Apply Promo**
- **Method & Path:** `POST /api/v1/cashier/orders/:id/promo`
- **Request Body:** `{ "promo_code": "JUMATBERKAH" }`
- **Response Sukses (200 OK):** Diskon diterapkan dan total akhir terupdate.

**7. Remove Promo**
- **Method & Path:** `DELETE /api/v1/cashier/orders/:id/promo`
- **Response Sukses (200 OK):** Promo dibatalkan dari transaksi.

**8. Get Shift Order History**
- **Method & Path:** `GET /api/v1/cashier/orders/history`
- **Response Sukses (200 OK):** Daftar pesanan di shift yang sedang aktif.

---

#### 5.11. Payments

**1. Process Checkout (Midtrans Snap)**
- **Method & Path:** `POST /api/v1/cashier/orders/:id/checkout`
- **Akses:** Cashier
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Checkout initiated",
    "data": {
      "snap_token": "midtrans-snap-token-xyz123",
      "redirect_url": "https://app.midtrans.com/snap/v2/vtweb/xyz123"
    }
  }
  ```

**2. Midtrans Webhook (Callback)**
- **Method & Path:** `POST /api/v1/webhooks/midtrans`
- **Akses:** Public (Midtrans Server)
- **Request Body:** Format standar notifikasi HTTP Midtrans.
- **Response Sukses:** `200 OK` tanpa body khusus.

**3. Get Payment Status**
- **Method & Path:** `GET /api/v1/cashier/orders/:id/payment-status`
- **Akses:** Cashier
- **Response Sukses (200 OK):**
  ```json
  {
    "success": true,
    "message": "Payment status retrieved",
    "data": { "order_id": "uuid", "status": "paid" }
  }
  ```
