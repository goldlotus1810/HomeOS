# WASM -- WebAssembly

born:      2017 (W3C standard)
type:      portable binary format, near-native speed
runs:      browser, Node.js, server, edge, embedded

## Vi tri trong SILK

    WASM ← C/Rust/Go   (compile targets)
    WASM → browser      (runs in JS engine sandbox)
    WASM ≡ portable-asm (universal execution)
    WASM ⊃ WASI         (system interface for server)

## WASM text format (WAT)

    ;; sum array -- ∫ accumulate
    (module
      (func $sum (param $ptr i32) (param $len i32) (result f64)
        (local $i i32)
        (local $acc f64)
        (block $break
          (loop $loop
            (br_if $break
              (i32.ge_u (local.get $i) (local.get $len)))
            (local.set $acc
              (f64.add (local.get $acc)
                (f64.load
                  (i32.add (local.get $ptr)
                    (i32.mul (local.get $i) (i32.const 8))))))
            (local.set $i (i32.add (local.get $i) (i32.const 1)))
            (br $loop)))
        (local.get $acc))
      (export "sum" (func $sum)))

## Compile to WASM

    Rust → WASM:
        #[wasm_bindgen]
        pub fn sdf_sphere(px:f32,py:f32,pz:f32,
                          cx:f32,cy:f32,cz:f32,r:f32) -> f32 {
            ((px-cx).powi(2)+(py-cy).powi(2)+(pz-cz).powi(2))
                .sqrt() - r
        }
        // cargo build --target wasm32-unknown-unknown

    Go → WASM:
        // GOARCH=wasm GOOS=js go build -o main.wasm

    C → WASM:
        // emcc source.c -o output.wasm

## SILK + WASM

    WASM = ideal compile target for SILK:
      SILK node definition → compile → WASM
      → runs in browser (HomeOS!)
      → runs on server
      → runs on edge devices

    1 SILK program = runs everywhere (∀ platforms)
