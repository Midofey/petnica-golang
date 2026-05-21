# gonum/mat — Kratki priručnik za radionicu 2

Ovaj fajl objašnjava sve `gonum/mat` funkcije koje vam mogu biti potrebne da završite `02_mnist_matrices`.

---

## Kreiranje matrica

### `mat.NewDense(redovi, kolone, podaci)`

Pravi novu matricu zadate veličine.

```go
// Matrica 2x3 sa zadatim vrednostima (red po red)
m := mat.NewDense(2, 3, []float64{
    1, 2, 3,
    4, 5, 6,
})

// Matrica 1x128 popunjena nulama (nil = sve nule)
z := mat.NewDense(1, 128, nil)
```

### `mat.DenseCopyOf(m)`

Pravi nezavisnu kopiju matrice. Promena kopije ne menja original.

```go
kopija := mat.DenseCopyOf(original)
```

---

## Čitanje vrednosti

### `m.At(red, kolona)`

Čita jednu vrednost iz matrice.

```go
v := m.At(0, 5) // Element u redu 0, koloni 5
```

### `m.Dims()`

Vraća dimenzije matrice kao `(redovi, kolone)`.

```go
r, c := m.Dims() // npr. r=1, c=784
```

---

## Matricno množenje (`.Mul`)

### `rezultat.Mul(A, B)`

Standardno matrično množenje: `rezultat = A × B`

**Unutrašnje dimenzije moraju da se poklapaju!**

```go
// A je [1x784], B je [784x128] → rezultat je [1x128]
rezultat := mat.NewDense(1, 128, nil)
rezultat.Mul(A, B)
```

Dimenzije: `[m×n] · [n×p] = [m×p]`

---

## Sabiranje i oduzimanje (`.Add`, `.Sub`)

### `rezultat.Add(A, B)` / `rezultat.Sub(A, B)`

Element-po-element sabiranje/oduzimanje. Obe matrice **moraju biti istog oblika**.

```go
// A i B su obe [1x10]
rezultat := mat.NewDense(1, 10, nil)
rezultat.Add(A, B)  // rezultat[i][j] = A[i][j] + B[i][j]
rezultat.Sub(A, B)  // rezultat[i][j] = A[i][j] - B[i][j]
```

> **Napomena:** Rezultat može biti ista matrica kao ulaz:
> ```go
> Z.Add(Z, B)  // Z = Z + B (in-place)
> ```

---

## Element-po-element množenje (`.MulElem`)

### `rezultat.MulElem(A, B)`

Množi **svaki element** matrice A sa odgovarajućim elementom matrice B. Obe matrice moraju biti **istog oblika**.

```go
// A i B su obe [1x10]
rezultat := mat.NewDense(1, 10, nil)
rezultat.MulElem(A, B)  // rezultat[i][j] = A[i][j] * B[i][j]
```

### Kad koristiti `.Mul()` vs `.MulElem()`?

| Operacija | Metod | Primer |
|---|---|---|
| Matrično množenje (različiti oblici) | `.Mul()` | `[1×784] · [784×128] = [1×128]` |
| Element-po-element (isti oblik) | `.MulElem()` | `[1×10] * [1×10] = [1×10]` |

---

## Primena funkcije na svaki element (`.Apply`)

### `rezultat.Apply(f, izvor)`

Primenjuje funkciju `f` na **svaki element** matrice `izvor` i čuva rezultat.

Funkcija **mora** imati tačno ovaj potpis:
```go
func(i, j int, v float64) float64
```

- `i, j` — red i kolona (često nepotrebni, ali Go zahteva da se navedu)
- `v` — trenutna vrednost elementa
- vraća — novu vrednost tog elementa

```go
// Primeni sigmoid na svaki element matrice Z
A := mat.NewDense(1, 128, nil)
A.Apply(func(i, j int, v float64) float64 {
    return 1.0 / (1.0 + math.Exp(-v))
}, Z)
```

> **Napomena:** Može se primeniti i in-place:
> ```go
> diff.Apply(func(i, j int, v float64) float64 {
>     return v * v  // Kvadrira svaki element
> }, diff)
> ```

---

## Transponovanje (`.T()`)

### `m.T()`

Zamenjuje redove i kolone matrice. **Ne kopira podatke** — vraća "pogled" (view).

```go
// A1 je [1x128] → A1.T() je [128x1]
dW.Mul(A1.T(), delta)  // [128x1] · [1x10] = [128x10]
```

Npr. koristimo transponovanje da se **dimenzije poklope** za matrično množenje u backward pass-u.

---

## Skaliranje (`.Scale`)

### `rezultat.Scale(skalar, M)`

Množi **svaki element** matrice skalarom (običnim brojem).

```go
scaled := mat.NewDense(128, 10, nil)
scaled.Scale(0.1, dW)  // scaled[i][j] = 0.1 * dW[i][j]
```

> **Napomena:** Može i in-place:
> ```go
> dW.Scale(learningRate, dW)  // dW = learningRate * dW
> ```

---

## Sumiranje svih elemenata (`mat.Sum`)

### `mat.Sum(m)`

Vraća sumu **svih elemenata** matrice kao `float64`.

```go
total := mat.Sum(diff)  // npr. za [1x10] matricu, sumira svih 10 vrednosti
```

---

## Pregled svih metoda za radionicu 2

| Metod | Tip operacije | Oblik | Koristi se u |
|---|---|---|---|
| `mat.NewDense(r, c, data)` | Kreiranje | `[r × c]` | Svuda |
| `.Mul(A, B)` | Matrično množenje | `[m×n]·[n×p]=[m×p]` | Forward, Backward |
| `.Add(A, B)` | Element-wise + | Isti oblici | Forward (bias) |
| `.Sub(A, B)` | Element-wise − | Isti oblici | Backward, Loss, Update |
| `.MulElem(A, B)` | Element-wise × | Isti oblici | Backward (delta) |
| `.Apply(f, src)` | Po elementu | Isti oblik | Sigmoid, Loss |
| `.T()` | Transponovanje | `[m×n]→[n×m]` | Backward |
| `.Scale(s, M)` | Skalar × matrica | Isti oblik | Update |
| `mat.Sum(M)` | Suma svih elem. | → `float64` | Loss |
| `.At(r, c)` | Čitanje jedne vr. | → `float64` | Predikcija |
| `mat.DenseCopyOf(M)` | Kopija | Isti oblik | Backward (dB) |
