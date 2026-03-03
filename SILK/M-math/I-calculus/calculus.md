# ∫ — INTEGRAL (Tích Phân)

```yaml
symbol:    ∫
codepoint: U+222B
name_en:   INTEGRAL
name_vi:   tích phân / tích lũy

branch:    𝕄-math
twig:      ∫-calculus
```

## Ngữ nghĩa thuần túy

Tích lũy liên tục một đại lượng trên một miền.
Phép toán nghịch của ∇ (vi phân).

## Quan hệ

```
∫ ⊃ ∑        — ∑ là trường hợp rời rạc của ∫
∫ ∘ ∇ = ◉   — tích phân của vi phân = identity
∫ ≡ fold     — trong functional programming
∫ ≡ reduce   — trong Python/JS
∫ ∈ 𝕄-math
```

## Projections

```c
// C
double sum = 0;
for(int i = 0; i < n; i++) sum += f(i);
```

```rust
// Rust
let sum = xs.iter().fold(0.0, |acc, x| acc + f(x));
```

```go
// Go
sum := 0.0
for _, x := range xs { sum += f(x) }
```

```js
// JavaScript
const sum = xs.reduce((s, x) => s + f(x), 0);
```

```python
# Python
sum(f(x) for x in xs)
```

```haskell
-- Haskell
foldl (+) 0 (map f xs)
```

```asm
; x86 Assembly
; EAX=sum, ESI=array ptr, ECX=count
XOR EAX, EAX
LOOP: ADD EAX, [ESI]
      ADD ESI, 4
      DEC ECX
      JNZ LOOP
```

## Instances trong thực tế

```
Physics:    ∫F·dx = work (công cơ học)
Biology:    ∫growth(t)dt = total biomass
Economics:  ∫price(t)dt = total cost
Computing:  ∫signal(t)dt = digital integration
```
