// ============================================================================
// Playground 3: MNIST sa Autodiff Engine-om
// ============================================================================
// Sada koristimo nas Value engine da AUTOMATSKI izracuna gradijente!
// Nema vise rucnog racunanja izvoda - samo pozovemo loss.Backward().
//
// Ovo je "Full Circle" momenat: isti MNIST zadatak kao Playground 2,
// ali sa automatskom diferencijacijom umesto rucne.
// ============================================================================

package main

import (
	"fmt"
	"math/rand"
	"os"

	"petnica-dl-workshop/helpers"
)

func main() {
	fmt.Println("============================================")
	fmt.Println("  Playground 3: MNIST sa Autodiff Engine-om")
	fmt.Println("============================================")
	fmt.Println()

	// ------------------------------------------------------------------
	// 1. Ucitavanje MNIST podataka
	// ------------------------------------------------------------------
	fmt.Println("[1] Ucitavanje MNIST podataka...")

	dataPath := "data/mnist_train.csv"
	if len(os.Args) > 1 {
		dataPath = os.Args[1]
	}

	inputs, labels, err := helpers.LoadMNIST(dataPath)
	if err != nil {
		fmt.Printf("    GRESKA: %v\n", err)
		fmt.Println("    Proverite da li postoji fajl:", dataPath)
		return
	}

	numSamples, inputSize := inputs.Dims()
	_, outputSize := labels.Dims()
	fmt.Printf("    Uzoraka: %d, Ulaz: %d, Izlaz: %d\n", numSamples, inputSize, outputSize)

	// Koristimo manji podskup za brzinu
	maxSamples := 500
	if numSamples > maxSamples {
		numSamples = maxSamples
		fmt.Printf("    Koristimo prvih %d uzoraka za brzinu\n", maxSamples)
	}
	fmt.Println()

	// ------------------------------------------------------------------
	// 2. Inicijalizacija tezina kao Value struktura
	// ------------------------------------------------------------------
	fmt.Println("[2] Inicijalizacija tezina kao Value struktura...")

	hiddenSize := 32 // Manji skriveni sloj za brzinu autodiff-a

	// Tezine prvog sloja: inputSize x hiddenSize
	W1 := make([][]*Value, inputSize)
	for i := range W1 {
		W1[i] = make([]*Value, hiddenSize)
		for j := range W1[i] {
			W1[i][j] = NewValue((rand.Float64()-0.5)*0.1, fmt.Sprintf("w1_%d_%d", i, j))
		}
	}
	B1 := make([]*Value, hiddenSize)
	for j := range B1 {
		B1[j] = NewValue(0.0, fmt.Sprintf("b1_%d", j))
	}

	// Tezine drugog sloja: hiddenSize x outputSize
	W2 := make([][]*Value, hiddenSize)
	for i := range W2 {
		W2[i] = make([]*Value, outputSize)
		for j := range W2[i] {
			W2[i][j] = NewValue((rand.Float64()-0.5)*0.1, fmt.Sprintf("w2_%d_%d", i, j))
		}
	}
	B2 := make([]*Value, outputSize)
	for j := range B2 {
		B2[j] = NewValue(0.0, fmt.Sprintf("b2_%d", j))
	}

	fmt.Printf("    W1: %dx%d Value-ova, B1: %d\n", inputSize, hiddenSize, hiddenSize)
	fmt.Printf("    W2: %dx%d Value-ova, B2: %d\n", hiddenSize, outputSize, outputSize)
	fmt.Println()

	// Hiperparametri
	learningRate := 0.1
	epochs := 5

	lossHistory := make([]float64, 0, epochs)

	// Prikupljamo SVE tezine u jedan slice za lakse azuriranje
	allWeights := make([]*Value, 0)
	for i := range W1 {
		allWeights = append(allWeights, W1[i]...)
	}
	allWeights = append(allWeights, B1...)
	for i := range W2 {
		allWeights = append(allWeights, W2[i]...)
	}
	allWeights = append(allWeights, B2...)

	fmt.Printf("    Ukupan broj parametara: %d\n", len(allWeights))
	fmt.Println()

	// ------------------------------------------------------------------
	// 3. Petlja treniranja
	// ------------------------------------------------------------------
	fmt.Println("[3] Pocinje treniranje...")
	fmt.Println()

	for epoch := 0; epoch < epochs; epoch++ {
		totalLoss := 0.0

		for s := 0; s < numSamples; s++ {
			// Izvlacimo ulazne piksele i ciljnu labelu
			xRow := make([]*Value, inputSize)
			for j := 0; j < inputSize; j++ {
				xRow[j] = NewValue(inputs.At(s, j), fmt.Sprintf("x_%d", j))
			}
			yTrue := make([]float64, outputSize)
			for j := 0; j < outputSize; j++ {
				yTrue[j] = labels.At(s, j)
			}

			// ==========================================================
			// TODO 1: Implementiraj forward pass koristeci Value operacije.
			// ==========================================================
			
			// ==========================================================
			// Korak A - Skriveni sloj:
			//
			// Setite se:
			// a = sigmoid(z)
			// z = w*x+b
			//
			//   Za svaki neuron j u skrivenom sloju:
			//		 z = B[j]
			//     za svaki ulaz i:
			//       z = z.Add(xRow[i].Mul(W1[i][j]))
			//     hidden[j] = z.Sigmoid()
			// ==========================================================

			hidden := make([]*Value, hiddenSize)
			for j := 0; j < hiddenSize; j++ {
				hidden[j] = NewValue(0.0, "h") // Zamenite forward pass-om iznad
			}

			// ==========================================================
			// Korak B - Izlazni sloj:
			//
			// Pokusajte sami... slicno kao za skriveni sloj
			//
			// HINT: koristite vektore B2 i hidden, i W2 matricu ; rezultat u output vektor
			//    ...nemojte zaboraviti sigmoid!
			// ==========================================================
			
			output := make([]*Value, outputSize)
			for k := 0; k < outputSize; k++ {
				output[k] = NewValue(0.0, "o") // Zamenite forward pass-om iznad
			}


			// ==========================================================
			// TODO 2: Izracunaj MSE loss koristeci Value operacije.
			//
			// loss = (1/outputSize) * sum((yTrue[k] - output[k])^2)
			//
			// Hint:
			//   loss := NewValue(0.0, "loss")
			//   for k := 0; k < outputSize; k++ {
			//       target := NewValue(yTrue[k], "target")
			//       ...
			//   }
			//
			// Koristite operacije definisane nad Value! (iz engine.go)
			// ==========================================================
			loss := NewValue(0.0, "loss") // Zamenite implementacijom racunanja loss-a

			totalLoss += loss.Data

			// ==========================================================
			// TODO 3: Pozovi loss.Backward() i azuriraj tezine.
			//
			// Ovo je MAGIJA automatske diferencijacije!
			// Jedan poziv racuna SVE gradijente automatski.
			//
			// Korak A: loss.Backward()
			//
			// Korak B: Azuriraj sve tezine:
			//   for _, w := range allWeights {
			//       w.Data -= learningRate * w.Grad
			//   }
			// ==========================================================
			_ = loss // Obrisite kada pozovete loss.Backward()
		}

		avgLoss := totalLoss / float64(numSamples)
		lossHistory = append(lossHistory, avgLoss)

		fmt.Printf("    Epoha %d/%d | Gubitak: %.6f\n", epoch+1, epochs, avgLoss)
	}

	fmt.Println()
	fmt.Println("    Treniranje zavrseno!")
	fmt.Println()

	// ------------------------------------------------------------------
	// 4. Vizuelizacija rezultata
	// ------------------------------------------------------------------
	fmt.Println("[4] Vizuelizacija rezultata...")

	if err := helpers.PlotLoss(lossHistory, "autodiff_loss.png"); err != nil {
		fmt.Printf("    GRESKA: %v\n", err)
	}

	// Generisemo predikcije za prvih 10 uzoraka
	if err := os.MkdirAll("autodiff_predictions", 0755); err != nil {
		fmt.Printf("    GRESKA: %v\n", err)
	}

	numPreview := 10
	if numSamples < numPreview {
		numPreview = numSamples
	}

	for i := 0; i < numPreview; i++ {
		pixels := make([]float64, inputSize)
		xRow := make([]*Value, inputSize)
		for j := 0; j < inputSize; j++ {
			pixels[j] = inputs.At(i, j)
			xRow[j] = NewValue(pixels[j], "x")
		}

		// Forward pass za predikciju
		hidden := make([]*Value, hiddenSize)
		for j := 0; j < hiddenSize; j++ {
			sum := B1[j]
			for ii := 0; ii < inputSize; ii++ {
				sum = sum.Add(xRow[ii].Mul(W1[ii][j]))
			}
			hidden[j] = sum.Sigmoid()
		}

		output := make([]*Value, outputSize)
		for k := 0; k < outputSize; k++ {
			sum := B2[k]
			for j := 0; j < hiddenSize; j++ {
				sum = sum.Add(hidden[j].Mul(W2[j][k]))
			}
			output[k] = sum.Sigmoid()
		}

		// Pronalazimo klasu sa najvecom verovatnocom
		prediction := 0
		maxProb := output[0].Data
		for k := 1; k < outputSize; k++ {
			if output[k].Data > maxProb {
				maxProb = output[k].Data
				prediction = k
			}
		}

		// Prava labela
		trueLabel := 0
		for j := 0; j < outputSize; j++ {
			if labels.At(i, j) == 1.0 {
				trueLabel = j
				break
			}
		}

		path := fmt.Sprintf("autodiff_predictions/uzorak_%02d_tacno_%d_pred_%d.png",
			i, trueLabel, prediction)
		if err := helpers.SaveDigitImage(pixels, prediction, path); err != nil {
			fmt.Printf("    GRESKA: %v\n", err)
		}
	}

	fmt.Println()
	fmt.Println("============================================")
	fmt.Println("  Gotovo! Proverite autodiff_loss.png i")
	fmt.Println("  autodiff_predictions/ folder.")
	fmt.Println()
	fmt.Println("  KLJUCNI UVID: Uporedite ovaj kod sa")
	fmt.Println("  Playground 2. Nema vise rucnih izvoda!")
	fmt.Println("  Samo loss.Backward() radi sav posao.")
	fmt.Println("============================================")

	// Sprecavamo upozorenja kompajlera
	_ = learningRate
	_ = allWeights
}
