# C — The Foundation Language

```yaml
symbol:    C
born:      1972
creator:   Dennis Ritchie
paradigm:  procedural, imperative
level:     systems (closest to metal)
```

## Vị trí trong SILK

C là ngôn ngữ gần máy nhất vẫn còn đọc được.
Mọi ngôn ngữ khác đều có thể compile xuống C hoặc
compile xuống machine code mà C đã định nghĩa.

```
◌ → 𝕃-language → 💻-programming → C
C → 🔢-machine → x86/ARM   (compile)
C ⊂ Rust                    (Rust supersedes C)
C ⊂ Go                      (Go born from C tradition)
```

## Mapping UTF-32 → C

```c
// ∫ (accumulate)
for(int i=0; i<n; i++) sum += arr[i];

// ∇ (gradient)
double dx = (f(x+h) - f(x-h)) / (2*h);

// ∪ (merge/union)
struct A { int x; };  // union type
union U { int i; float f; };

// ∀ (for all)
for(int i=0; i<n; i++) { /* apply to all */ }

// ∃ (exists)
int found = 0;
for(int i=0; i<n; i++) if(pred(arr[i])) { found=1; break; }

// ● (sphere/ball - geometric primitive)
typedef struct { float x,y,z,r; } Sphere;
float sdf_sphere(vec3 p, Sphere s) {
    return length(p - s.center) - s.r;
}

// 🛡 (security gate)
assert(ptr != NULL);
if(!validate(input)) return ERROR;

// 🌱 (grow/evolve state)
age = fmin(1.0f, age + rate * dt);

// ♻ (cycle)
t = fmod(t + dt, period);
```

## Core concepts → UTF-32

```
pointer    → ∂ (partial reference, points to)
malloc     → ◌→● (void becomes object)
free       → ●→◌ (object returns to void)
struct     → ⊕ (composition of fields)
union      → ∪ (one of many types)
function   → → (maps input to output)
array      → ∀ (collection of same type)
null       → ◌ (void/empty)
```
