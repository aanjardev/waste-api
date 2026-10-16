# Audit Proyek: Waste Collection API

## Konteks

Saya mengerjakan take-home test untuk posisi Backend Developer (Go). Saya minta kamu mengaudit proyek saya **sebagai reviewer teknis**, bukan menulis ulang kode.

## Requirement dari Tes

### Entitas

**Household**
- id (UUID, PK)
- owner_name (string, wajib)
- address (string, wajib)
- created_at, updated_at

**Waste Pickup**
- id (UUID, PK)
- household_id (UUID, FK)
- type: organic | plastic | paper | electronic
- status: pending | scheduled | completed | canceled
- pickup_date (nullable)
- safety_check (bool, wajib untuk electronic)
- created_at, updated_at

**Payment**
- id (UUID, PK)
- household_id (UUID, FK)
- waste_id (UUID, FK, unik)
- amount (decimal)
- payment_date (nullable)
- status: pending | paid | failed
- proof_file_url (nullable)
- created_at, updated_at

### Endpoint

Household:
- POST /api/households
- GET /api/households (pagination)
- GET /api/households/:id
- DELETE /api/households/:id

Pickup:
- POST /api/pickups
- GET /api/pickups (filter status, household_id)
- PUT /api/pickups/:id/schedule
- PUT /api/pickups/:id/complete
- PUT /api/pickups/:id/cancel

Payment:
- POST /api/payments
- GET /api/payments (filter status, household_id)
- PUT /api/payments/:id/confirm (upload file)

Report:
- GET /api/reports/waste-summary
- GET /api/reports/payment-summary

### Business Rules

1. Household tidak bisa buat pickup baru jika punya payment pending.
2. Pickup hanya bisa dijadwalkan jika status pending.
3. Electronic tidak bisa dijadwalkan kecuali safety_check = true.
4. Pickup completed otomatis buat payment: organic/plastic/paper = 50.000, electronic = 100.000.
5. Confirm payment wajib upload file bukti, simpan lokal, path disimpan di DB.

### Persyaratan Teknis

- Docker + docker-compose, satu perintah jalan
- Response konsisten, HTTP status code sesuai
- Validasi input
- Database migration
- Bonus: graceful shutdown, unit test business rules, filter tanggal di GET /api/payments

## Struktur Proyek

(Paste hasil `tree` atau struktur folder di sini)

## Yang Saya Minta Kamu Audit

Tolong audit dengan fokus:

1. **Kebenaran business rules** — apakah 5 rule di atas benar-benar terimplementasi? Ada edge case yang terlewat?
2. **Pemisahan layer** — apakah handler, service, repository sudah terpisah bersih? Ada logic yang bocor ke layer salah?
3. **API design** — apakah RESTful? Status code sesuai? Response konsisten?
4. **Error handling** — apakah error dari DB diterjemahkan ke HTTP status yang tepat? Ada error yang bocor ke user?
5. **Validasi input** — apakah cukup? Ada field yang tidak divalidasi?
6. **Keamanan dasar** — ada SQL injection? Path traversal di upload? File type validation?
7. **Konkurensi & transaksi** — ada race condition? Operasi yang butuh transaksi tapi tidak pakai?
8. **Kualitas kode Go** — idiomatik? Ada yang tidak perlu? Ada yang bisa disederhanakan?
9. **Test coverage** — apakah business rules sudah tercakup? Ada skenario penting yang belum dites?

## Format Output yang Saya Mau

Untuk setiap temuan, beri:

- **Severity**: Critical / High / Medium / Low
- **Lokasi**: file dan fungsi
- **Masalah**: apa yang salah
- **Kenapa penting**: dampaknya
- **Saran perbaikan**: konsep, bukan kode jadi (saya mau tulis sendiri)

Urutkan dari severity tertinggi. Di akhir, beri:

- **Ringkasan**: 3 hal terbaik dari proyek ini
- **3 hal yang paling perlu diperbaiki** sebelum submit
- **Hal yang sudah bagus dan jangan diubah**

## Aturan untuk Kamu (AI)

- Jangan tulis ulang seluruh file. Beri saran konsep.
- Jangan berasumsi ada kode yang tidak saya tunjukkan. Kalau perlu file lain, tanya.
- Kalau ada yang tidak jelas, tanya dulu sebelum menilai.
- Fokus ke requirement tes, bukan preferensi pribadi.
- Jujur kalau ada yang kurang. Saya lebih suka tahu sekarang daripada ditolak recruiter.