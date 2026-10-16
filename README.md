# Waste Collection API

Backend API untuk Sistem Pengumpulan Sampah Komunitas. Dibangun dengan Go, PostgreSQL, dan Docker.

## Tech Stack

- Language: Go 1.27
- Database: PostgreSQL 16
- HTTP Framework: net/http (Go 1.27 routing)
- DB Driver: pgx/v5
- Containerization: Docker + Docker Compose
- Validation: go-playground/validator
- Testing: testify

## Fitur

- CRUD Household dengan pagination
- Waste Pickup dengan filter status & household
- Payment dengan upload bukti bayar
- Reporting: waste summary & payment summary
- Business rules:
  1. Household tidak bisa buat pickup baru jika punya payment pending
  2. Pickup hanya bisa dijadwalkan jika status pending
  3. Electronic tidak bisa dijadwalkan kecuali safety_check = true
  4. Pickup completed otomatis membuat payment (50000 untuk organic/plastic/paper, 100000 untuk electronic)
  5. Confirm payment wajib upload file bukti

## Struktur Proyek

    waste-api/
    ├── cmd/api/                    # Entry point
    ├── internal/
    │   ├── config/                 # Config & DB connection
    │   ├── handler/                # HTTP handler
    │   ├── service/                # Business logic
    │   │   └── mocks/              # Mock untuk unit test
    │   ├── repository/             # Database query
    │   ├── model/                  # Entity & error
    │   ├── dto/                    # Request/Response DTO
    │   ├── middleware/             # HTTP middleware
    │   └── response/               # Response helper
    ├── migrations/                 # SQL migration & seed
    ├── uploads/                    # File upload (Docker volume)
    ├── docker-compose.yml
    ├── Dockerfile
    ├── Makefile
    └── README.md

Alur layer: handler -> service -> repository. Handler hanya urus HTTP, service urus business logic, repository urus SQL.

## Setup & Menjalankan

### Prasyarat

- Docker & Docker Compose
- Go 1.22+ (hanya jika mau jalankan tanpa Docker)

### 1. Clone repository

    git clone https://github.com/aanjardev/waste-api.git
    cd waste-api

Ganti `aanjardev` dengan username GitHub kamu.

### 2. Copy environment file

    cp .env.example .env

### 3. Jalankan Docker

    docker compose up --build

Tunggu sampai muncul log:

    waste_app | connected to database
    waste_app | server running on :8080

Migration dan seed otomatis dijalankan saat container pertama kali start.

### 4. Tes

    curl http://localhost:8080/health

Harus `OK`.

## Environment Variables

| Variable | Default | Deskripsi |
|---|---|---|
| APP_PORT | 8080 | Port HTTP server |
| APP_ENV | development | Environment |
| BASE_URL | http://localhost:8080 | Base URL untuk file |
| UPLOAD_DIR | ./uploads | Folder upload |
| DB_HOST | db | Host PostgreSQL |
| DB_PORT | 5432 | Port PostgreSQL |
| DB_USER | postgres | User PostgreSQL |
| DB_PASSWORD | postgres | Password PostgreSQL |
| DB_NAME | waste_db | Nama database |
| DB_SSLMODE | disable | SSL mode |

## Menjalankan Migration & Seeding

Migration dan seed otomatis dijalankan saat container app pertama kali start. Kamu tidak perlu menjalankan perintah manual.

Kalau mau reset database:

    docker compose down -v
    docker compose up --build

Kalau mau jalankan migration manual (untuk development):

    make migrate-up
    make migrate-down

Manual (legacy):

    docker compose exec -T db psql -U postgres -d waste_db < migrations/000001_init.up.sql
    docker compose exec -T db psql -U postgres -d waste_db < migrations/000002_seed.up.sql

## Menjalankan Test

    go test ./... -v

Unit test untuk business rules ada di `internal/service/*_test.go`.

## API Endpoints

### Household

| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | /api/households | Buat household baru |
| GET | /api/households?page=&limit= | List household |
| GET | /api/households/{id} | Detail household |
| DELETE | /api/households/{id} | Hapus household |

### Waste Pickup

| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | /api/pickups | Buat pickup |
| GET | /api/pickups?status=&household_id=&page=&limit= | List pickup |
| PUT | /api/pickups/{id}/schedule | Jadwalkan pickup |
| PUT | /api/pickups/{id}/complete | Tandai selesai |
| PUT | /api/pickups/{id}/cancel | Batalkan |

### Payment

| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | /api/payments | Buat payment |
| GET | /api/payments?status=&household_id=&from=&to=&page=&limit= | List payment |
| PUT | /api/payments/{id}/confirm | Konfirmasi + upload bukti |

### Report

| Method | Endpoint | Deskripsi |
|---|---|---|
| GET | /api/reports/waste-summary | Agregasi pickup per type & status |
| GET | /api/reports/payment-summary | Total payment per status & revenue |

## Contoh Request

### Create Household

    curl -X POST http://localhost:8080/api/households \
      -H "Content-Type: application/json" \
      -d '{"owner_name":"Budi","address":"Jl. Merdeka No. 1"}'

### Create Pickup

    curl -X POST http://localhost:8080/api/pickups \
      -H "Content-Type: application/json" \
      -d '{"household_id":"<uuid>","type":"organic"}'

### Schedule Pickup

    curl -X PUT http://localhost:8080/api/pickups/<id>/schedule \
      -H "Content-Type: application/json" \
      -d '{"pickup_date":"2026-10-20T10:00:00Z"}'

### Confirm Payment

    curl -X PUT http://localhost:8080/api/payments/<id>/confirm \
      -F "proof_file=@bukti.png"

## Response Format

### Sukses

    {
      "success": true,
      "message": "household created",
      "data": { ... }
    }

### Sukses dengan pagination

    {
      "success": true,
      "message": "households retrieved",
      "data": [ ... ],
      "meta": {
        "page": 1,
        "limit": 10,
        "total": 3,
        "total_pages": 1
      }
    }

### Error

    {
      "success": false,
      "message": "validation error",
      "errors": {
        "owner_name": "is required"
      }
    }

## Status Code

- 200 OK
- 201 Created
- 400 Bad Request
- 404 Not Found
- 409 Conflict
- 422 Unprocessable Entity
- 500 Internal Server Error

## Postman Collection

Import file `postman_collection.json` ke Postman. Environment variable `base_url` sudah diset ke `http://localhost:8080`.

## Known Limitations

- **Partial migration failure**: Jika `000001_init.up.sql` sukses tapi `000002_seed.up.sql` gagal, restart container berikutnya akan skip migration karena tabel `households` sudah ada. Untuk production, gunakan migration tool seperti `golang-migrate` yang melacak versi di tabel `schema_migrations`.
- **Upload file**: Validasi MIME menggunakan `http.DetectContentType` (512 byte pertama). Polyglot file tingkat lanjut masih bisa lolos. Untuk production, tambahkan scanning atau validasi lebih ketat.
- **ON DELETE CASCADE pada payments**: Menghapus household akan menghapus semua payment terkait. Untuk audit trail keuangan di production, pertimbangkan ON DELETE RESTRICT.

## Catatan

- File upload disimpan di folder `uploads/` dan di-mount sebagai Docker volume.
- Enum di PostgreSQL digunakan untuk `type`, `status` pickup, dan `status` payment.
- Unit test untuk business rules ada di `internal/service/*_test.go`.