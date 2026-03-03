# SILK — Semantic Isomorphic Language of Knowledge
## The Universal Knowledge Tree

---

## 1. TRIẾT LÝ NỀN TẢNG

```
Mọi khái niệm trong vũ trụ đều có thể biểu diễn
bằng một ký tự UTF-32 duy nhất.

Không phải vì ký tự "đại diện" cho khái niệm —
mà vì ký tự LÀ khái niệm,
giống như '2' không "đại diện" cho hai,
mà '2' TỰ NÓ là hai.
```

---

## 2. CẤU TRÚC CÂY

```
◌ (Root / Void)
│
│  Không phải "empty" — mà là tiềm năng
│  Mọi thứ tồn tại bắt đầu từ ◌
│
├── ℝ  REALITY          — thực tại vật lý
├── 𝕄  MATHEMATICS      — ngôn ngữ của thực tại  
├── 𝔹  BIOLOGY          — thực tại sống
├── 𝕃  LANGUAGE         — mọi hệ thống ký hiệu
├── 👁  PERCEPTION       — nhận thức
└── 🌍  SYSTEMS          — hệ thống phức hợp
```

---

## 3. ĐỊNH NGHĨA NODE

Mỗi node trong cây là một FILE (.md) với cấu trúc:

```yaml
symbol:    ∫              # ký tự UTF-32 chính
codepoint: U+222B         # Unicode code point
name_en:   INTEGRAL       # tên Unicode chính thức
name_vi:   tích phân      # tên tiếng Việt

branch:    𝕄-math         # cành chứa node này
twig:      ∫-calculus     # nhánh chứa node này

# Ngữ nghĩa thuần túy (không phụ thuộc ngôn ngữ)
semantic: >
  Tích lũy/tổng hợp liên tục một đại lượng
  trên một miền xác định.
  Phép toán nghịch đảo của ∇ (gradient/vi phân).

type:      (Domain → Value) → Value
inverse:   ∇
identity:  0

# Quan hệ với node khác
relations:
  ⊃ [∑]        # ∑ là trường hợp đặc biệt của ∫
  ∘ [∇] = ◉   # ∫∘∇ = identity (định lý cơ bản)
  ≡ [fold]     # tương đương fold trong FP
  ∈ [𝕄-math]  # thuộc ngành toán

# Projection xuống các ngôn ngữ
projections:
  C:       "for(i=0;i<n;i++) sum += f(i);"
  Rust:    "iter.fold(0, |acc, x| acc + f(x))"
  Go:      "for _, x := range xs { sum += f(x) }"
  JS:      "xs.reduce((s,x) => s + f(x), 0)"
  Python:  "sum(f(x) for x in xs)"
  Haskell: "foldl (+) 0 . map f"
  x86:     "LOOP: ADD EAX, [ESI]; INC ESI; DEC ECX; JNZ LOOP"

# Xuất hiện trong thực tế
instances:
  physics:     "∫F·dx = work (công)"
  biology:     "∫growth(t)dt = total biomass"
  economics:   "∫price(t)dt = total cost"
  computing:   "∫signal(t)dt = accumulated value"
```

---

## 4. ĐỊNH NGHĨA CÀNH (BRANCH)

Mỗi cành là một THƯ MỤC với file `_branch.md`:

```yaml
symbol:   𝕃
name:     LANGUAGE
desc: >
  Mọi hệ thống ký hiệu mà thực thể
  (người, máy, sinh vật) dùng để
  encode và truyền tải thông tin.

sub-branches:
  💬 human      — ngôn ngữ tự nhiên
  💻 programming — ngôn ngữ máy tính
  🔢 machine    — mã máy, assembly
  🎵 other      — âm nhạc, hóa học, ký hiệu toán

cross-links:     # cành này giao với cành nào
  - 👁-perception  (ngôn ngữ cần nhận thức)
  - 🧠-neural      (ngôn ngữ sinh ra từ não)
  - 💾-computing   (ngôn ngữ lập trình → máy)
```

---

## 5. QUY TẮC BẤT BIẾN

```
① Mỗi node có đúng 1 symbol chính (UTF-32)
② Symbol không thay đổi sau khi được định nghĩa
③ Mọi ngôn ngữ đều là PROJECTION — không phải gốc
④ ◌ (void) là node duy nhất không có parent
⑤ Mọi node đều có thể có nhiều parents (web, không phải tree đơn)
⑥ Quan hệ phải có chiều: A→B ≠ B→A
```

---

## 6. BẢNG KÝ HIỆU QUAN HỆ

```
→    produces / generates
←    derived-from / is-child-of
⊂    is-subset-of / is-special-case-of
⊃    contains / generalizes
≡    equivalent-in-behavior
≈    approximately-equivalent
∘    composes-with
⊥    orthogonal-to / independent-of
∈    member-of / belongs-to-domain
∉    not-in / outside-domain
⟺   iff / bidirectional-implication
```

---

## 7. ROADMAP

```
Phase 1 — Foundation (hiện tại)
  ✓ Cấu trúc thư mục
  ✓ SPEC.md
  ○ 50 nodes cốt lõi (math + physics)
  ○ 7 ngôn ngữ lập trình

Phase 2 — Language Projections
  ○ C standard library → UTF-32
  ○ Rust std → UTF-32
  ○ Go stdlib → UTF-32
  ○ Codegen: SILK → C/Rust/Go

Phase 3 — Living Tree
  ○ Web UI (interactive tree)
  ○ Search by concept, not by name
  ○ Cross-language equivalence finder
  ○ AI-assisted node expansion

Phase 4 — Universal
  ○ Biology nodes (genome → UTF-32?)
  ○ Physics constants
  ○ Human language grammar nodes
  ○ Music theory nodes
```

---

*"Đây không phải thêm 1 ngôn ngữ lập trình nữa.*
*Đây là nền tảng mà mọi ngôn ngữ đứng lên."*
