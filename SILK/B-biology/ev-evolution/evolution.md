# E -- EVOLUTION

symbol:    recycling (cycle + selection)
name_vi:   tien hoa / chu ky
branch:    B-biology

## Ngu nghia

Chu ky chon loc: bien the → canh tranh → ke thich nghi ton tai.
Khong co dich den — chi co thich nghi voi moi truong hien tai.

## 3 Dieu kien

    1. Variation:  ca the khac nhau  (mutation)
    2. Heredity:   con giong cha me   (DNA)
    3. Selection:  ke thich nghi song (fitness)

    → Du 3 dieu kien = tien hoa xay ra
    → Khong can ke hoach, khong can tri tue

## Quan he

    evolution ≡ cycle(t, period)
    evolution ⊃ mutation       (fuel)
    evolution ⊃ growth         (result)
    evolution → ecosystem
    evolution ≡ iteration      (loop + selection in computing)

## Projections

Python (evolutionary algorithm):
    def evolve(population, fitness_fn, generations=100):
        for _ in range(generations):
            scored = sorted(population,
                            key=fitness_fn, reverse=True)
            survivors = scored[:len(scored)//2]
            children = [crossover(survivors[i], survivors[i+1])
                        for i in range(0, len(survivors)-1, 2)]
            population = survivors + [mutate(c) for c in children]
        return max(population, key=fitness_fn)

Go (oscillation / day-night cycle):
    func Cycle(t, period float64) float64 {
        return (math.Sin(2*math.Pi*t/period) + 1) * 0.5
    }
    func DayNight(hour float64) float64 {
        return Cycle(hour, 24.0)
    }

C:
    float cycle(float t, float period) {
        return (sinf(2*3.14159f*t/period) + 1.0f) * 0.5f;
    }
    t = fmodf(t + dt, period);   // wrap around

## Instances

    Biology:  species adapt over millions of years
    Software: agile iterations, A/B testing, CI/CD
    AI:       reinforcement learning
    SILK:     tree structure evolves each session
    HomeOS:   day/night cycle(t, 24) drives lighting
