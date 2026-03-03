# T -- THERMODYNAMICS

symbol:    T (temperature, thermodynamics)
name_vi:   nhiet dong luc hoc
branch:    R-reality

## Ngu nghia

Entropy luon tang. He thong tien ve trang thai hon loan nhat.
Thong tin va nhiet la hai mat cua cung mot thuc tai.

## 4 Dinh luat

    0th: Equilibrium is transitive  (A=B, B=C → A=C)
    1st: ΔU = Q - W                 (energy conserved)
    2nd: ΔS ≥ 0                     (entropy never decreases)
    3rd: S → 0 as T → 0K

## Entropy = Information

    Shannon:        H = -∑ p(x)·log2(p(x))   [bits]
    Thermodynamic:  S = k·ln(Ω)              [J/K]

    Landauer: erase 1 bit = k·T·ln(2) joules
              → Information IS physical
              → Computation has energy cost

Python:
    def entropy(probs):
        import math
        return -sum(p * math.log2(p) for p in probs if p > 0)
    # uniform [.25,.25,.25,.25] → 2.0 bits (max entropy)

Go:
    func Entropy(states int) float64 {
        k := 1.380649e-23
        return k * math.Log(float64(states))
    }

## SILK connection

    Adding a node to SILK reduces entropy of knowledge space.
    SILK = entropy reduction engine for human knowledge.
