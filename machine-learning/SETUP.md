# Uputstvo za pripremu okruženja

Kratki vodič za postavljanje svega što vam je potrebno za radionicu.

---

## 1. Instalacija Go-a

Potreban je **Go 1.22+**.

### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install golang-go
```

### macOS
```bash
brew install go
```

### Windows
Preuzmite installer sa [go.dev/dl](https://go.dev/dl/) i pratite korake.

### Provera instalacije
```bash
go version
# Treba da ispiše: go version go1.22.x ...
```

## 2. Preuzimanje dataseta

U folderu `data/` su vam potrebna 2 fajla:

| Fajl | Veličina | Izvor |
|------|----------|-------|
| `student_performance_finalscore.csv` | ~700 KB | [Kaggle](https://www.kaggle.com/datasets/sarveshchhetri/student-lifestyle-vs-academic-performance-dataset) |
| `mnist_train.csv` | ~110 MB | [Kaggle](https://www.kaggle.com/datasets/oddrationale/mnist-in-csv) |

Stavite ih u `data/` folder:

**PROVERITE PRVO DA LI SU VEC TU**
```
ml-workshop/
└── data/
    ├── student_performance_finalscore.csv
    └── mnist_train.csv
```

> **Napomena:** `mnist_train.csv` je veliki fajl (~110 MB). Ako imate spor internet, dostavićemo ga na USB-u.

---

## 3. Instalacija zavisnosti

Iz korenog foldera projekta:
```bash
cd ml-workshop
go mod download
```

Ovo preuzima `gonum` (biblioteka za matrice) i `gonum/plot` (za grafikone).

### Provera

```bash
go build ./...
```

Ako nema grešaka, sve je spremno!

---

## 4. Pokretanje vežbi

Svaka vežba se pokreće pojedinačno:

```bash
# Vežba 1: Skalarna mreža
go run 01_scalar_network/main.go

# Vežba 2: MNIST sa matricama
go run 02_mnist_matrices/main.go

# Vežba 3: Autodiff engine
go run 03_autodiff/engine.go 03_autodiff/main.go
```

---

## Česte greške

### `cannot find package "gonum.org/v1/gonum/mat"`
```bash
go mod download
```

### `open data/mnist_train.csv: no such file or directory`
Preuzmite MNIST dataset i stavite ga u `data/` folder.

### `address already in use` (port 8080)
Zatvorite prethodni proces (Ctrl+C) ili će server automatski izabrati slobodan port.
