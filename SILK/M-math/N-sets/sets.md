# S -- SET THEORY

symbol:    ∞ (infinity / sets)
name_vi:   ly thuyet tap hop / vo han
branch:    M-math

## Ngu nghia

Tap hop: bo suu tap cac phan tu phan biet.
Nen tang cua toan bo toan hoc hien dai.

## Ky hieu → UTF-32 mapping

    ∈   element of          (thuoc)
    ∉   not element of      (khong thuoc)
    ⊂   subset              (tap con)
    ⊃   superset            (tap cha)
    ∪   union               (hop)
    ∩   intersection        (giao)
    ∅   empty set           (tap rong = void)
    ∀   for all             (moi phan tu)
    ∃   there exists        (ton tai)
    ℕ   natural numbers
    ℤ   integers
    ℚ   rationals
    ℝ   reals
    ℂ   complex numbers

## Projections

C (set as bit array):
    uint64_t set = 0;
    set |= (1ULL << elem);       // add elem
    set &= ~(1ULL << elem);      // remove elem
    int has = (set >> elem) & 1; // membership test

Rust:
    use std::collections::HashSet;
    let mut s: HashSet<i32> = HashSet::new();
    s.insert(42);
    let union: HashSet<_> = a.union(&b).collect();
    let inter: HashSet<_> = a.intersection(&b).collect();

Python:
    a = {1, 2, 3}
    b = {2, 3, 4}
    a | b    # union
    a & b    # intersection
    a - b    # difference
    a <= b   # subset check

Go:
    type Set[T comparable] map[T]struct{}
    func (s Set[T]) Add(v T)      { s[v] = struct{}{} }
    func (s Set[T]) Has(v T) bool { _, ok := s[v]; return ok }

## Cantor's Infinities

    |N| = countably infinite
    |R| = uncountably infinite
    |R| > |N|   (Cantor diagonal argument)
    2^|N| = |R| (continuum hypothesis)
