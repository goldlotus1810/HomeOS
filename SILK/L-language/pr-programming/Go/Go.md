# Go — Simple Concurrency

```yaml
symbol:    ⊗→  (goroutine + channel)
born:      2009
creator:   Rob Pike, Ken Thompson, Robert Griesemer
paradigm:  concurrent, procedural, minimalist
level:     systems-adjacent, high productivity
```

## Vị trí trong SILK

```
C → Go     (Go = C + garbage collection + goroutines)
Go ≈ Rust  (cùng mục đích, khác tradeoffs)
Go ⊃ ⊗    (goroutine = Go's defining concept)
```

## Mapping UTF-32 → Go

```go
// ∫ (accumulate)
sum := 0.0
for _, x := range xs { sum += x }

// ∀ (for all)
for i, x := range xs {
    fmt.Println(i, x)
}

// ⊗ (goroutine = concurrent execution)
go func() {
    result <- compute()
}()

// → (channel = directed message flow)
ch := make(chan int)
go func() { ch <- 42 }()
val := <-ch

// ∪ (interface = behavioral union)
type Shape interface {
    SDF(p Vec3) float64   // ●, ⌀, □ all implement this
    DNA() string
}

// 🛡 (error as value = explicit security gate)
result, err := riskyOperation()
if err != nil {
    return fmt.Errorf("🛡 blocked: %w", err)
}

// ∀ parallel (goroutines = ∀ concurrent)
var wg sync.WaitGroup
for _, gene := range genes {
    wg.Add(1)
    go func(g Gene) {
        defer wg.Done()
        g.Animate(t)      // ∀ genes animate in parallel
    }(gene)
}
wg.Wait()

// 🌱 (grow — state mutation over time)
func (g *Gene) Animate(t float64) {
    g.Age = math.Min(1.0, g.Age + 0.0001*t)
}
```

## Core concepts → UTF-32

```
goroutine  → ⊗ (composed parallel execution)
channel    → → (directed flow of values)
interface  → ≡ (behavioral contract)
struct     → ⊕ (product type)
defer      → 𝑡↩ (execute at end of time scope)
panic      → ✗✗ (unrecoverable error)
nil        → ◌ (void/empty)
map        → ∀→ (for all keys, maps to value)
slice      → ∀∂ (partial view of array)
```
