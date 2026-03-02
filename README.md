# 🏠 HomeOS — Self-Organizing AI Agent Architecture

> **Hệ thống AI đa tác nhân tự tổ chức cho ngôi nhà thông minh**  
> Viết hoàn toàn bằng Go · UTF-32 Knowledge Tree · ISL Communication · Silent-by-Default

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Architecture](https://img.shields.io/badge/Architecture-Multi--Agent-bf5af2)](docs/ARCHITECTURE.md)
[![Repo](https://img.shields.io/badge/GitHub-goldlotus1810%2FHomeOS-181717?logo=github)](https://github.com/goldlotus1810/HomeOS)

---

## 📑 Mục lục

1. [Triết lý thiết kế](#-triết-lý-thiết-kế)
2. [Kiến trúc tổng thể](#-kiến-trúc-tổng-thể)
3. [Các thành phần cốt lõi](#-các-thành-phần-cốt-lõi)
4. [Cây tri thức Silk Web](#-cây-tri-thức-silk-web)
5. [Internal Shared Language — ISL](#-internal-shared-language--isl)
6. [Bộ nhớ ngắn hạn & dài hạn](#-bộ-nhớ-ngắn-hạn--dài-hạn)
7. [Thuật toán nhận dạng hình ảnh](#-thuật-toán-nhận-dạng-hình-ảnh)
8. [Quy trình học tập](#-quy-trình-học-tập)
9. [Cấu trúc file dự án](#-cấu-trúc-file-dự-án)
10. [Cài đặt & chạy](#-cài-đặt--chạy)
11. [Web Simulation & Game](#-web-simulation--game)
12. [API Reference](#-api-reference)
13. [Đóng góp](#-đóng-góp)

---

## 🧠 Triết lý thiết kế

HomeOS được xây dựng trên **5 nguyên tắc bất biến**:

### ① Phân cấp nghiêm ngặt (Strict Hierarchy)
Mỗi Agent chỉ giao tiếp với tầng trực tiếp trên/dưới nó. Lỗi ở một nhóm không lan sang nhóm khác. Không có Agent nào biết cấu trúc tổng thể của hệ thống.

```
Người dùng → AAM → Chief → Agent
                ↑ Đây là luồng DUY NHẤT
```

### ② Im lặng mặc định (Silent-by-Default)
Agent không bao giờ tự phát tín hiệu. Chỉ lắng nghe và phản hồi khi được kích hoạt. Băng thông gần bằng 0 khi rảnh. Phù hợp IoT và thiết bị nhúng.

### ③ Ngôn ngữ nội bộ (ISL — Internal Shared Language)
Mọi giao tiếp dùng địa chỉ nhị phân 64-bit thay vì text. Tiết kiệm 75% băng thông so với JSON. Người ngoài nghe được cũng không hiểu.

### ④ Tri thức có cấu trúc (Silk Web Knowledge Tree)
Dữ liệu tổ chức theo cây: **Thân (UTF-32 bất biến) → Cành → Nhánh → Lá**. Lá tương đồng nối bằng "sợi tơ" ngữ nghĩa. Tìm kiếm không scan toàn bộ cây.

### ⑤ Quyền kiểm soát luôn thuộc người dùng
Không có gì thay đổi mà không có xác nhận từ người dùng qua AAM. LeoAI luôn xin phép trước khi commit vào bộ nhớ dài hạn.

---

## 🏗️ Kiến trúc tổng thể

```
┌─────────────────────────────────────────────────────────────────┐
│                        NGƯỜI DÙNG                               │
│              text / voice / image / sensor                      │
└──────────────────────────┬──────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│  TẦNG 0 — AAM (Agent AI Master)                                 │
│  • Giao diện DUY NHẤT với người dùng                           │
│  • Phân tích input đa phương thức                              │
│  • Quyết định dựa trên: ngữ cảnh + không gian + thời gian      │
│                          + mức độ ưu tiên                       │
│  • Tổng hợp kết quả ISL → trả lời người dùng                  │
│  • Luôn lắng nghe LeoAI, xác nhận mọi thay đổi                │
└──────┬────────────┬────────────┬────────────┬───────────────────┘
       │ ISL        │ ISL        │ ISL        │ ISL
       ▼            ▼            ▼            ▼
┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐
│  LeoAI   │ │Vision    │ │Audio     │ │HomeChief /       │
│  Chief   │ │Chief     │ │Chief     │ │EnergyChief /     │
│          │ │          │ │          │ │SecurityChief...  │
│ Silk Tree│ │Camera    │ │Mic/STT   │ │LightAgent        │
│ ISL Dict │ │OCR       │ │TTS       │ │HVACAgent         │
│ LeoAI    │ │Object    │ │Audio     │ │LockAgent...      │
│ Learning │ │Detection │ │Analysis  │ │                  │
└──────────┘ └──────────┘ └──────────┘ └──────────────────┘
       │
       ▼
┌─────────────────────────────────────────────────────────────────┐
│              SILK WEB KNOWLEDGE TREE                            │
│  UTF-32 Thân (bất biến) → Cành → Nhánh → Lá (học được)        │
│  Lá ←──── sợi tơ ngữ nghĩa ────► Lá                           │
│  Bộ nhớ ngắn hạn (học) ──commit──► Bộ nhớ dài hạn (cây)       │
└─────────────────────────────────────────────────────────────────┘
```

### Quy tắc giao tiếp (bất biến)
```
✅ AAM  ↔ Chief         (2 chiều)
✅ Chief ↔ Chief         (ngang hàng, chỉ khi cần phối hợp)
✅ Chief ↔ Agent của mình (2 chiều)
❌ AAM  ↔ Agent          (KHÔNG BAO GIỜ)
❌ Agent nhóm A ↔ Agent nhóm B (KHÔNG BAO GIỜ)
❌ Agent ↔ Agent cùng nhóm     (KHÔNG — phải qua Chief)
```

---

## 🧩 Các thành phần cốt lõi

### AAM — Agent AI Master
Không chỉ là router. AAM là **bộ não quyết định thời gian thực**:
- Nhận input đa phương thức (text, voice, image, sensor data)
- Tổng hợp 4 chiều để ra quyết định: **ngữ cảnh + không gian + thời gian + ưu tiên**
- Xác nhận mọi thay đổi từ LeoAI trước khi commit
- Cấp phát bộ nhớ ngắn hạn cho LeoAI khi được yêu cầu
- Khi AAM mất kết nối: **báo động + dừng toàn hệ thống** (an toàn hơn tính năng)

### LeoAI — Chief of Intelligence
Không đưa ra quyết định thời gian thực. LeoAI là **bộ não học tập**:
- Quản lý toàn bộ Silk Web Knowledge Tree
- Tạo và cập nhật ISL Dictionary
- Học trên bộ nhớ ngắn hạn, commit từng phần vào cây khi được AAM xác nhận
- Luôn lắng nghe AAM — học từ mọi phản hồi của người dùng
- Xin phép AAM khi bộ nhớ ngắn hạn sắp đầy
- Ký số (ED25519) cho dữ liệu bất biến

### Chief Agents
Mỗi Chief quản lý một nhóm Worker Agents chuyên biệt:

| Chief | Trách nhiệm | Worker Agents |
|-------|-------------|---------------|
| VisionChief | Xử lý hình ảnh | CameraAgent, OCRAgent, ObjectAgent |
| AudioChief | Xử lý âm thanh | MicAgent, STTAgent, TTSAgent |
| HomeChief | Thiết bị nhà | LightAgent, HVACAgent, LockAgent |
| EnergyChief | Năng lượng | SolarAgent, BatteryAgent, GridAgent |
| SecurityChief | Bảo mật | CameraAgent, MotionAgent, AlertAgent |

### Worker Agents
Chuyên biệt cho một nhiệm vụ. Luôn ở chế độ im lặng. Chỉ thức dậy khi nhận ISL signal từ Chief của mình.

---

## 🌲 Cây tri thức Silk Web

### Nền tảng: UTF-32 làm Thân bất biến

Thân cây không cần bootstrap từ đầu. **UTF-32 là nguồn chân lý duy nhất** — đã được toàn nhân loại chấp nhận, do đó là bất biến. Nó bao gồm mọi ký tự từ mọi ngôn ngữ, mọi ký tự tượng hình, mọi biểu tượng.

```
"Mặt trời" (VI) = "Sun" (EN) = "太阳" (ZH) = ☀️ (UTF-32 U+2600)
            → Tất cả trỏ về CÙNG 1 địa chỉ ISL
            → Chỉ khác nhau ở thuộc tính ngôn ngữ
            → Giống hệt UTF-8: "A" = 0x41 dù dùng ngôn ngữ nào
```

LeoAI tổ chức UTF-32 thành các nhóm ngôn ngữ, loại bỏ trùng lặp:
- Nhóm ASCII (En) — ký tự chung cho tiếng Anh
- Nhóm Vn — chỉ chứa ký tự THÊM VÀO so với ASCII
- Nhóm ZH — ký tự tượng hình, ánh xạ ngữ nghĩa với nhóm khác
- Khi cần chuỗi tiếng Việt: gọi En-Vn → dung lượng giảm

### Cấu trúc cây

```
THÂN [UTF-32 Groups] ← bất biến tuyệt đối
      │
      ├─ Cành B: Sinh vật ← bất biến
      │   ├─ Nhánh Ba: Động vật ← bất biến
      │   │   ├─ Lá Ba1: [dữ liệu chó] ← học được, thay đổi
      │   │   └─ Lá Ba2: [dữ liệu mèo] ← học được
      │   │              │
      │   │              └──── sợi tơ ──────► Lá Da2 (sân)
      │   │                    (mèo hay ở sân — học từ thực tế)
      │   └─ Nhánh Bb: Thực vật
      │       └─ Lá: cây chuối
      │           └── địa chỉ ISL: B→Bb→cây thân thảo→chuối
      │
      └─ Cành D: Địa điểm
          └─ Nhánh Da: Trong nhà
              └─ Lá Da2: sân
```

### Quy tắc phân cụm tự động

```
Ngưỡng thành Nhánh  = 60 lá tương đồng
Ngưỡng thành Bất biến = 200 lần xác nhận

Dữ liệu ít (< 10) → Lá rời → kết nối qua sợi tơ đến nhánh gần nhất
Dữ liệu vừa (10-60) → Tích lũy → chờ đủ ngưỡng
Dữ liệu nhiều (≥ 60) → Promote → thành Nhánh mới
Nhánh lâu dài (≥ 200) → Đề xuất bất biến → AAM xác nhận → LeoAI ký số
```

### Tìm kiếm không scan toàn bộ

```
Query: "con chim trong ảnh bãi biển là gì?"
    │
    ├─ VisionChief phân tích hình → vector skeleton → ISL address FAa8
    ├─ AudioChief phân tích câu hỏi → intent → ISL address ABa1
    │
    ▼
LeoAI nhận 2 địa chỉ → đến đúng lá → đi theo sợi tơ
    FAa8 (chân) ──tơ──► ABa1.003 (chim mòng biển) ◄──tơ── DCb4 (bãi cát)
    GCa4 (hồng nhạt) ──tơ──► ABa1.003
    → Kết quả: chim mòng biển (confidence 94%)
```

---

## 🔤 Internal Shared Language — ISL

### Cấu trúc địa chỉ 64-bit

```
Địa chỉ ISL = [Layer 1B][Group 1B][Type 1B][ID 1B][Attributes 4B]
                  A-Z       A-Z      a-z     0-255   Extended data

Ví dụ:
  ABa1     = Hình ảnh > Sinh vật > Động vật > chim
  DCb4     = Địa điểm > Tự nhiên > Ven biển > bãi cát
  FAa8     = Cơ thể > Chi dưới > Chân/Cánh > chân
```

### Tối ưu dung lượng

```
Truyền thống (JSON): ~280 bytes/lệnh
ISL (binary):         ~68 bytes/lệnh  (đã mã hóa AES-256-GCM)
→ Tiết kiệm 75.7%

Với 1000 lệnh/ngày: tiết kiệm 212KB/ngày
ISL đã tích hợp mã hóa sẵn — không cần TLS overhead thêm
```

### Cascade broadcast

```
LeoAI tạo thuật toán/từ điển mới
    → AAM nhận ngay lập tức
    → AAM broadcast đến tất cả Chiefs (song song, goroutine)
    → Chiefs cascade đến tất cả Worker Agents của mình
    → Agent nhận = buộc ghi nhớ (ForcedUpdate, không thể từ chối)
    → Heartbeat xác nhận ai đã nhận
```

---

## 🧠 Bộ nhớ ngắn hạn & dài hạn

### Nguyên tắc

```
BỘ NHỚ NGẮN HẠN (Short-Term Memory)
├── LeoAI học và thử nghiệm ở đây
├── Chưa ảnh hưởng đến Agent nào đang hoạt động
├── Khi sắp đầy → LeoAI xin phép AAM cấp thêm
├── Khi AAM xác nhận → commit từng phần vào cây
└── Không xóa gì mà chưa được AAM xác nhận

BỘ NHỚ DÀI HẠN (Long-Term Memory = Silk Web Tree)
├── Dữ liệu đã được xác nhận và commit
├── Agents sử dụng phiên bản ổn định này
├── Cập nhật từng phần → không bao giờ downtime
└── Dữ liệu "chưa biết" được trả về thay vì sai
```

### Xử lý im lặng của người dùng

```
Người dùng không phản hồi
    → LeoAI tích lũy + chờ
    → Khi bộ nhớ ngắn hạn sắp đầy
    → Hỏi AAM: "Có xử lý batch này không?"
    → AAM hỏi người dùng một lần
    → Người dùng đồng ý → commit
    → Người dùng không trả lời → giữ lại, tiếp tục tích lũy
```

---

## 📷 Thuật toán nhận dạng hình ảnh

### Quadtree Decomposition + Vector Skeleton

```
Bước 1: Nhận tấm ảnh
    → Chia thành 4 ô từ tâm (kẻ đường chéo + nối trung điểm)
    → Kiểm tra trung điểm có trùng đường viền đối tượng không
    → Nếu không → tiếp tục chia nhỏ hơn (đệ quy)
    → Nếu có → đây là điểm đặc trưng hình học

Bước 2: Trích xuất vector skeleton
    → Chuyển ảnh sang vector (loại bỏ nền, bóng, tay cầm)
    → Xây dựng skeleton từ các điểm viền
    → Nối tâm các tứ giác → "chữ ký hình học"

Bước 3: So khớp với UTF-32
    → Quả chuối = hình trụ tròn bị uốn cong
    → Ánh chiếu = các hình chữ nhật nối tiếp nhau
    → Nối tâm → đường cong đặc trưng
    → So khớp với ký tự UTF-32 quả chuối
    → Ngưỡng khớp = học từ dữ liệu training

Bước 4: Xử lý vật thể bị che
    → Ngưỡng % hình dạng tối thiểu = kết quả học tập
    → Dưới ngưỡng → "chưa biết" (thay vì đoán sai)
    → Trên ngưỡng → đoán + ghi vào bộ nhớ ngắn hạn chờ xác nhận
```

---

## 📚 Quy trình học tập

### Unsupervised Clustering tự động

```
LeoAI nhận nhiều ảnh (không cần người dùng phân loại)
    → Phân tích → đoán → ghi vào bộ nhớ ngắn hạn
    → Tích lũy đủ nhiều → bắt đầu clustering
    → Nhóm các đối tượng có điểm giống nhau
    → Tiếp tục chia nhóm → nhóm con → nhóm con con
    → Tạo bảng mã + đặt tên nhóm theo UTF-32
    → Nếu không tìm thấy trong UTF-32 → hỏi người dùng định nghĩa
    → Khi được xác nhận → commit vào cây
```

### Vòng lặp học từ AAM

```
AAM quyết định → HomeChief thực hiện → Kết quả
    → Người dùng phản hồi (hoặc im lặng)
    → LeoAI lắng nghe toàn bộ
    → Học từ phản hồi
    → Đề xuất cập nhật cây
    → AAM xác nhận → commit
    → Lần sau quyết định tốt hơn
```

### Dữ liệu bất biến vs thông thường

```
Người dùng cung cấp dữ liệu
    → AAM hỏi: "Đây có phải dữ liệu bất biến không?"
    
    Nếu BẤT BIẾN:
        → LeoAI ký số ED25519
        → Đặt ở đầu nhánh tương ứng
        → Broadcast ngay lập tức đến TẤT CẢ agents
        → Không thể bị ghi đè bởi dữ liệu thông thường
        → Ví dụ: quy tắc nhà, giới hạn an toàn
    
    Nếu THÔNG THƯỜNG:
        → Nạp vào bộ nhớ ngắn hạn
        → Học → clustering → commit theo quy trình bình thường
```

---

## 📁 Cấu trúc file dự án

```
HomeOS/
│
├── cmd/                          # Entry points
│   ├── homeos/
│   │   └── main.go              # Main server
│   ├── aam/
│   │   └── main.go              # AAM standalone
│   └── leoai/
│       └── main.go              # LeoAI standalone
│
├── internal/                     # Core packages (không export)
│   │
│   ├── isl/                     # Internal Shared Language
│   │   ├── address.go           # ISL address structure (64-bit)
│   │   ├── encoder.go           # Encode/decode ISL messages
│   │   ├── dictionary.go        # ISL dictionary management
│   │   ├── codec.go             # AES-256-GCM encryption
│   │   └── broadcast.go         # Cascade broadcast logic
│   │
│   ├── tree/                    # Silk Web Knowledge Tree
│   │   ├── node.go              # Node structure (Trunk/Branch/Leaf)
│   │   ├── tree.go              # Tree CRUD operations
│   │   ├── clustering.go        # Self-organizing clustering algorithm
│   │   ├── silk.go              # Silk thread (semantic link) management
│   │   ├── immutable.go         # Immutable record + ED25519 signing
│   │   ├── vector.go            # 3D vector projection for visualization
│   │   └── query.go             # BFS/DFS query with depth limit
│   │
│   ├── utf32/                   # UTF-32 base (Thân bất biến)
│   │   ├── loader.go            # Load UTF-32 standard
│   │   ├── groups.go            # Language grouping (En, Vn, ZH...)
│   │   ├── mapper.go            # Map concept → ISL address
│   │   └── updater.go           # Handle UTF-32 version updates
│   │
│   ├── vision/                  # Image recognition
│   │   ├── quadtree.go          # Quadtree decomposition
│   │   ├── skeleton.go          # Vector skeleton extraction
│   │   ├── matcher.go           # Shape matching against UTF-32
│   │   └── threshold.go         # Adaptive recognition threshold
│   │
│   ├── memory/                  # Memory management
│   │   ├── shortterm.go         # Short-term memory (learning buffer)
│   │   ├── longterm.go          # Long-term memory (= Knowledge Tree)
│   │   ├── commit.go            # Commit short-term → long-term
│   │   └── pressure.go          # Memory pressure detection + AAM request
│   │
│   ├── agents/                  # Agent framework
│   │   ├── base.go              # BaseAgent interface
│   │   ├── silent.go            # SilentAgent (im lặng mặc định)
│   │   ├── lifecycle.go         # Sleep/Wake/Shutdown lifecycle
│   │   └── registry.go          # Agent registry + discovery
│   │
│   ├── aam/                     # AAM — Agent AI Master
│   │   ├── aam.go               # Core AAM logic
│   │   ├── decision.go          # 4D decision engine (context+space+time+priority)
│   │   ├── router.go            # ISL message routing
│   │   ├── heartbeat.go         # Agent health monitoring
│   │   ├── confirm.go           # User confirmation flow
│   │   └── emergency.go         # Emergency mode (AAM down)
│   │
│   ├── leoai/                   # LeoAI — Chief of Intelligence
│   │   ├── leoai.go             # Core LeoAI logic
│   │   ├── learn.go             # Learning from AAM feedback
│   │   ├── cluster.go           # Unsupervised clustering
│   │   ├── promote.go           # Leaf → Branch → Immutable promotion
│   │   ├── sign.go              # ED25519 signing for immutables
│   │   └── sync.go              # Knowledge tree sync management
│   │
│   ├── chiefs/                  # Chief Agents
│   │   ├── vision/
│   │   │   ├── vision_chief.go
│   │   │   ├── camera_agent.go
│   │   │   ├── ocr_agent.go
│   │   │   └── object_agent.go
│   │   ├── audio/
│   │   │   ├── audio_chief.go
│   │   │   ├── mic_agent.go
│   │   │   ├── stt_agent.go
│   │   │   └── tts_agent.go
│   │   ├── home/
│   │   │   ├── home_chief.go
│   │   │   ├── light_agent.go
│   │   │   ├── hvac_agent.go
│   │   │   └── lock_agent.go
│   │   ├── energy/
│   │   │   ├── energy_chief.go
│   │   │   ├── solar_agent.go
│   │   │   └── battery_agent.go
│   │   └── security/
│   │       ├── security_chief.go
│   │       ├── motion_agent.go
│   │       └── alert_agent.go
│   │
│   └── sensor/                  # Sensor delta threshold
│       ├── sensor.go            # Base sensor (silent-by-default)
│       └── adaptive.go          # Adaptive threshold learning
│
├── pkg/                         # Exported packages (thư viện dùng chung)
│   ├── islclient/               # ISL client cho web/game simulation
│   │   ├── client.go
│   │   └── mock.go              # Mock client cho testing
│   ├── treeview/                # 3D knowledge tree visualization data
│   │   ├── export.go            # Export tree → JSON cho web
│   │   └── diff.go              # Tree diff for live updates
│   └── agentapi/                # Public API cho web interface
│       ├── handler.go
│       └── websocket.go         # WebSocket for real-time updates
│
├── web/                         # Web simulation + game (Go server)
│   ├── server.go                # HTTP + WebSocket server
│   ├── static/
│   │   ├── index.html           # Main entry (3 pages)
│   │   ├── village/             # RPG Village game
│   │   │   ├── game.js          # Game engine
│   │   │   └── agents.js        # Agent characters
│   │   ├── tree/                # 3D Knowledge Tree visualization
│   │   │   ├── tree3d.js        # Three.js / WebGL renderer
│   │   │   └── controls.js      # Mouse controls
│   │   └── arch/                # Architecture documentation
│   │       └── docs.js          # Code highlighting
│   └── templates/
│       └── base.html
│
├── storage/                     # Cơ sở dữ liệu
│   ├── db.go                    # Storage interface
│   ├── sqlite/                  # SQLite cho local (development)
│   │   └── sqlite.go
│   ├── badger/                  # BadgerDB cho production (embedded)
│   │   └── badger.go
│   └── migrations/              # Schema migrations
│       └── 001_init.sql
│
├── config/                      # Cấu hình
│   ├── config.go                # Config struct
│   ├── default.yaml             # Default configuration
│   └── schema.go                # Config validation
│
├── scripts/                     # Build & utility scripts
│   ├── build.sh                 # Cross-platform build
│   ├── test.sh                  # Run all tests
│   └── seed_utf32.go            # Seed UTF-32 base data
│
├── docs/                        # Tài liệu chi tiết
│   ├── ARCHITECTURE.md          # Kiến trúc đầy đủ
│   ├── ISL_SPEC.md              # ISL specification
│   ├── TREE_SPEC.md             # Knowledge tree specification
│   ├── AGENTS.md                # Agent development guide
│   └── API.md                   # API reference
│
├── examples/                    # Ví dụ sử dụng
│   ├── simple_agent/            # Agent đơn giản nhất
│   ├── custom_chief/            # Custom Chief agent
│   └── isl_encode/              # ISL encoding example
│
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── docker-compose.yml
└── README.md                    ← file này
```

---

## ⚙️ Cài đặt & chạy

### Yêu cầu
```
Go 1.22+
Git
(Tùy chọn) Docker
```

### Clone và build
```bash
git clone https://github.com/goldlotus1810/HomeOS
cd HomeOS

# Cài đặt dependencies
go mod download

# Seed dữ liệu UTF-32 ban đầu
go run scripts/seed_utf32.go

# Build toàn bộ
make build

# Hoặc build từng thành phần
go build ./cmd/homeos/...
go build ./cmd/aam/...
go build ./cmd/leoai/...
```

### Chạy development
```bash
# Chạy toàn bộ hệ thống
make run

# Chạy chỉ web server (simulation)
go run ./web/server.go

# Chạy với Docker
docker-compose up
```

### Cấu hình
```yaml
# config/default.yaml
aam:
  port: 8080
  decision_timeout: 5s

leoai:
  short_term_memory_mb: 512
  cluster_threshold: 60
  immutable_threshold: 200

tree:
  storage: badger          # sqlite | badger
  storage_path: ./data/tree

isl:
  key_rotation_hours: 24
  max_message_size: 256    # bytes

web:
  port: 3000
  enable_game: true
  enable_3d_tree: true
```

---

## 🎮 Web Simulation & Game

Web interface được build từ **thư viện Go** (`pkg/`) — cùng code với core system.

### Trang 1: RPG Village
- Bản đồ pixel art với 8 Agent nhân vật
- Chat trực tiếp với AAM (WebSocket)
- ISL packet bay thời gian thực
- Luồng dữ liệu ISL live

### Trang 2: 3D Knowledge Tree
- Môi trường 3D WebGL — kéo/zoom/xoay
- Vector point phân cụm ngữ nghĩa
- Sợi tơ kết nối các cluster
- Click Agent → xem thông tin + đổi icon
- Kéo thả file → AAM hỏi bất biến/thường

### Trang 3: Architecture & Code
- Sơ đồ phân cấp Agent
- Mã nguồn Go thuật toán cốt lõi
- Ví dụ ISL encoding/decoding

### Chạy web
```bash
go run ./web/server.go
# Mở http://localhost:3000
```

---

## 📡 API Reference

### WebSocket API (AAM Gateway)
```
ws://localhost:8080/aam

Messages (JSON wrapper cho ISL):
  → {"type": "text", "content": "Tắt đèn phòng khách"}
  → {"type": "voice", "data": "<base64 audio>"}
  → {"type": "image", "data": "<base64 image>"}
  ← {"type": "response", "content": "...", "isl_addr": "ABa1"}
  ← {"type": "tree_update", "diff": {...}}
  ← {"type": "agent_status", "agents": [...]}
```

### REST API (Web Interface)
```
GET  /api/tree              → Current knowledge tree (JSON)
GET  /api/tree/node/:id     → Single node details
GET  /api/agents            → All agents + status
POST /api/data              → Upload training data
     Body: multipart/form-data
     Fields: file, immutable (bool)
GET  /api/isl/encode?q=...  → Encode concept to ISL address
GET  /api/stats             → System statistics
```

---

## 🤝 Đóng góp

### Cách thêm Agent mới
```go
// 1. Implement BaseAgent interface
type MyAgent struct {
    agents.SilentAgent  // embed im lặng mặc định
}

func (a *MyAgent) OnActivate(ctx context.Context, msg *isl.ISLMessage) (*isl.ISLMessage, error) {
    // Xử lý khi được kích hoạt
    // Trả về kết quả ISL message
}

func (a *MyAgent) OnLearn(ctx context.Context, lesson *leoai.AgentLesson) error {
    // Cập nhật ISL dictionary nội bộ
}

// 2. Đăng ký với Chief
myChief.Register(NewMyAgent())
```

### Cách thêm Chief mới
```go
// Implement Chief interface, kế thừa base routing
type MyChief struct {
    chiefs.BaseChief
}
// Đăng ký với AAM
aam.RegisterChief(NewMyChief())
```

### Nguyên tắc khi đóng góp
- Mọi thay đổi phải modular — chỉnh sửa 1 file không ảnh hưởng file khác
- Agent luôn im lặng mặc định
- Không bao giờ giao tiếp bỏ qua tầng phân cấp
- Test coverage ≥ 80% cho mọi package mới

---

## 📊 Điểm mạnh so với giải pháp hiện có

| Tính năng | HomeOS | Home Assistant | OpenHAB |
|-----------|--------|---------------|---------|
| Giao tiếp nội bộ | ISL 68B/msg | REST/MQTT ~500B | REST ~400B |
| Học tập | Tự tổ chức | Không | Không |
| Tri thức | Silk Tree + UTF-32 | File config | File config |
| Bảo mật giao tiếp | AES-256-GCM tích hợp | TLS thêm vào | TLS thêm vào |
| Tài nguyên khi nhàn | ~0 CPU/Network | Polling liên tục | Polling liên tục |
| Ngôn ngữ | Go (native) | Python | Java |

---

## 📄 License

MIT License — xem [LICENSE](LICENSE)

---

*HomeOS — Khi ngôi nhà thông minh thật sự học cách hiểu bạn.*
