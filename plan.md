# Plan: Upload Foto Court ke Supabase Storage (terintegrasi di Create/Update)

## Konteks

- Bucket Supabase: **public** (`SUPABASE_STORAGE_BUCKET=court`, sudah ada di `.env`).
- Hak akses: **hanya admin** (guard `isAdminFromContext` di handler).
- **Tidak ada endpoint upload terpisah** — backend menerima file langsung di `POST /api/courts` & `PUT /api/courts/:id`, BE yang upload ke Supabase.
- SDK: `github.com/supabase-community/storage-go v0.7.0` (sudah indirect dep, dijadikan direct dep).

## Kontrak API (breaking change)

- `POST /api/courts` & `PUT /api/courts/:id` berubah dari `application/json` → **multipart/form-data**.
- Field: `name`, `type`, `price_per_hour`, `location` + file `image` (opsional).
- `GET /api/courts/:id/image` (JWT): stream file asli dari Supabase Storage (`c.Data`), dipakai untuk mengambil image dari storage; `image_url` tetap dikembalikan pada `GET` list/detail.
- `PUT /api/courts/:id` tetap **partial update** — field/file yang tidak dikirim = nilai lama dipertahankan.
- Batasan file: ekstensi `jpg/jpeg/png/webp`, maksimal 5MB.

## Alur

- **Create**: insert court (tanpa image, dapat `id`) → jika `image` ada: upload ke `courts/<id>/<hex>.<ext>` → simpan `image_url` ke DB.
- **Update**: terapkan field parsial → jika `image` baru: upload dulu → simpan `image_url` → hapus objek lama di bucket (best-effort).
- **Delete court**: hapus objek `image_url` dari bucket (best-effort) bersama penghapusan row.

## Perubahan per file

| File | Perubahan |
|---|---|
| `configs/config.go` | Tambah env `SUPABASE_STORAGE_URL` (bind + `requiredEnv`); fallback derive dari `SUPABASE_URL` strip `/rest/v1/`. |
| `go.mod` / `go.sum` | `storage-go v0.7.0` jadi direct dep; `go mod tidy`. |
| `storage/storage.go` (baru) | `StorageClient` wrap SDK; interface `Storage{ Upload(relativePath, io.Reader, contentType) (publicURL, error); Download(publicURL) ([]byte, contentType, error); Delete(publicURL) error }`; helper `ContentType(path)` & ekstrak path dari URL. |
| `databases/migration/sql_migration/000006_migration_add_court_image.sql` (baru) | `ALTER TABLE courts ADD COLUMN IF NOT EXISTS image_url TEXT;` |
| `model/entity/courts.go` | Tambah `ImageURL *string` (`db:"image_url"`). |
| `model/dto/request/courtRequest.go` | Ganti tag `json` → `form` di `CourtRequest` & `UpdateCourtRequest` (field pointer tetap untuk partial update). |
| `model/dto/response/courtResponse.go` | Tambah `ImageURL *string` (`json:"image_url,omitempty"`). |
| `repository/courtRepository.go` | Tambah `image_url` di `Create`/`FindAll`/`FindByID`/`Update`; method baru `UpdateImageURL(id, *string)`. |
| `service/courtService.go` | `NewCourtService(courtRepo, storage.Storage)`; alur sesuai deskripsi; `GetImage` download & stream; update `toCourtResponse`. |
| `handler/courtHandler.go` | `Create` & `Update` pakai `c.ShouldBind` (form) + `c.FormFile("image")`; validasi ekstensi (`jpg/jpeg/png/webp`) & ukuran ≤ 5MB; guard admin; `GetImage` stream via `c.Data`. |
| `routes/routes.go` | Tambah `GET /api/courts/:id/image`. |
| `main.go` & `api/index.go` | Injeksi `storage.NewStorageClient(...)` ke `NewCourtService`. |
| `.env` / `README.md` | Tambah `SUPABASE_STORAGE_URL`; update dokumentasi kontrak multipart. |

## Validasi

- `go build ./...` dan `go vet ./...`.
- Manual:
  ```bash
  curl -F "name=GOR A" -F "type=futsal" -F "price_per_hour=150000" \
       -F "location=Bandung" -F "image=@foto.jpg" \
       -H "Authorization: Bearer <ADMIN_TOKEN>" \
       -X POST http://localhost:8080/api/courts
  ```

## Catatan / keputusan

- **Ganti gambar**: upload baru dulu, baru hapus objek lama (best-effort) supaya tidak ada celah URL kehilangan.
- **Hapus file saat court di-delete**: ikut dihapus (best-effort).