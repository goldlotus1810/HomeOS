# Rust — Memory Safety Without GC

```yaml
symbol:    🛡+C  (C with shield)
born:      2010
creator:   Graydon Hoare / Mozilla
paradigm:  systems, functional, ownership
level:     systems (same as C, safer)
```

## Vị trí trong SILK

```
C → Rust   (Rust = C + 🛡 ownership system)
Rust ⊃ C   (Rust có thể làm mọi thứ C làm + an toàn hơn)
Rust ≈ Go  (cùng era, khác philosophy)
```

## Mapping UTF-32 → Rust

```rust
// ∫ (accumulate / fold)
let sum: f64 = xs.iter().fold(0.0, |acc, x| acc + x);
let sum: f64 = xs.iter().sum();

// ∇ (gradient)
fn grad(f: impl Fn(f64)->f64, x: f64) -> f64 {
    (f(x + 1e-5) - f(x - 1e-5)) / 2e-5
}

// ∪ (union / enum — Rust's killer feature)
enum Shape {
    Sphere { center: Vec3, r: f32 },   // ●
    Capsule { a: Vec3, b: Vec3, r: f32 }, // ⌀
    Void,                               // ◌
}

// 🛡 (ownership = built-in security gate)
fn process(data: String) -> String {  // owns data
    data.to_uppercase()               // transfers ownership
}  // data dropped here automatically

// ∀ (for all / iterator)
xs.iter().for_each(|x| println!("{}", x));
xs.iter().map(|x| x * 2).collect::<Vec<_>>();

// ♻ (cycle / infinite iterator)
(0..).cycle().take(n)

// 🌱 (grow — builder pattern)
let world = World::new()
    .with_seed(888)
    .with_terrain(Terrain::fbm(6))
    .build();

// ∪ smooth union (from HomeOS opTable!)
fn smooth_union(d1: f32, d2: f32, k: f32) -> f32 {
    let h = (k - (d1 - d2).abs()).max(0.0) / k;
    d1.min(d2) - h * h * k * 0.25
}
```

## Core concepts → UTF-32

```
ownership  → 🛡 (security, exactly 1 owner)
borrow &   → ∂ (partial access, temporary)
lifetime   → 𝑡 (time-bounded reference)
Option<T>  → ◌|● (void or value)
Result<T>  → ✓|✗ (ok or error)
trait      → ≡ (behavioral equivalence contract)
enum       → ∪ (one of N variants)
struct     → ⊕ (product of fields)
Box<T>     → ●heap (object on heap)
Arc<T>     → ⊗ (shared ownership, composed)
```
