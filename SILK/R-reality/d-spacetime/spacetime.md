# ∂ — SPACETIME (Không-Thời Gian)

```yaml
symbol:    ∂
codepoint: U+2202
name_en:   PARTIAL DIFFERENTIAL
name_vi:   vi phân riêng / không-thời gian

branch:    ℝ-reality
twig:      ∂-spacetime
```

## Ngữ nghĩa thuần túy

Sự thay đổi theo MỘT chiều trong khi
giữ các chiều khác cố định.
Nền tảng của mọi biến đổi cục bộ.

## Quan hệ

```
∂ ⊂ ∇     — ∇ là tổ hợp của nhiều ∂
∂ → 𝑡     — ∂/∂t = đạo hàm theo thời gian
∂ → pointer(C)  — pointer là ∂ trong lập trình
∂ ≡ borrow(Rust) — borrow là tham chiếu cục bộ
```

## Projections

```c
// C — partial derivative ≡ pointer
int* ptr = &arr[i];  // ∂arr/∂i — partial view

// Rust — borrow = partial access
fn f(x: &Vec3) -> f64 { ... }  // borrow, not own

// Physics
// ∂f/∂x = (f(x+h) - f(x-h)) / 2h
```
