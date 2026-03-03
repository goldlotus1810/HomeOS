# N -- NEURAL

symbol:    brain (neural network)
name_vi:   mang than kinh / nao
branch:    B-biology

## Ngu nghia

Mang luoi cac node xu ly thong tin, hoc tu kinh nghiem.
Tin hieu lan truyen, khuech dai hoac uc che.

## Quan he

    neural ≡ gradient·weights  (hoc = gradient descent)
    neural ⊃ integral          (neuron = ∫ inputs·weights)
    neural → cognitive         (brain → thought)
    neural ← genetics          (brain structure in DNA)
    neural ≈ computing         (biological ≈ artificial NN)

## Neuron → math

    y = σ( ∑ wᵢxᵢ + b )
    y = σ( ∫ weights·inputs )

    σ (activation): tanh, ReLU, sigmoid
    ∑ (sum):        weighted inputs
    ∇ (gradient):   backpropagation signal
    W (weights):    learned parameters

## Projections

Python:
    import numpy as np

    def neuron(x, w, b):
        return np.tanh(np.dot(w, x) + b)

    def train_step(w, x, y_true, lr=0.01):
        y = np.tanh(np.dot(w, x))
        grad = 2*(y - y_true) * (1 - y**2)   # ∇ loss
        return w - lr * grad * x              # update

Rust:
    fn neuron(x: &[f64], w: &[f64], b: f64) -> f64 {
        let z: f64 = w.iter().zip(x)
                      .map(|(wi, xi)| wi * xi)
                      .sum::<f64>() + b;
        z.tanh()
    }

Go:
    type Neuron struct{ W []float64; B float64 }

    func (n *Neuron) Forward(x []float64) float64 {
        sum := n.B
        for i, w := range n.W { sum += w * x[i] }
        return math.Tanh(sum)
    }

## SILK connection

    SILK tree ≈ neural network:
      Nodes    = neurons
      Edges    = weighted connections
      Learning = adding / updating nodes
      Query    = forward pass through tree
    Difference:
      Neural: weights = numbers
      SILK:   weights = typed relations (∈, ⊃, ≡, →)
