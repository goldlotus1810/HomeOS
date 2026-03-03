# A -- ALGEBRA

symbol:    ∑ (summation / algebra)
name_vi:   dai so / tong
branch:    M-math

## Ngu nghia

Quan he giua cac phan tu theo quy tac nhat dinh.
Nen tang cua moi phep tinh, moi cau truc so.

## Quan he

    ∑ ⊂ ∫      (∑ la truong hop roi rac cua ∫)
    ∑ ⊃ +      (cong la ∑ cua 2 phan tu)
    ∑ ≡ fold   (fold trong FP)
    ∑ ≡ reduce (reduce trong JS/Python)

## Cau truc dai so

    Group  (Nhom):   (G, .) -- closure, assoc, identity, inverse
    Ring   (Vanh):   (R, +, .) -- group + distributive
    Field  (Truong): (F, +, .) -- ring + mult inverse
    Vector space:    scalars x vectors -- linear algebra

    Computing:
      integers = ring   (no division)
      floats   = field  (with rounding errors)
      matrices = ring   (non-commutative multiply)
      booleans = field mod 2

## Projections

C:
    int sum = 0;
    for(int i = 0; i < n; i++) sum += arr[i];

    // matrix multiply: ∑ a[i][j] * b[j][k]
    for(int i=0;i<N;i++)
      for(int j=0;j<N;j++)
        for(int k=0;k<N;k++)
          C[i][k] += A[i][j] * B[j][k];

Rust:
    let sum: i32 = arr.iter().sum();
    let dot: f64 = a.iter().zip(&b).map(|(x,y)| x*y).sum();

Go:
    sum := 0
    for _, v := range arr { sum += v }

Python:
    total = sum(arr)
    dot = sum(a*b for a,b in zip(va, vb))

Haskell:
    sum xs
    foldl (+) 0 xs
