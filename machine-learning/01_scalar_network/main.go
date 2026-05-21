// ============================================================================
// Playground 1: Skalarna neuronska mreza (1-1-1)
// ============================================================================
// Ovo je vas prvi korak u duboko ucenje! Gradimo najjednostavniju mogucu
// neuronsku mrezu: jedan ulaz -> jedan skriveni neuron -> jedan izlaz.
//
// Struktura mreze:
//   x --[w1]--> [sigmoid] --[w2]--> [sigmoid] --> predikcija
//
// Cilj: na osnovu sati spavanja studenta, predvideti njegov finalni rezultat.
// Dataset: Student Lifestyle & Academic Performance (Kaggle)
// ============================================================================

package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"petnica-dl-workshop/helpers"
)

// ==========================================================
// TODO 1: Implementiraj sigmoid funkciju.
// Sigmoid pretvara bilo koji broj u vrednost izmedju 0 i 1.
//
// Hint: koristite math.Exp()
// ==========================================================
func sigmoid(z float64) float64 {
	rez := 1.0 / (1.0 + math.Exp(-z))
	return rez // Implementirajte sigmoid funkciju ovde
}

// ==========================================================
// TODO 2: Implementiraj izvod sigmoid funkcije.
// Ovo nam govori "koliko brzo" se sigmoid menja u tacki z.
//
// Hint: mozete koristiti sigmoid() funkciju koju ste vec napisali
// ==========================================================
func sigmoidDerivative(z float64) float64 {
	rez := sigmoid(z) * (1 - sigmoid(z)) // Obrisite ovu liniju
	return rez                           // Implementirajte izvod sigmoid funkcije ovde
}

// ==========================================================
// TODO 3: Implementiraj forward pass (prolaz unapred).
// Pokusaj da implementiras kako smo objasnjavali na tabli
//
// Funkcija vraca sva 4 medjurezultata jer ce nam biti potrebni za backward pass.
// ==========================================================
func forward(x, w1, b1, w2, b2 float64) (z1, a1, z2, a2 float64) {
	z1 = w1*x + b1   // Zameni jednacinom za z1
	a1 = sigmoid(z1) // Zameni jednacinom za a1
	z2 = w2*a1 + b2  // Zameni jednacinom za z2
	a2 = sigmoid(z2) // Zameni jednacinom za a2
	return
}

func main() {
	fmt.Println("============================================")
	fmt.Println("  Playground 1: Skalarna mreza (1-1-1)")
	fmt.Println("  Zadatak: Sati spavanja -> Finalni rezultat")
	fmt.Println("============================================")
	fmt.Println()

	// ------------------------------------------------------------------
	// 1. Ucitavanje podataka
	// ------------------------------------------------------------------
	fmt.Println("[1] Ucitavanje podataka...")

	dataPath := "data/student_performance_finalscore.csv"
	if len(os.Args) > 1 {
		dataPath = os.Args[1]
	}

	// Ucitavamo podatke: sati spavanja (ulaz) i finalni rezultat (cilj)
	// Obe vrednosti se automatski normalizuju na opseg [0.0, 1.0]
	// jer sigmoid moze da vraca samo vrednosti izmedju 0 i 1.
	xData, yData, scoreMin, scoreMax, err := helpers.LoadStudentPerformance(dataPath)
	if err != nil {
		fmt.Printf("    GRESKA: %v\n", err)
		fmt.Println("    Proverite da li postoji fajl:", dataPath)
		return
	}
	fmt.Printf("    Broj uzoraka: %d\n", len(xData))
	fmt.Println()

	// ------------------------------------------------------------------
	// 2. Inicijalizacija tezina i parametara
	// ------------------------------------------------------------------
	fmt.Println("[2] Inicijalizacija tezina...")

	// Nasumicne pocetne tezine (mali brojevi)
	w1 := rand.Float64()*2 - 1 // Tezina: ulaz -> skriveni sloj
	b1 := rand.Float64()*2 - 1 // Bias skrivenog sloja
	w2 := rand.Float64()*2 - 1 // Tezina: skriveni sloj -> izlaz
	b2 := rand.Float64()*2 - 1 // Bias izlaznog sloja

	fmt.Printf("    w1=%.4f, b1=%.4f, w2=%.4f, b2=%.4f\n", w1, b1, w2, b2)
	fmt.Println()

	// Hiperparametri
	learningRate := 0.005
	epochs := 5000

	// Istorija gubitka za grafikon
	lossHistory := make([]float64, 0, epochs)

	// Animator za GIF animaciju fitovanja
	animator := helpers.NewAnimator(xData, yData)

	// ------------------------------------------------------------------
	// 3. Petlja treniranja
	// ------------------------------------------------------------------
	fmt.Println("[3] Pocinje treniranje...")
	fmt.Println()

	for epoch := 0; epoch < epochs; epoch++ {
		totalLoss := 0.0

		// Prikupljamo predikcije za animaciju
		predictions := make([]float64, len(xData))

		for i := 0; i < len(xData); i++ {
			x := xData[i]
			yTrue := yData[i]

			// Forward pass - koristimo funkciju iz TODO 3
			z1, a1, z2, a2 := forward(x, w1, b1, w2, b2)

			predictions[i] = a2

			// ==========================================================
			// TODO 4: Izracunaj Mean Squared Error (MSE) za ovu tacku.
			// Formulu smo pokazali na tabli
			//
			// Hint: koristite math.Pow( _ , 2)
			// ==========================================================
			loss := math.Pow(a2-yTrue, 2) / 2 // Zameni formulom za loss
			totalLoss += loss

			// ==========================================================
			// TODO 5: Implementiraj backward pass (prolaz unazad).
			// Ovo je "srce" ucenja! Racunamo koliko je svaka tezina (i bias)
			// doprinela gresci, i u kom pravcu treba da je pomerimo.
			//
			// ==========================================================
			delta2 := (a2 - yTrue) * sigmoidDerivative(z2) // Zameni formulom za delta2
			dW2 := delta2 * a1                             // Zameni formulom za dW2
			dB2 := delta2                                  // Zameni formulom za dB2

			delta1 := delta2 * w2 * sigmoidDerivative(z1) // Zameni formulom za delta1
			dW1 := delta1 * x                             // Zameni formulom za dW1
			dB1 := delta1                                 // Zameni formulom za dB1

			// ==========================================================
			// TODO 6: Azuriraj tezine koristeci learning rate.
			// Svaku tezinu pomeramo SUPROTNO od izvoda (gradientni spust).
			//
			// ==========================================================
			// Implementirajte azuriranje tezina ovde:
			w1 = w1 - dW1*learningRate
			b1 = b1 - dB1*learningRate
			w2 = w2 - dW2*learningRate
			b2 = b2 - dB2*learningRate
		}

		// Srednji gubitak po epohi
		avgLoss := totalLoss / float64(len(xData))
		lossHistory = append(lossHistory, avgLoss)

		// Snimamo kadar za animaciju svakih 50 epoha
		if epoch%50 == 0 {
			animator.CaptureFrame(epoch, predictions)
		}

		// Stampamo napredak svakih 500 epoha
		if epoch%500 == 0 {
			fmt.Printf("    Epoha %4d/%d | Gubitak: %.6f\n", epoch, epochs, avgLoss)
		}
	}

	fmt.Println()
	fmt.Printf("    Zavrseno treniranje! Finalni gubitak: %.6f\n", lossHistory[len(lossHistory)-1])
	fmt.Printf("    Finalne tezine: w1=%.4f, b1=%.4f, w2=%.4f, b2=%.4f\n", w1, b1, w2, b2)
	fmt.Println()

	// ------------------------------------------------------------------
	// 4. Primer predikcije (de-normalizacija nazad u skalu rezultata)
	// ------------------------------------------------------------------
	fmt.Println("[4] Primer predikcija (nakon treniranja):")
	fmt.Println("    NAPOMENA: Predikcije ce biti korektne tek kada zavrsiste sve TODO-ove!")
	fmt.Println()

	testSleepHours := []float64{5.0, 6.5, 7.0, 8.0}
	for _, sleepHrs := range testSleepHours {
		// Normalizujemo ulaz na isti nacin kao trening podatke
		xNorm := helpers.NormalizeMinMax(sleepHrs, 4.0, 10.0)

		// Koristimo ISTU forward() funkciju iz TODO 3
		_, _, _, predNorm := forward(xNorm, w1, b1, w2, b2)

		// De-normalizujemo predikciju nazad u skalu finalnog rezultata
		predScore := helpers.Denormalize(predNorm, scoreMin, scoreMax)
		fmt.Printf("    Spavanje: %.1f sati -> Predvidjeni rezultat: %.1f\n", sleepHrs, predScore)
	}
	fmt.Println()

	// ------------------------------------------------------------------
	// 5. Cuvanje rezultata
	// ------------------------------------------------------------------
	fmt.Println("[5] Cuvanje rezultata...")

	// Grafikon gubitka
	if err := helpers.PlotLoss(lossHistory, "loss_plot.png"); err != nil {
		fmt.Printf("    GRESKA pri cuvanju grafikona: %v\n", err)
	}

	// GIF animacija fitovanja krive
	if err := animator.SaveAnimation("fitting_animation.gif"); err != nil {
		fmt.Printf("    GRESKA pri cuvanju animacije: %v\n", err)
	}

	fmt.Println()
	fmt.Println("============================================")
	fmt.Println("  Gotovo! Proverite loss_plot.png i")
	fmt.Println("  fitting_animation.gif u ovom folderu.")
	fmt.Println("============================================")

	// Sprecavamo upozorenja kompajlera za neiskoriscene importe
	_ = math.Exp
	_ = learningRate
}
