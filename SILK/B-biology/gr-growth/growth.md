# G -- GROWTH

symbol:    seedling (growth)
name_vi:   tang truong / nay mam
branch:    B-biology

## Ngu nghia

Tich luy dan dan theo thoi gian.
Bat dau tu nho → phuc tap hon → truong thanh.

## Quan he

    growth ≡ ∫(rate·dt)   (= integral of rate over time)
    growth ← genetics      (program encoded in DNA)
    growth → ecosystem     (growth → fills environment)
    growth ∘ cycle = life  (growth + repetition = life)

## Mo hinh tang truong

    Linear:      f(t) = a·t
    Exponential: f(t) = A·e^(rt)           (bacteria, interest)
    Logistic:    f(t) = K/(1+e^(-r(t-t0))) (with carrying cap)
    Sigmoidal:   f(t) = 1/(1+e^(-t))       (neural activation)
    SILK/HomeOS: f(t) = min(1, age+rate·dt) (bounded growth)

## Projections

C:
    // HomeOS tree growth
    age = fminf(1.0f, age + rate * dt);

    // logistic
    float logistic(float t, float K, float r, float t0) {
        return K / (1.0f + expf(-r * (t - t0)));
    }

Rust:
    fn grow(age: f64, rate: f64, dt: f64) -> f64 {
        (age + rate * dt).min(1.0)
    }
    fn logistic(t: f64, k: f64, r: f64, t0: f64) -> f64 {
        k / (1.0 + (-r*(t-t0)).exp())
    }

Go:
    func Grow(age, rate, dt float64) float64 {
        return math.Min(1.0, age+rate*dt)
    }

Python:
    grow = lambda age, rate, dt: min(1.0, age + rate*dt)
    import numpy as np
    t = np.linspace(0, 10, 100)
    pop = 100 * np.exp(0.3 * t)   # exponential

## Instances

    Biology:  cell division, organ dev, aging
    Plants:   photosynthesis → biomass
    Code:     codebase grows with features
    SILK:     tree grows with new nodes each session
    HomeOS:   trees grow each animation frame
    AI:       model improves with training epochs
