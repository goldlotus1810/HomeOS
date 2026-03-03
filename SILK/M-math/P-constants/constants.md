# C -- CONSTANTS

symbol:    π (pi / constants)
name_vi:   hang so toan hoc
branch:    M-math

## Ngu nghia

Cac gia tri bat bien cua vu tru.
Khong phu thuoc vao quan sat, van hoa, hay ngon ngu.
Ngon ngu duy nhat thuc su universal.

## Cac hang so

    π = 3.14159265358979...   ty le chu vi / duong kinh
    e = 2.71828182845904...   co so logarithm tu nhien
    φ = 1.61803398874989...   ty le vang (golden ratio)
    √2= 1.41421356237309...   duong cheo hinh vuong don vi
    i = √(-1)                 don vi ao

    Vat ly:
    c  = 299,792,458 m/s      toc do anh sang
    h  = 6.626e-34 J·s        Planck constant
    G  = 6.674e-11 N·m^2/kg^2 hang so hap dan
    kB = 1.380e-23 J/K        Boltzmann constant

## Projections

C:
    #define PI    3.14159265358979323846
    #define E     2.71828182845904523536
    #define PHI   1.61803398874989484820

Rust:
    use std::f64::consts::{PI, E, SQRT_2};
    const PHI: f64 = 1.618033988749895;

Go:
    math.Pi     // π
    math.E      // e
    math.Sqrt2  // √2
    math.Phi    // φ (Go 1.21+)

Python:
    import math
    math.pi    # π
    math.e     # e
    math.tau   # 2π

## Euler Identity

    e^(iπ) + 1 = 0

    Ket hop 5 hang so quan trong nhat:
      e  -- growth, calculus
      i  -- rotation, imaginary
      π  -- circles, geometry
      1  -- multiplicative identity
      0  -- additive identity (void)

    Day la "Hello World" cua toan hoc.
    SILK equivalent: void → everything  (◌ → ∀)
