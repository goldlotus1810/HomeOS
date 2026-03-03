# M -- MICROLIFE

symbol:    microbe
name_vi:   vi sinh vat / dot bien
branch:    B-biology

## Ngu nghia

Su song o cap do vi mo.
Dot bien ngau nhien + chon loc = tien hoa.
Don gian nhat nhung da dang nhat.

## Quan he

    microbe ← genetics    (defined by DNA)
    microbe → evolution   (mutation fuels evolution)
    microbe ∘ growth      (microbes grow)
    microbe ≡ mutation    (introduces variation)

## Dot bien → Computing

    Biological:
      DNA copy error → new variant
      Selection → better variants survive

    Computing (genetic algorithm):
      population = [genome1, genome2, ...]
      fitness(genome) → score
      crossover(g1, g2) → child
      mutate(child, rate) → variant

## Projections

Python:
    import random

    def mutate(genome, rate=0.01):
        return [
            gene if random.random() > rate
            else random.choice(alphabet)
            for gene in genome
        ]

    def crossover(g1, g2):
        cut = random.randint(0, len(g1))
        return g1[:cut] + g2[cut:]

Rust:
    fn mutate(v: f64, rate: f64, noise: f64) -> f64 {
        if rand::random::<f64>() < rate {
            v + (rand::random::<f64>() - 0.5) * noise
        } else { v }
    }

Go:
    func Mutate(v, rate float64) float64 {
        if rand.Float64() < rate {
            return v + (rand.Float64()-0.5)*0.2
        }
        return v
    }

## SILK connection

    SILK nodes can mutate:
      When a new interpretation is discovered
      for a symbol → node gains new relation.
    Mutation rate = rate of new research/insight.
