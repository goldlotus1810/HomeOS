# F -- FIELDS

symbol:    F (scalar/vector field)
name_vi:   truong vat ly
branch:    R-reality

## Ngu nghia

Anh huong lan rong qua khong gian.
Moi diem trong khong gian co mot gia tri (scalar or vector).

## Quan he

    field ≡ ∇·potential    (field = gradient of potential)
    field ⊃ energy         (fields carry energy)
    field ≈ SDF            (Signed Distance Field = scalar field!)
    field → ∂              (fields vary with space/time)

## Cac loai truong

    Scalar: f(x,y,z) → number       (temperature, pressure, SDF)
    Vector: f(x,y,z) → vec3         (wind, gravity, EM)
    Tensor: f(x,y,z) → matrix       (stress, spacetime curvature)
    SDF:    f(x,y,z) → distance      (HomeOS terrain!)

## Projections

GLSL (SDF terrain field):
    float terrainSDF(vec3 p) {
        return p.y - fbm(p.x, p.z);
    }

Python (scalar field on grid):
    import numpy as np
    field = np.array([[temp(i*0.1, j*0.1)
                       for j in range(100)]
                      for i in range(100)])

Go (vector field):
    type Field func(x, y, z float64) [3]float64
    var gravity Field = func(x, y, z float64) [3]float64 {
        return [3]float64{0, -9.8, 0}
    }

HomeOS connection:
    H(x,z)   = fbm(x,z)          -- height field
    normal   = gradient(H, x, z) -- ∇ of field
    SDF(p)   = p.y - H(p.x,p.z)  -- signed distance
