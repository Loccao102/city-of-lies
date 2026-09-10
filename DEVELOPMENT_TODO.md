# CITY OF LIES — Kế Hoạch Phát Triển & Danh Sách Hạng Mục Chưa Làm (Development Backlog)

> **Tài liệu tham chiếu gốc:** [CITY_OF_LIES_MASTER_BIBLE_GO.md](file:///c:/Users/Admin/city-of-lies/CITY_OF_LIES_MASTER_BIBLE_GO.md)  
> **Trạng thái hiện tại:** Gameplay V1 hoàn chỉnh (Core Simulation + Next.js 15 / Three.js 3D Web App đã hoạt động và build thành công).

---

## I. TỔNG QUAN TIẾN ĐỘ THEO CÁC PHASE

| Phase | Nội dung | Trạng thái | Ghi chú |
|---|---|---|---|
| **Phase 0** | Cấu trúc repo, tài liệu, ADR, Scenario JSON | ✅ **Hoàn thành** | Đầy đủ 9 dataset Riverside |
| **Phase 1** | Go Skeleton, Config, Health, Graceful Shutdown | ✅ **Hoàn thành** | Go 1.24+ / Chi HTTP / Slog |
| **Phase 2** | Pure Domain Math (Propagation, Narrative, Evidence) | ✅ **Hoàn thành** | Khớp công thức Master Bible |
| **Phase 3** | Scenario Loader & Validator CLI | ✅ **Hoàn thành** | `bin/scenario-check.exe` PASS |
| **Phase 4** | PostgreSQL Persistence (Migrations & sqlc) | 🟡 **Cần kích hoạt DB** | Đã có migrations & queries SQL; hiện chạy In-Memory Engine |
| **Phase 5** | Simulation Engine & Headless Calibration | ✅ **Hoàn thành** | `bin/sim.exe` đạt 75% tin giả lúc 18:34 |
| **Phase 6** | HTTP API + WebSocket Hub | ✅ **Hoàn thành** | Chi router + Coder WebSocket |
| **Phase 7** | Frontend Shell & State Management | ✅ **Hoàn thành** | Next.js 15, React 19, Zustand |
| **Phase 8** | Investigation Mechanics (Interview, Notebook, Truth) | ✅ **Hoàn thành** | Đầy đủ modal & tương tác |
| **Phase 9** | 3D World (Isometric Three.js Canvas) | ✅ **Hoàn thành** | 12 địa điểm, 20 agents, rumor pulse |
| **Phase 10** | External LLM Integration (OpenAI / Claude / Gemini) | 🟡 **Chờ API Key** | Template Provider đã chạy tốt; chờ cắm Key thật |
| **Phase 11** | Cân bằng tự động (Soak Test), 3D Model & Âm thanh | ⏳ **Chưa làm** | Chi tiết bên dưới |
| **Phase 12** | CI/CD GitHub Actions & Demo Deployment | ⏳ **Chưa làm** | Chi tiết bên dưới |

---

## II. CHI TIẾT CÁC HẠNG MỤC CHƯA LÀM (TODO BACKLOG)

### 1. Phase 4 — Kích hoạt PostgreSQL Persistence (Tùy chọn)
Hiện tại game đang chạy trên **In-Memory Simulation Engine** siêu tốc (không cần cài PostgreSQL vẫn chơi mượt mà 100%). Nếu muốn lưu session lâu dài vào PostgreSQL:
- [ ] Khởi chạy PostgreSQL: `docker compose up -d postgres`.
- [ ] Chạy migration bảng dữ liệu: `backend/db/migrations/000001_init.up.sql`.
- [ ] Chạy `sqlc generate` trong thư mục `backend/` để sinh Go repository code tương tác DB.
- [ ] Cấu hình biến môi trường `CITYOFLIES_DATABASE_URL=postgres://cityoflies:cityoflies_dev_only@localhost:5432/cityoflies?sslmode=disable`.
- [ ] Ghi nhận snapshot và lịch sử `world_events` vào PostgreSQL sau mỗi lượt chơi.

---

### 2. Phase 10 — Tích hợp LLM Động (OpenAI / Gemini / Claude)
Game hiện đang sử dụng **Deterministic Template Dialogue Provider** (miễn phí, không tốn token, phản hồi tức thì và an toàn tuyệt đối theo luật Locked Decisions). Khi muốn bật AI đối thoại tự nhiên:
- [ ] Điền API Key vào file `.env`:
  ```bash
  CITYOFLIES_LLM_ENABLED=true
  CITYOFLIES_LLM_PROVIDER=openai
  CITYOFLIES_LLM_BASE_URL=https://api.openai.com/v1
  CITYOFLIES_LLM_API_KEY=sk-...
  CITYOFLIES_LLM_MODEL=gpt-4o-mini
  ```
- [ ] Kiểm thử `KnowledgeGuard` để đảm bảo model không rò rỉ Ground Truth bí mật của Server hoặc bịa đặt số liệu người chết sai lệch với trí nhớ của NPC.
- [ ] Kiểm tra cơ chế Fallback: Nếu API bị timeout hoặc rate limit, hệ thống tự động rơi về Template Provider mà không làm đứt đoạn màn chơi.

---

### 3. Phase 11 — Tinh chỉnh Cân bằng, 3D Assets & Âm thanh (Audio/Visual Polish)
- [ ] **Soak Test 100 Seeds tự động**:
  - Viết script chạy `bin/sim.exe` qua 100 seeds khác nhau (từ seed 1 đến 100).
  - Thu thập thống kê: tỉ lệ thua cuộc tự nhiên không có người chơi (mục tiêu > 90%), thời gian trung vị đạt 75% tin giả (mục tiêu 12–18 phút).
- [ ] **Mô hình 3D Nhân vật Bespoke (GLTF/GLB)**:
  - Thay thế các hình khối Cylinder/Sphere cơ bản của 20 nhân vật bằng 4–6 bộ model low-poly người (Bác sĩ/Y tá, Cứu hỏa, Công nhân, Thường dân, Bảo vệ, Phóng viên).
  - Thêm animation: Idle, Walk giữa các địa điểm, Talk khi có tin đồn phát xung.
- [ ] **Thiết kế Âm thanh (Soundscape)**:
  - Thêm tiếng còi hú xe cứu hỏa lúc 18:03 và xe cấp cứu lúc 18:04.
  - Tiếng chuông thông báo khi có tin đồn mới trên mạng xã hội.
  - Tiếng lật trang sổ tay điều tra khi mở Notebook.
  - Nhạc nền điều tra phong cách Noir/Ambient hồi hộp, dồn dập hơn khi tỉ lệ tin giả vượt 60% và 70%.

---

### 4. Phase 12 — Tự động hóa CI/CD & Triển khai Demo
- [ ] **GitHub Actions Workflow** (`.github/workflows/ci.yml`):
  - Kiểm tra Go test, `go vet ./...` và `scenario-check`.
  - Kiểm tra Next.js TypeScript check và `npm run build`.
- [ ] **Docker Compose Production**:
  - Chạy thử nghiệm toàn bộ hệ thống bằng 1 lệnh:
    ```bash
    docker compose up --build
    ```
- [ ] **Triển khai Demo lên Cloud**:
  - Deploy Backend Go lên Render / Fly.io / Railway.
  - Deploy Frontend Next.js lên Vercel / Cloudflare Pages.
- [ ] **Chụp ảnh Screenshots & Quay video Demo**:
  - Bổ sung ảnh chụp giao diện bản đồ 3D, thanh HUD cảnh báo, sổ tay điều tra và màn hình báo cáo After-Action Report vào `README.md`.

---

## III. HƯỚNG DẪN KHỞI CHẠY NHANH NGAY BÂY GIỜ

```bash
# 1. Chạy Backend Go:
cd backend
.\bin\api.exe

# 2. Mở một cửa sổ Terminal khác, chạy Frontend Next.js:
cd frontend
npm run dev

# 3. Mở trình duyệt tại:
http://localhost:3000
```
