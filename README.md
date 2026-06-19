<p align="right">
  <strong>English</strong> | <a href="README.zh.md">中文</a>
</p>

# Bitlog

**Less is More.**

Bitlog is a minimal blogging system inspired by the BearBlog style. Built for those who love writing, it strips away all unnecessary distractions to provide a pure, quiet creative space.

This project has been completely rebuilt with a **React SPA (frontend) + Go (backend) separation architecture**, retaining the original minimal interface and zero-configuration philosophy while significantly improving runtime performance, responsiveness, and cross-platform compilation capabilities.

---

## 🏗️ Architecture & Tech Stack

The refactored system employs a modern, high-performance tech stack:

- **Frontend (frontend/)**: A modern single-page application (SPA) built with **React 19** + **Vite 8** + **TypeScript 6** + **TailwindCSS v4**. Uses React Query for data caching and React Router 7 for dynamic routing, delivering seamless, instant page transitions.
- **Backend (backend/)**: A high-performance RESTful API service built with **Go 1.25** + **Gin framework** + **GORM ORM**.
- **Database**: Uses **SQLite embedded database**. To support `CGO_ENABLED=0` pure Go cross-compilation, the underlying driver uses the pure Go implementation `github.com/glebarez/sqlite`, completely eliminating any dependency on a host GCC compiler.
- **Same-Origin Minimal Deployment**: In production, the frontend build output is bundled into the Go backend's static asset directory. A single Go binary executable handles API services, static file hosting, and file uploads — everything in one deployable unit.

---

## ✨ Core Features

- **Minimalism First**: No ads, no tracking scripts, extremely lightweight — pages load in milliseconds.
- **Seamless Database Compatibility**: The data schema is **100% compatible** with the original Next.js Prisma SQLite database. You can directly swap in the physical database file for zero-migration migration.
- **Powerful Editor**:
  - Full Markdown syntax support.
  - **Local image upload**: Paste (Ctrl+V) or drag-and-drop images directly into the editor.
  - **Smart compression**: Uploaded images are automatically converted to JPEG, proportionally resized (max width 1200px), downsampled with Catmull-Rom smoothing, and compressed at 80% quality. GIF files are preserved as-is, balancing quality and loading speed.
  - **Security validation**: Built-in magic bytes security check and extension sanitization to prevent arbitrary file uploads and malicious script execution.
  - **External video optimization**: Full support for Bilibili and YouTube embeds, automatically adapted to 16:9 aspect ratio with autoplay disabled by default.
- **Flexible Slug Management**: Titles and URLs (slugs) sync automatically, with support for manual lock-in of custom modifications.
- **Security & Rate Limiting**: Login endpoint features built-in IP rate-limiting middleware (max 5 attempts per IP within 5 minutes; a successful login resets the counter).
- **Data Ownership**:
  - **Import/Export**: One-click export of all posts as a zip archive of standard `.md` files with metadata. One-click batch import of `.zip` archives or individual `.md` files.
- **Stealth Mode**: Site visibility can be toggled to **"Private Mode"** with a single click. In this mode, unauthenticated visitors receive fake **404 Not Found** responses for every route, providing perfect cover for a personal notebook.

---

## 📸 Screenshots

| 🏠 Home | 📝 New Post |
| :---: | :---: |
| ![Home](./screenshots/Home.png) | ![New Post](./screenshots/NewPost.png) |

| 📚 Blog List | ⚙️ Settings |
| :---: | :---: |
| ![Blog List](./screenshots/blog.png) | ![Settings](./screenshots/setting1.png) |

---

## 📂 Project Structure

```text
├── frontend/               # React SPA frontend
│   ├── src/                # Components, views, routes, context
│   ├── public/             # Static assets (incl. favicon.ico)
│   └── vite.config.ts      # Dev proxy and Vite configuration
├── backend/                # Go backend
│   ├── db/                 # DB initialization, SQLite singleton, auto-migration & seed
│   ├── handlers/           # Controllers (post CRUD, import/export, image upload, config)
│   ├── middleware/         # Middleware (JWT auth, IP rate limiter, stealth mode, CORS)
│   ├── models/             # Prisma-compatible GORM data models
│   ├── utils/              # Utilities (image resize/compress, JWT signing, zip parsing)
│   └── main.go             # Go API entry point & SPA fallback routing
├── Dockerfile              # Multi-stage Docker build
└── docker-compose.yml      # Deployment with persistent volume mounts
```

---

## 🚀 Production Deployment

### Docker Compose (Recommended)

The simplest, fastest, and cleanest deployment method. The multi-stage build compiles the frontend and statically links the Go backend, then packages everything into a minimal Alpine container.

```bash
git clone https://github.com/sunyan7902/Bitlog.git
cd Bitlog
sudo docker compose up -d --build
```

Visit `http://localhost:3456` to start using.

**Persistence & Volume Mounts:**
- Data is stored in the host's local `./data` directory.
- Database path: `./data/bitlog.db`
- Uploaded images path: `./data/uploads/`

---

## 💻 Local Development

Local development requires [Node.js](https://nodejs.org/) >= 20 and [Go](https://go.dev/) >= 1.25.

### 1. Start the Go Backend
```bash
cd backend
# Automatically downloads dependencies and starts the server (listens on port 3000)
go run main.go
```
On startup, the backend checks whether `backend/data/bitlog.db` exists. If it doesn't, it automatically creates the tables and seeds the default admin account.

### 2. Start the React Frontend
```bash
cd frontend
npm install
npm run dev
```
The frontend runs at `http://localhost:5173`. Vite is pre-configured with a dev proxy that seamlessly forwards `/api` and `/uploads` requests to `http://localhost:3000`.

---

## 🔐 Admin Panel

The public-facing pages do not display an obvious login link, preserving the minimal and clean aesthetic.

- **Login URL**: `/login`
- **Default Username**: `admin`
- **Default Password**: `123456`

> [!IMPORTANT]
> After deployment, **immediately** change the default username and password in **Settings → Account Settings** to ensure security.
>
> The default mode is **public**. If you intend to use it as a personal private notebook, switch the "Site Visibility" to **"Private Mode"** under **Settings**. The system will then return **404 Not Found** for all public-facing requests, achieving perfect stealth.

---

## 📜 License

This project is open-sourced under the MIT License.

**Powered by [Bitlog](https://github.com/sunyan7902/Bitlog)**
