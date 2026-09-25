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
| **Phase 11** | Cân bằng tự động (Soak Test), Visual Polish & Âm thanh | ✅ **Hoàn thành** | 100-seed soak test PASS (100%), Web Audio engine, strobe lights |
| **Phase 12** | CI/CD GitHub Actions & Demo Deployment | ✅ **Hoàn thành** | GitHub Actions matrix (Go + Next.js build + Soak check) |
| **Phase 13** | Multi-Scenario System & Random Scenario Picker | ✅ **Hoàn thành** | 5 kịch bản hoàn chỉnh (Riverside, Metro Hospital, Midtown Bank, Subway Line 3, City Water), chế độ bốc ngẫu nhiên mỗi lần chơi |


---

## II. CHI TIẾT CÁC HẠNG MỤC CHƯA LÀM & ĐÃ HOÀN THÀNH

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

### 3. Phase 11 — Tinh chỉnh Cân bằng, Visual Polish & Âm thanh (Audio/Visual Polish)
- [x] **Soak Test 100 Seeds tự động**:
  - Đã xây dựng công cụ đa luồng `backend/cmd/soak-test` và các scripts `scripts/soak-test.ps1`, `scripts/soak-test.sh`.
  - Kết quả 100 seeds: 100/100 ván thua tự nhiên (100.0% defeat rate, vượt mục tiêu > 90%), thời gian trung vị 1957.5s (~32:37), tốc độ cực nhanh (11.22 ms/game).
- [x] **Thiết kế Âm thanh Procedural (Web Audio API)**:
  - Bộ tổng hợp âm thanh `frontend/src/services/sound.ts` hoàn toàn không phụ thuộc file mp3 ngoài, 0ms latency.
  - Còi hú xe cứu hỏa / xe cứu thương 2 âm sắc điều tần lúc 18:03 và 18:04.
  - Chuông thông báo 3 nốt rực rỡ khi có tin đồn mới truyền lan trên mạng xã hội.
  - Tiếng sột soạt lật trang sổ tay điều tra khi chuyển tab Notebook.
  - Nhạc nền Noir ambient drone thay đổi độ căng thẳng và tần số cộng hưởng theo tỉ lệ tin giả (vượt 50%, 60%, 70%).
  - Nút bật/tắt âm thanh (Mute/Unmute) trên TopBar HUD.
- [x] **Hiệu ứng Thị giác 3D (Visual Polish)**:
  - Đèn beacon chớp nháy luân phiên đỏ/xanh (stroboscopic strobe) tại cổng Nhà Máy và Bệnh Viện.
  - Phân màu đặc thù theo 6 nhóm vai trò của 20 nhân vật (Cứu hỏa, Y tế, Báo chí, An ninh, Công nhân, Thường dân).
  - Hiệu ứng âm thanh phản hồi xúc giác khi click vào nhân vật hoặc công trình trên bản đồ.

---

### 4. Phase 12 — Tự động hóa CI/CD & Triển khai Demo
- [x] **GitHub Actions Workflow** (`.github/workflows/ci.yml`):
  - Job Backend: `go vet ./...`, `go test -v ./...`, `scenario-check` và Fast Soak Test 10 seeds.
  - Job Frontend: Node 20+, `npm ci`, Next.js 15 production build (`npm run build`).
- [x] **Docker Compose Production**:
  - Đồng bộ `Dockerfile` backend và frontend sẵn sàng chạy với `docker compose up --build`.
- [ ] **Triển khai Demo lên Cloud**:
  - Deploy Backend Go lên Render / Fly.io / Railway.
  - Deploy Frontend Next.js lên Vercel / Cloudflare Pages.
- [ ] **Chụp ảnh Screenshots & Quay video Demo**:
  - Bổ sung ảnh chụp giao diện bản đồ 3D, thanh HUD cảnh báo, sổ tay điều tra và màn hình báo cáo After-Action Report vào `README.md`.

---

### 5. Phase 13 — Hệ Thống Đa Kịch Bản (Multi-Scenario System) & Cân Bằng Toàn Diện (Soak Calibrated)
- [x] **5 Kịch Bản Hoàn Chỉnh (Full Scenarios)**:
  1. `riverside-factory`: Vụ Cháy Kho Nhà Máy Riverside (Tin đồn nổ hóa chất & công ty giấu xác).
  2. `metro-hospital-outbreak`: Báo Động Bệnh Viện Metro (Ngộ độc histamine cá ngừ vs Tin đồn rò rỉ virus phòng lab).
  3. `midtown-bank-run`: Cơn Hoảng Loạn Ngân Hàng Midtown (Deadlock nâng cấp IT vs Tin đồn vỡ nợ, sếp ôm vàng bỏ trốn).
  4. `subway-line3-standstill`: Chuyến Tàu Ngầm Tuyến Số 3 (Chập cáp tín hiệu ray dừng tàu vs Tin đồn khủng bố hơi ngạt chết người).
  5. `city-water-panic`: Khủng Hoảng Nguồn Nước Thành Phố (Sự cố van áp lực sục cặn oxit sắt vs Tin đồn xyanua đầu độc nguồn nước).
- [x] **Kiến Trúc Lan Truyền Tin Đồn Động (Data-driven Belief Seeds & Event Timelines)**:
  - Cấu trúc `BeliefSeed` và `PublicEventSeed` hỗ trợ nhúng hạt giống tin đồn trực tiếp vào sự kiện thời gian thực.
  - Bộ điều phối `Scheduler` tự động phân giải trọng số kích động cảm xúc (Emotional Weights: 0.98 cho tin tử vong, 0.88 cho tin bưng bít, 0.80 cho tin sự cố ban đầu) theo vai trò `primary_false` và `weight` của từng kịch bản.
- [x] **Cân Bằng Tự Động Toàn Diện (100% Soak Test PASS trên cả 5 kịch bản)**:
  - Cả 5 kịch bản đều đạt **100.0% defeat rate** khi không có người chơi can thiệp (vượt chuẩn > 90%).
  - Thời gian thất bại trung vị (Median Time to Defeat) đạt chuẩn từ **31 phút 09 giây** đến **32 phút 23 giây**, đồng nhất tuyệt đối với thiết kế Master Bible.
- [x] **Tích Hợp Giao Diện Frontend (Dynamic Scenario UI)**:
  - Menu chuyển đổi kịch bản mượt mà trên thanh TopBar và màn hình khởi đầu.
  - Chế độ "🎲 Bốc Ngẫu Nhiên" (Random Scenario Picker) cho mỗi lượt chơi mới.
  - Modal nộp sự thật (Submit Truth Modal) tự động tạo form câu hỏi và đáp án động khớp 100% với `truth-form.json` của kịch bản đang chọn.
  - Bản đồ 3D tự động thích ứng tên 12 địa điểm và màu sắc trang phục theo vai trò nhân vật của từng kịch bản.

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
