# ∇ — NABLA / GRADIENT (Vi Phân / Độ Dốc)

```yaml
symbol:    ∇
codepoint: U+2207
name_en:   NABLA
name_vi:   nabla / gradient / vi phân

branch:    𝕄-math
twig:      ∇-geometry
```

## Ngữ nghĩa thuần túy

Hướng và tốc độ thay đổi của một trường tại một điểm.
"Nơi nào dốc nhất, theo hướng nào."

## Quan hệ

```
∇ ⊃ ∂        — ∂ (partial) là thành phần của ∇
∇ ∘ ∫ = ◉   — vi phân của tích phân = identity
∇ ≡ diff     — trong numerical computing
∇ → backprop — trong neural networks (gradient descent)
```

## Projections

```c
// C — finite difference approximation
double grad = (f(x + h) - f(x - h)) / (2 * h);
```

```rust
// Rust
fn gradient(f: impl Fn(f64)->f64, x: f64, h: f64) -> f64 {
    (f(x + h) - f(x - h)) / (2.0 * h)
}
```

```python
# Python / NumPy
grad = np.gradient(field)
```

```js
// JavaScript
const grad = (f, x, h=1e-5) => (f(x+h) - f(x-h)) / (2*h);
```

## Instances

```
Physics:    ∇T = heat flow direction
ML:         ∇Loss = backpropagation
Graphics:   ∇height = surface normal (HomeOS!)
Economics:  ∇profit = marginal gain direction
```
