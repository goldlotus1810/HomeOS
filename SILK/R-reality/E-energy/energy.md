# E -- ENERGY

symbol:    E (→ U+26A1 in display)
name_vi:   nang luong
branch:    R-reality

## Ngu nghia

Kha nang thuc hien cong. Bat bien trong he kin.
Khong tao ra, khong mat di -- chi chuyen doi hinh thuc.

## Quan he

    E = integral of Power over time
    E → thermodynamics (heat)
    E ≡ information (Landauer: erase 1 bit = kT·ln2 joules)
    E ⊃ all energy forms: kinetic, potential, photon, rest-mass

## Cac dang nang luong

    Kinetic:   E = 0.5 * m * v^2
    Potential: E = m * g * h
    Photon:    E = h * f
    Rest mass: E = m * c^2
    Thermal:   E = k * T

## Projections

C:
    float energy_decay(float E0, float t, float r) {
        return E0 * expf(-r * t);
    }

Rust:
    fn decay(e0: f64, t: f64, r: f64) -> f64 {
        e0 * (-r * t).exp()
    }

Go:
    func Decay(e0, t, r float64) float64 {
        return e0 * math.Exp(-r*t)
    }

Python:
    energy = lambda e0, t, r: e0 * math.exp(-r * t)

## Computing instances

    CPU:  Joules/operation
    GPU:  Watts TDP
    Net:  bits/joule efficiency
    SILK: each symbol lookup has energy cost
