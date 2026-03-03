# ARM -- Mobile/Embedded Architecture

born:      1983 (Acorn RISC Machine)
type:      RISC -- Reduced Instruction Set Computing
dominant:  mobile, embedded, Apple Silicon
variants:  ARMv7 (32-bit), ARMv8/AArch64 (64-bit)

## RISC vs CISC

    CISC (x86):  complex instructions
                 MOV EAX, [EBX+ECX*4+8]  -- 1 instruction

    RISC (ARM):  simple uniform instructions
                 LDR R0, [R1, R2, LSL #2]
                 fewer transistors → less power → mobile wins

## Instructions → symbols

    ; ARM AArch64

    ; ∫ (sum array)
        MOV  X2, #0
    LOOP:
        LDR  W3, [X0], #4    ; load and advance pointer (∂)
        ADD  X2, X2, X3      ; accumulate (∫)
        SUBS X1, X1, #1
        BNE  LOOP            ; cycle (♻)

    ; ∀ (NEON SIMD -- ARM's parallel unit)
        LD1  {V0.4S}, [X0]
        FADD V0.4S, V0.4S, V1.4S   ; 4 floats at once

    ; → (function call)
        BL   function_name   ; branch + link
        RET

## Apple Silicon (M-series)

    ARM + Apple custom:
      Unified Memory:   CPU + GPU share same RAM
      Neural Engine:    dedicated neural hardware
      Performance cores + Efficiency cores = hybrid

    M3 Max: 16 perf + 4 eff = 20 cores → parallel (∀)
    Performance/watt = best in industry (2024)
