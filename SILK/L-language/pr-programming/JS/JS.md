# JS -- JavaScript

symbol:    JS
born:      1995 (Brendan Eich, Netscape)
paradigm:  multi-paradigm, event-driven, async
runs:      browser + Node.js + Deno + edge

## Vi tri trong SILK

    JS = language of the web
    JS ⊃ async/await  (non-blocking cycle)
    JS → WASM         (compile target)
    JS ≈ Python       (scripting, dynamic typing)

## UTF-32 → JS

    // ∫ (reduce / accumulate)
    const sum = xs.reduce((s, x) => s + f(x), 0);

    // ∀ (map / forEach)
    const doubled = xs.map(x => x * 2);
    xs.forEach(x => process(x));

    // ∃ (some = there exists)
    const hasEven = xs.some(x => x % 2 === 0);

    // ♻ (cycle / event loop)
    setInterval(() => tick(world, dt), 16);   // 60fps

    // parallel (Promise.all)
    const [a, b] = await Promise.all([fa(), fb()]);

    // ∪ (union type -- TypeScript)
    type Shape = Sphere | Capsule | Box;

    // void
    null | undefined | void 0

    // 🛡 (try/catch)
    try {
        const data = JSON.parse(input);
    } catch(e) {
        console.error('blocked:', e);
    }

    // world object (HomeOS)
    const W = {
        seed: 888, time: 10.0, genes: [],
        animate(dt) { this.time += dt; }
    };

## Core concepts → symbols

    closure       -- captures partial environment (∂)
    Promise       -- void becoming value, async  (◌→●)
    async/await   -- non-blocking cycle           (♻)
    event loop    -- eternal cycle                (♻)
    null          -- void                         (◌)
    spread ...    -- apply to all                 (∀)
    lambda =>     -- anonymous function           (→)
