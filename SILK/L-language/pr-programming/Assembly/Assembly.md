# ASM -- Assembly

symbol:    Assembly (direct metal)
paradigm:  imperative, register-based, no abstraction
level:     lowest human-readable (1 step above binary)
variants:  x86, ARM, MIPS, RISC-V, AVR

## Vi tri trong SILK

    Assembly ← machine code    (1:1 mapping, just readable)
    Assembly ← C               (C compiles to assembly)
    Assembly ≡ ∂∂∂             (direct partial access to everything)

## UTF-32 → x86 Assembly

    ; ∫ (accumulate / loop sum)
    ; EAX=sum, ECX=count, ESI=array ptr
        XOR  EAX, EAX
    LOOP:
        ADD  EAX, [ESI + ECX*4]
        INC  ECX
        CMP  ECX, N
        JL   LOOP

    ; ∀ vectorized (SIMD/AVX)
        VMOVAPS YMM0, [arr]
        VADDPS  YMM0, YMM0, YMM1   ; 8 floats at once

    ; ● (allocate on stack)
        SUB  RSP, 64               ; reserve space
        MOV  [RSP], rax

    ; → (function call)
        PUSH RDI
        CALL function_name
        POP  RDI

    ; 🛡 (bounds check)
        CMP  ECX, N
        JAE  out_of_bounds

    ; ∂ (pointer arithmetic)
        LEA  RAX, [RBX + RCX*8 + 16]

    ; ♻ (loop cycle)
        DEC  ECX
        JNZ  LOOP

## Registers → concepts

    RAX   accumulator result    (∫)
    RBX   base pointer          (∂)
    RCX   counter for loops     (∀ index)
    RDX   data / intermediate   (●)
    RSP   stack pointer         (∂ in call stack)
    RSI   source index          (∂ read)
    RDI   destination index     (∂ write)

## Bai hoc cho SILK

    Assembly proves:
      All computation = combination of:
        ∂  (load/store memory)
        ∫  (accumulate in register)
        ∀  (loop/branch)
        →  (call/jump)

      4 operations describe every program ever written.
      SILK: few root symbols + composition = everything.
