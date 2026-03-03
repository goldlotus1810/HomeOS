# PY -- Python

symbol:    Python
born:      1991 (Guido van Rossum)
paradigm:  multi-paradigm, readable, batteries-included
dominates: AI/ML, data science, scripting, web

## Vi tri trong SILK

    Python ⊃ numpy/scipy   (∫∇ in practice)
    Python ⊃ pytorch/tf    (neural training)
    Python ≈ JS            (scripting tier)
    Python ← C             (CPython written in C)

## UTF-32 → Python

    # ∫ (reduce / accumulate)
    total = sum(f(x) for x in xs)
    from functools import reduce
    product = reduce(lambda a,b: a*b, xs, 1)

    # ∀ (list comprehension)
    squares = [x**2 for x in xs]
    evens   = [x for x in xs if x % 2 == 0]

    # ∃ (any)
    has_prime = any(is_prime(x) for x in xs)

    # ∇ (gradient -- numpy)
    import numpy as np
    grad = np.gradient(field, axis=0)

    # ∪ (union types -- Python 3.10+)
    def f(x: int | str | None) -> str: ...

    # 🛡 (context manager = safe gate)
    with open('file.txt') as f:
        data = f.read()

    # ∀∇ (vectorized)
    result = np.dot(matrix, vector)

    # ♻ (cycle -- itertools)
    from itertools import cycle
    colors = cycle(['red', 'green', 'blue'])

## Core concepts → symbols

    list comprehension  -- for-all transform     (∀)
    generator           -- lazy cycle            (♻∂)
    decorator           -- function composition  (→∘)
    context manager     -- safe resource gate    (🛡)
    *args               -- collect all           (∀)
    None                -- void                  (◌)
    lambda              -- anonymous mapping     (→)
    yield               -- partial result        (∂t)
