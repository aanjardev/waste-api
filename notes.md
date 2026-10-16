# Catatan Proyek: Waste Collection API

## Tujuan
API untuk sistem pengumpulan sampah komunitas.
Stack: Go + PostgreSQL + Docker.

## Entitas

### Household
- id (UUID, PK)
- owner_name (string, wajib)
- address (string, wajib)
- created_at, updated_at

### Waste Pickup
- id (UUID, PK)
- household_id (UUID, FK)
- type: organic | plastic | paper | electronic
- status: pending | scheduled | completed | canceled
- pickup_date (nullable)
- safety_check (bool, wajib untuk electronic)
- created_at, updated_at

### Payment
- id (UUID, PK)
- household_id (UUID, FK)
- waste_id (UUID, FK, unik)
- amount (decimal)
- payment_date (nullable)
- status: pending | paid | failed
- proof_file_url (nullable, diisi saat confirm)
- created_at, updated_at

## Endpoint

### Household
- POST   /api/households
- GET    /api/households?page=&limit=
- GET    /api/households/:id
- DELETE /api/households/:id

### Pickup
- POST /api/pickups
- GET  /api/pickups?status=&household_id=
- PUT  /api/pickups/:id/schedule
- PUT  /api/pickups/:id/complete
- PUT  /api/pickups/:id/cancel

### Payment
- POST /api/payments
- GET  /api/payments?status=&household_id=
- PUT  /api/payments/:id/confirm  (multipart, upload file)

### Report
- GET /api/reports/waste-summary
- GET /api/reports/payment-summary

## Business Rules
1. Household tidak bisa buat pickup baru jika punya payment status pending.
2. Pickup hanya bisa dijadwalkan jika status = pending.
3. Electronic tidak bisa dijadwalkan kecuali safety_check = true.
4. Pickup completed -> auto create payment:
   - organic/plastic/paper = 50.000
   - electronic = 100.000
5. Confirm payment wajib upload file bukti. Disimpan lokal, path disimpan di DB.

## Response Format
Sukses:
{
  "success": true,
  "message": "...",
  "data": {...},
  "meta": { "page": 1, "limit": 10, "total": 0, "total_pages": 0 }
}

Error:
{
  "success": false,
  "message": "...",
  "errors": {...}
}

## Status Code
- 200 OK
- 201 Created
- 400 Bad Request
- 404 Not Found
- 409 Conflict
- 422 Unprocessable Entity
- 500 Internal Server Error

## Rencana 5 Hari
- Hari 1: setup Docker, migration, config, model, Household API
- Hari 2: Pickup API + rule 1-3
- Hari 3: Payment API + upload + rule 4-5
- Hari 4: Report, filter, pagination, validasi
- Hari 5: Test, README, Postman, publish

## Catatan Pribadi
(tulis progress harian di sini)