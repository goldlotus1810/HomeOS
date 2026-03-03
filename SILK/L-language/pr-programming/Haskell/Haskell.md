# HS -- Haskell

symbol:    lambda (pure function)
born:      1990 (Simon Peyton Jones et al.)
paradigm:  purely functional, lazy, strongly typed
superpower: types that prove correctness at compile time

## Vi tri trong SILK

    Haskell ⊃ lambda-calculus  (theoretical foundation)
    Haskell → Rust             (influenced ownership/types)
    Haskell ≡ math             (code IS math notation)
    Haskell ⊃ infinity         (lazy = infinite data structures)

## UTF-32 → Haskell

    -- ∫ (fold)
    sum xs              = foldl (+) 0 xs
    integrate f xs      = foldl (\acc x -> acc + f x) 0 xs

    -- ∀ (map)
    doubles = map (*2) xs
    squares = map (^2) [1..]   -- infinite list!

    -- ∃ (any)
    hasEven = any even xs

    -- ∇ (derivative)
    deriv f x h = (f (x+h) - f (x-h)) / (2*h)

    -- ∪ (sum type / tagged union)
    data Shape
        = Sphere  { center :: Vec3, r :: Float }
        | Capsule { a :: Vec3, b :: Vec3, r :: Float }
        | Void

    -- ∞ (infinite lazy list)
    nats  = [1..]
    fibs  = 0 : 1 : zipWith (+) fibs (tail fibs)

    -- void/option (Maybe)
    safeDiv :: Int -> Int -> Maybe Int
    safeDiv _ 0 = Nothing
    safeDiv x y = Just (x `div` y)

    -- pure function (→)
    f :: Int -> Int
    f x = x * x + 1

## Core concepts → symbols

    pure function  -- maps input to output, no side effects (→)
    Maybe a        -- Nothing (void) or Just value          (◌|●)
    lazy list      -- potentially infinite                  (∞)
    type class     -- behavioral contract                   (≡)
    functor        -- map over container                    (∀→)
    monad          -- sequenced computation                 (→∘)
    fold           -- accumulate over structure             (∫)
    fix point      -- recursive cycle                       (♻)

## Euler identity in Haskell

    import Data.Complex
    euler :: Complex Double
    euler = exp (0 :+ pi) + 1   -- should be ~0

    Haskell: code = math, no translation loss
    SILK:    symbol = concept,  no translation loss
