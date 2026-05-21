// ============================================================================
// Playground 2: MNIST matricna neuronska mreza
// ============================================================================
// Prelazimo sa skalara na matrice! Umesto jednog neurona, sada imamo
// citave SLOJEVE neurona koji rade paralelno koristeci matricno mnozenje.
//
// Struktura mreze:
//   [784 piksela] --[W1: 784x128]--> [sigmoid] --[W2: 128x10]--> [sigmoid] --> [10 klasa]
//
// Cilj: prepoznavanje rukom pisanih cifara (MNIST dataset)
// ============================================================================

package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"gonum.org/v1/gonum/mat"

	"petnica-dl-workshop/helpers"
)

// ==========================================================
// TODO 1a: Napisi sigmoidMatrix funkciju koju ce koristiti mat.Apply.
// mat.Apply primenjuje funkciju na SVAKI element matrice.
//
// mat.Apply zahteva TACNO ovaj potpis funkcije:
//
//	func(i, j int, v float64) float64
//
// Go nas tera da navedemo i, j (red i kolona elementa) cak i
// ako ih ne koristimo. Mi samo treba da transformisemo v.
//
// Hint: Pogledajte kako smo uradili u 1. delu (01_scalar_network/main.go)
// ==========================================================
func sigmoidMatrix(i, j int, v float64) float64 {
	return 1 / (1 + math.Exp(-v))
}

// ==========================================================
// TODO 1b: Napisi izvod sigmoid funkcije za matrice.
// Isti potpis kao TODO 1a — samo drugacija formula.
//
// Hint: Pogledajte kako smo uradili u 1. delu (01_scalar_network/main.go)
// ==========================================================
func sigmoidDerivMatrix(i, j int, v float64) float64 {
	return sigmoidMatrix(i, j, v) * (1 - sigmoidMatrix(i, j, v))
}

// ==========================================================
// TODO 2: Implementiraj forward pass sa matricama.
// Ovo je matricna verzija forward pass-a iz prvog dela.
//
// Korak A - Skriveni sloj:
//
//	Z1 = X @ W1 + B1       (matricno mnozenje + bias)
//	A1 = sigmoid(Z1)        (primeniti na svaki element)
//
// Korak B - Izlazni sloj:
//
//	Z2 = A1 @ W2 + B2
//	A2 = sigmoid(Z2)        (finalna predikcija)
//
// Hint: koristite .Mul(), .Add(), .Apply()
// ==========================================================
func forward(xMat *mat.Dense, W1, B1, W2, B2 *mat.Dense) (Z1, A1, Z2, A2 *mat.Dense) {
	_, hiddenSize := W1.Dims()
	_, outputSize := W2.Dims()

	Z1 = mat.NewDense(1, hiddenSize, nil)
	Z1.Mul(xMat, W1)
	Z1.Add(Z1, B1)

	A1 = mat.NewDense(1, hiddenSize, nil)
	A1.Apply(sigmoidMatrix, Z1)

	Z2 = mat.NewDense(1, outputSize, nil)
	Z2.Mul(A1, W2)
	Z2.Add(Z2, B2)

	A2 = mat.NewDense(1, outputSize, nil)
	A2.Apply(sigmoidMatrix, Z2)

	return
}

func main() {
	fmt.Println("============================================")
	fmt.Println("  Playground 2: MNIST matricna mreza")
	fmt.Println("============================================")
	fmt.Println()

	// ------------------------------------------------------------------
	// 1. Ucitavanje MNIST podataka
	// ------------------------------------------------------------------
	fmt.Println("[1] Ucitavanje MNIST podataka...")

	// Ocekujemo CSV fajl: prva kolona = labela, sledecih 784 = pikseli
	dataPath := "data/mnist_train.csv"
	if len(os.Args) > 1 {
		dataPath = os.Args[1]
	}

	inputs, labels, err := helpers.LoadMNIST(dataPath)
	if err != nil {
		fmt.Printf("    GRESKA: %v\n", err)
		fmt.Println("    Proverite da li postoji fajl:", dataPath)
		fmt.Println("    Format: label,pixel0,pixel1,...,pixel783")
		return
	}

	numSamples, inputSize := inputs.Dims()
	_, outputSize := labels.Dims()
	fmt.Printf("    Uzoraka: %d, Ulaz: %d piksela, Izlaz: %d klasa\n",
		numSamples, inputSize, outputSize)
	fmt.Println()

	// ------------------------------------------------------------------
	// 2. Inicijalizacija tezina (nasumicne matrice)
	// ------------------------------------------------------------------
	fmt.Println("[2] Inicijalizacija tezina...")

	hiddenSize := 128 // Broj neurona u skrivenom sloju

	// Tezine i biasi prvog sloja: [784 x 128]
	w1Data := make([]float64, inputSize*hiddenSize)
	for i := range w1Data {
		w1Data[i] = (rand.Float64() - 0.5) * 0.1 // Mali nasumicni brojevi
	}
	// mat.NewDense(redovi, kolone, podaci) pravi novu matricu.
	//   - redovi: broj redova matrice
	//   - kolone: broj kolona matrice
	//   - podaci: slice float64 vrednosti koji popunjava matricu red po red,
	//             ili nil ako zelimo matricu popunjenu nulama.
	W1 := mat.NewDense(inputSize, hiddenSize, w1Data)
	B1 := mat.NewDense(1, hiddenSize, nil) // nil = sve nule

	// Tezine i biasi drugog sloja: [128 x 10]
	w2Data := make([]float64, hiddenSize*outputSize)
	for i := range w2Data {
		w2Data[i] = (rand.Float64() - 0.5) * 0.1
	}
	W2 := mat.NewDense(hiddenSize, outputSize, w2Data)
	B2 := mat.NewDense(1, outputSize, nil)

	fmt.Printf("    W1: %dx%d, B1: 1x%d\n", inputSize, hiddenSize, hiddenSize)
	fmt.Printf("    W2: %dx%d, B2: 1x%d\n", hiddenSize, outputSize, outputSize)
	fmt.Println()

	// Hiperparametri
	learningRate := 0.1
	epochs := 5

	lossHistory := make([]float64, 0, epochs)

	// ------------------------------------------------------------------
	// 3. Petlja treniranja
	// ------------------------------------------------------------------
	fmt.Println("[3] Pocinje treniranje...")
	fmt.Println()

	for epoch := 0; epoch < epochs; epoch++ {
		// Za svaki uzorak pojedinacno (stohasticki gradijentni spust)
		totalLoss := 0.0

		for s := 0; s < numSamples; s++ {
			// Izvlacimo s-ti uzorak iz dataseta i pretvaramo ga u matricu [1 x N].
			// Ovo je potrebno jer mat.Mul radi samo sa matricama, ne sa vektorima.
			xRow := inputs.RowView(s)
			xMat := mat.NewDense(1, inputSize, nil)
			for j := 0; j < inputSize; j++ {
				// xRow.AtVec(j) cita j-ti element vektora.
				// xMat.Set(red, kolona, vrednost) upisuje vrednost u matricu.
				xMat.Set(0, j, xRow.AtVec(j))
			}

			// Ciljna labela kao matrica [1 x 10] -- "verovatnoca" da je cifra upravo ta
			yRow := labels.RowView(s)
			yMat := mat.NewDense(1, outputSize, nil)
			for j := 0; j < outputSize; j++ {
				yMat.Set(0, j, yRow.AtVec(j))
			}

			// Forward pass - koristimo funkciju iz TODO 2
			Z1, A1, Z2, A2 := forward(xMat, W1, B1, W2, B2)

			// ==========================================================
			// TODO 3: Izracunaj loss za ovaj uzorak.
			// Ovaj put moramo da saberemo gubitke za svaki neuron u izlaznom sloju.
			//
			// Formula: loss = (1/n) * sum((yTrue - yPred)^2)
			//  -- n je broj neurona u izlaznom sloju (outputSize)
			//  -- yTrue je prava vrednost (yMat)
			//  -- yPred je predikcija (A2)
			//
			// Hint:
			//   Pocnite sa definisanjem potrebne matrice koja se koristi za medjuvrednost (razliku)
			//   diff := mat.NewDense(1, outputSize, nil)
			//
			//   ...pokusajte ostatak da resite sami koristeci formulu i hint iznad.
			//
			// ==========================================================
			loss := 0.0 // Implementirajte prema uputstvu iznad
			totalLoss += loss

			// ==========================================================
			// TODO 4: Implementiraj backward pass sa matricama.
			// Ovo je ista logika kao u prvoj vezbi, ali sa matricama!
			//
			// Setite se backward pass-a iz prve vezbe. Ista matematika,
			// samo sa matricama umesto skalara:
			//
			// LEGENDA:  @ = matricno mnozenje (.Mul)
			//           * = element-po-element (.MulElem)
			//
			//   outputError = A2 - Y                              [1x10]
			//
			//   delta2      = outputError * d_sigmoid(Z2)         [1x10]
			//   dW2         = A1^T @ delta2           [128x1] @ [1x10] = [128x10]
			//   dB2         = delta2                                      [1x10]
			//
			//   hiddenError = delta2 @ W2^T            [1x10] @ [10x128] = [1x128]
			//   delta1      = hiddenError * d_sigmoid(Z1)                  [1x128]
			//   dW1         = X^T @ delta1            [784x1] @ [1x128] = [784x128]
			//   dB1         = delta1                                       [1x128]
			//
			// NAPOMENA: ^T predstavlja transponovanje matrice (zamena redova i kolona).
			// Ovo je potrebno da bi se dimenzije matrica poklopile za mnozenje.
			//
			// Hint: koristite .Sub(), .MulElem(), .Mul(), .Apply(), .T()
			//   Pocnite sa outputError := mat.NewDense(1, outputSize, nil)
			// ==========================================================

			// ==========================================================
			// TODO 5: Azuriraj tezine (gradijentni spust).
			// Setite se formule iz prve vezbe, samo sada sa matricama:
			//
			//   W2 = W2 - learningRate * dW2
			//   B2 = B2 - learningRate * dB2
			//   W1 = W1 - learningRate * dW1
			//   B1 = B1 - learningRate * dB1
			//
			// Hint: koristite .Scale() i .Sub()
			// ==========================================================

			// Sprecavamo upozorenja kompajlera (obrisite kada zavrsiste TODO-ove)
			_ = Z1
			_ = A1
			_ = Z2
			_ = A2
			_ = yMat
		}

		avgLoss := totalLoss / float64(numSamples)
		lossHistory = append(lossHistory, avgLoss)

		fmt.Printf("    Epoha %2d/%d | Gubitak: %.6f\n", epoch+1, epochs, avgLoss)
	}

	fmt.Println()
	fmt.Println("    Treniranje zavrseno!")
	fmt.Println()

	// ------------------------------------------------------------------
	// 4. Vizuelizacija rezultata
	// ------------------------------------------------------------------
	fmt.Println("[4] Vizuelizacija rezultata...")

	// Cuvamo grafikon gubitka
	if err := helpers.PlotLoss(lossHistory, "mnist_loss.png"); err != nil {
		fmt.Printf("    GRESKA: %v\n", err)
	}

	// Cuvamo prvih 10 predikcija kao slike
	fmt.Println("    Generisanje predikcija za prvih 10 uzoraka...")
	numPreview := 10
	if numSamples < numPreview {
		numPreview = numSamples
	}

	if err := os.MkdirAll("predictions", 0755); err != nil {
		fmt.Printf("    GRESKA pri kreiranju foldera: %v\n", err)
	}

	for i := 0; i < numPreview; i++ {
		// Izvlacimo piksele uzorka
		pixelRow := inputs.RowView(i)
		pixels := make([]float64, 784)
		for j := 0; j < 784; j++ {
			pixels[j] = pixelRow.AtVec(j)
		}

		// Koristimo ISTU forward() funkciju iz TODO 2
		xMat := mat.NewDense(1, inputSize, pixels)
		_, _, _, A2 := forward(xMat, W1, B1, W2, B2)

		// Pronalazimo klasu sa najvecom verovatnocom
		prediction := 0
		maxProb := A2.At(0, 0)
		for j := 1; j < outputSize; j++ {
			if A2.At(0, j) > maxProb {
				maxProb = A2.At(0, j)
				prediction = j
			}
		}

		// Prava labela
		trueLabel := 0
		labelRow := labels.RowView(i)
		for j := 0; j < outputSize; j++ {
			if labelRow.AtVec(j) == 1.0 {
				trueLabel = j
				break
			}
		}

		path := fmt.Sprintf("predictions/uzorak_%02d_tacno_%d_pred_%d.png", i, trueLabel, prediction)
		if err := helpers.SaveDigitImage(pixels, prediction, path); err != nil {
			fmt.Printf("    GRESKA: %v\n", err)
		}
	}

	fmt.Println()
	fmt.Println("============================================")
	fmt.Println("  Gotovo! Proverite mnist_loss.png i")
	fmt.Println("  predictions/ folder za rezultate.")
	fmt.Println("============================================")

	// ------------------------------------------------------------------
	// 5. Interaktivni mod — crtanje cifara u browseru
	// ------------------------------------------------------------------
	// Pokrecemo web server gde mozete crtati cifre i videti predikcije uzivo!
	predict := func(pixels []float64) (int, []float64) {
		xMat := mat.NewDense(1, inputSize, pixels)
		_, _, _, A2 := forward(xMat, W1, B1, W2, B2)

		// Pronalazimo klasu sa najvecom verovatnocom
		prediction := 0
		maxProb := A2.At(0, 0)
		probs := make([]float64, outputSize)
		probs[0] = A2.At(0, 0)
		for j := 1; j < outputSize; j++ {
			probs[j] = A2.At(0, j)
			if probs[j] > maxProb {
				maxProb = probs[j]
				prediction = j
			}
		}
		return prediction, probs
	}

	if err := helpers.StartDrawingServer(8080, predict); err != nil {
		fmt.Printf("GRESKA: %v\n", err)
	}

	// Sprecavamo upozorenja kompajlera
	_ = sigmoidDerivMatrix
	_ = learningRate
	_ = math.Exp
}
