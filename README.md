# STK Menu Tree System

A fullstack hierarchical menu management system with CRUD, drag-and-drop, and nested tree visualization.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.26, Gin, GORM, PostgreSQL |
| **Frontend** | Next.js 16 (App Router), React 19, TypeScript |
| **State** | Zustand |
| **Styling** | Tailwind CSS v4, Lucide Icons |
| **API Docs** | Swagger/OpenAPI |

## Project Structure

```
stk-project/
├── backend/
│   ├── cmd/server/main.go        # Entry point
│   ├── internal/
│   │   ├── config/               # Environment config
│   │   ├── model/                # GORM model
│   │   ├── dto/                  # Request/Response DTOs
│   │   ├── repository/           # Database layer
│   │   ├── service/              # Business logic + tree building
│   │   ├── handler/              # HTTP handlers
│   │   ├── middleware/           # CORS + error recovery
│   │   └── router/               # Route definitions
│   ├── migrations/               # SQL migration reference
│   ├── docs/                     # Swagger generated docs
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── app/                  # Pages (App Router)
│   │   ├── components/
│   │   │   ├── layout/           # Sidebar, Header
│   │   │   └── menus/            # Tree, nodes, modals, detail panel
│   │   ├── lib/                  # Types, API client
│   │   └── store/                # Zustand store
│   └── Dockerfile
├── docker-compose.yml            # Production
├── docker-compose.dev.yml        # Development
└── .env.example
```

## Demo

<video src="https://github.com/yachzu/menu-tree-system/raw/main/docs/demo.mp4" controls width="100%">
  <a href="https://github.com/yachzu/menu-tree-system/raw/main/docs/demo.mp4">Download demo video</a>
</video>

## Quick Start

### Prerequisites

- Go 1.26+
- Node.js 22+
- PostgreSQL 16+ (atau akun [Neon](https://neon.tech) untuk cloud PostgreSQL)
- Docker (optional, for containerized setup)

> **Catatan untuk Reviewer:** File `.env.example` sudah berisi kredensial
> database Neon yang aktif. Cukup rename ke `.env` lalu jalankan aplikasi —
> tidak perlu setup PostgreSQL lokal.

### Development Mode (Without Docker)

**1. Backend**

```bash
cd backend
cp .env.example .env
# Edit .env with your PostgreSQL credentials
go mod tidy
go run ./cmd/server
# Server starts on http://localhost:8080
```

**Menggunakan Neon (Cloud PostgreSQL):**

1. Buat akun di https://neon.tech
2. Buat project baru, dapatkan `DATABASE_URL` dari dashboard
3. Set `DATABASE_URL` di `backend/.env`:
   ```env
   DATABASE_URL=postgresql://user:password@ep-xxx.aws.neon.tech/neondb?sslmode=require
   ```

> **Catatan:** Dengan Neon, Anda tidak perlu menjalankan PostgreSQL lokal. SSL wajib diaktifkan (`sslmode=require`).

**2. Frontend**

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
# App starts on http://localhost:3000
```

### Development Mode (With Docker)

Koneksi database menggunakan **Neon PostgreSQL cloud** — tidak perlu container Postgres lokal.

```bash
docker compose -f docker-compose.dev.yml up --build
# Backend: http://localhost:8080
# Frontend: http://localhost:3000
# Swagger: http://localhost:8080/swagger/index.html
```

### Production Mode (With Docker)

```bash
docker compose up --build
```

> **Catatan:** Service Postgres di-comment karena database menggunakan Neon cloud.
> Jika ingin PostgreSQL lokal, uncomment service `postgres` di file compose dan comment atau hapus `DATABASE_URL`.
>
> **⚠️ Security:** Ganti `DATABASE_URL` di file compose dengan env var `${DATABASE_URL}` dan set nilainya di file `.env` root untuk menghindari hardcoded credentials.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/menus` | Get full menu tree (nested) |
| `GET` | `/api/menus/:id` | Get single menu item |
| `POST` | `/api/menus` | Create new menu item |
| `PUT` | `/api/menus/:id` | Update menu item name |
| `DELETE` | `/api/menus/:id` | Delete menu + cascade children |
| `PATCH` | `/api/menus/:id/move` | Move item to different parent |
| `PATCH` | `/api/menus/:id/reorder` | Reorder within siblings |

API docs available at `/swagger/index.html` when the server is running.

## API Examples

```bash
# Get tree
curl http://localhost:8080/api/menus

# Create root menu
curl -X POST http://localhost:8080/api/menus \
  -H "Content-Type: application/json" \
  -d '{"name": "Dashboard", "parent_id": null}'

# Create child menu
curl -X POST http://localhost:8080/api/menus \
  -H "Content-Type: application/json" \
  -d '{"name": "Analytics", "parent_id": "<parent-uuid>"}'

# Update name
curl -X PUT http://localhost:8080/api/menus/<uuid> \
  -H "Content-Type: application/json" \
  -d '{"name": "New Name"}'

# Delete (cascades to children)
curl -X DELETE http://localhost:8080/api/menus/<uuid>
```

## Database Schema

```sql
CREATE TABLE menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    parent_id UUID REFERENCES menus(id) ON DELETE CASCADE,
    depth INTEGER NOT NULL DEFAULT 0,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_menus_parent_id ON menus(parent_id);
CREATE INDEX idx_parent_order ON menus(parent_id, order_index);
```

## Frontend Features

- **Tree View**: Recursive nested menu display with visual hierarchy lines
- **CRUD**: Add, edit, delete menus via modals
- **Detail Panel**: Right-side panel for editing selected menu
- **Drag & Drop**: Move menus between parents via HTML5 drag-and-drop
- **Search**: Real-time search/filter with Ctrl+F shortcut
- **Expand/Collapse**: Individual or bulk expand/collapse
- **Responsive**: Mobile sidebar, adaptive layout
- **States**: Loading spinners, error messages, empty state placeholders

## Running Tests

```bash
# Backend unit tests (26 test cases)
cd backend
go test ./internal/service/ -v

# All backend tests
go test ./... -v
```

## Architecture Decisions

- **Adjacency List** pattern for hierarchy (parent_id self-reference)
- **In-memory tree building** (fetch all → build map → attach children) avoids recursive DB CTEs
- **Denormalized depth** column for O(1) read performance
- **GORM AutoMigrate** for schema management
- **Interface-based repository** pattern for testability
- **Zustand** for lightweight state management without boilerplate
- **Custom API client** with fetch (no Axios dependency)
- **HTML5 native drag-and-drop** (no external library dependency)

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Backend server port |
| `DATABASE_URL` | `postgresql://...` | Full PostgreSQL connection string (Neon atau PostgreSQL apapun) |
| `DATABASE_HOST` | `localhost` | DB host (if not using DATABASE_URL) |
| `DATABASE_PORT` | `5432` | DB port |
| `DATABASE_USER` | `postgres` | DB user |
| `DATABASE_PASSWORD` | `postgres` | DB password |
| `DATABASE_NAME` | `menu_tree` | DB name |
| `DATABASE_SSLMODE` | `disable` | SSL mode |
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080/api` | API base URL (frontend) |

> **Neon Users:** Dapatkan `DATABASE_URL` dari dashboard Neon project Anda.
> Pastikan sertakan `?sslmode=require` karena Neon mewajibkan SSL untuk koneksi.
