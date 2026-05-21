// Pomocne funkcije za ucitavanje podataka.
// Ucenici NE menjaju ovaj fajl.

package helpers

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// Load1DData ucitava 1D podatke iz CSV fajla (dve kolone: x, y).
// Vraca dva slice-a float64 vrednosti.
func Load1DData(filepath string) ([]float64, []float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("greska pri otvaranju fajla %s: %w", filepath, err)
	}
	defer file.Close()

	var xs, ys []float64
	scanner := bufio.NewScanner(file)

	lineNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lineNum++

		// Preskoci zaglavlje (prva linija koja nije numerička)
		if lineNum == 1 {
			parts := strings.Split(line, ",")
			if len(parts) >= 2 {
				_, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
				_, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err1 != nil || err2 != nil {
					continue // Zaglavlje, preskacemo
				}
			}
		}

		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}

		x, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return nil, nil, fmt.Errorf("greska pri parsiranju x vrednosti u liniji %d: %w", lineNum, err)
		}

		y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return nil, nil, fmt.Errorf("greska pri parsiranju y vrednosti u liniji %d: %w", lineNum, err)
		}

		xs = append(xs, x)
		ys = append(ys, y)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("greska pri citanju fajla: %w", err)
	}

	fmt.Printf("  [Loader] Ucitano %d tacaka iz %s\n", len(xs), filepath)
	return xs, ys, nil
}

// LoadMNIST ucitava MNIST podatke iz CSV formata.
// Ocekuje fajl gde je prva kolona labela (0-9), a ostatak su 784 piksela (0-255).
// Vraca matricu ulaza (normalizovanu 0.0-1.0) i matricu labela (one-hot encoded).
func LoadMNIST(filepath string) (*mat.Dense, *mat.Dense, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("greska pri otvaranju MNIST fajla %s: %w", filepath, err)
	}
	defer file.Close()

	var allInputs []float64
	var allLabels []float64
	numSamples := 0

	scanner := bufio.NewScanner(file)
	// Povecavamo velicinu bafera za vece linije
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lineNum++

		// Preskoci zaglavlje
		if lineNum == 1 {
			parts := strings.Split(line, ",")
			if len(parts) > 0 {
				_, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
				if err != nil {
					continue // Zaglavlje, preskacemo
				}
			}
		}

		parts := strings.Split(line, ",")
		if len(parts) < 785 { // 1 labela + 784 piksela
			continue
		}

		// Prva kolona je labela
		label, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}

		// One-hot kodiranje labele
		oneHot := OneHotEncode(label)
		allLabels = append(allLabels, oneHot...)

		// Preostalih 784 vrednosti su pikseli (normalizujemo na 0.0 - 1.0)
		for i := 1; i <= 784; i++ {
			pixelVal, err := strconv.ParseFloat(strings.TrimSpace(parts[i]), 64)
			if err != nil {
				pixelVal = 0.0
			}
			allInputs = append(allInputs, pixelVal/255.0)
		}

		numSamples++
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("greska pri citanju MNIST fajla: %w", err)
	}

	if numSamples == 0 {
		return nil, nil, fmt.Errorf("nisu pronadjeni podaci u fajlu %s", filepath)
	}

	// Pravimo matrice: svaki red je jedan uzorak
	inputs := mat.NewDense(numSamples, 784, allInputs)
	labels := mat.NewDense(numSamples, 10, allLabels)

	fmt.Printf("  [Loader] Ucitano %d MNIST uzoraka iz %s\n", numSamples, filepath)
	return inputs, labels, nil
}

// OneHotEncode pretvara celobrojnu labelu (0-9) u one-hot vektor duzine 10.
// Npr. labela 3 -> [0, 0, 0, 1, 0, 0, 0, 0, 0, 0]
func OneHotEncode(label int) []float64 {
	oneHot := make([]float64, 10)
	if label >= 0 && label < 10 {
		oneHot[label] = 1.0
	}
	return oneHot
}

// LoadStudentData ucitava studentski CSV i izvlaci dve kolone po imenu zaglavlja.
// Vraca normalizovane vrednosti (0.0-1.0) i originalne min/max za de-normalizaciju.
// Tipicna upotreba: LoadStudentData("data/students.csv", "Sleep_Hours", "Previous_GPA")
func LoadStudentData(filepath, xColumn, yColumn string) (xs, ys []float64, xMin, xMax, yMin, yMax float64, err error) {
	file, openErr := os.Open(filepath)
	if openErr != nil {
		err = fmt.Errorf("greska pri otvaranju fajla %s: %w", filepath, openErr)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Citamo zaglavlje i trazimo indekse kolona
	if !scanner.Scan() {
		err = fmt.Errorf("fajl %s je prazan", filepath)
		return
	}
	header := strings.Split(strings.TrimSpace(scanner.Text()), ",")
	xIdx, yIdx := -1, -1
	for i, col := range header {
		col = strings.TrimSpace(col)
		if col == xColumn {
			xIdx = i
		}
		if col == yColumn {
			yIdx = i
		}
	}
	if xIdx == -1 {
		err = fmt.Errorf("kolona '%s' nije pronadjena u zaglavlju", xColumn)
		return
	}
	if yIdx == -1 {
		err = fmt.Errorf("kolona '%s' nije pronadjena u zaglavlju", yColumn)
		return
	}

	// Citamo podatke
	var rawX, rawY []float64
	lineNum := 1
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		maxIdx := xIdx
		if yIdx > maxIdx {
			maxIdx = yIdx
		}
		if len(parts) <= maxIdx {
			continue
		}

		xVal, xErr := strconv.ParseFloat(strings.TrimSpace(parts[xIdx]), 64)
		if xErr != nil {
			continue
		}
		yVal, yErr := strconv.ParseFloat(strings.TrimSpace(parts[yIdx]), 64)
		if yErr != nil {
			continue
		}
		rawX = append(rawX, xVal)
		rawY = append(rawY, yVal)
	}

	if scanErr := scanner.Err(); scanErr != nil {
		err = fmt.Errorf("greska pri citanju fajla: %w", scanErr)
		return
	}

	if len(rawX) == 0 {
		err = fmt.Errorf("nema validnih podataka u fajlu %s", filepath)
		return
	}

	// Normalizujemo obe kolone na opseg [0.0, 1.0] (min-max normalizacija)
	xMin, xMax = rawX[0], rawX[0]
	yMin, yMax = rawY[0], rawY[0]
	for _, v := range rawX {
		if v < xMin {
			xMin = v
		}
		if v > xMax {
			xMax = v
		}
	}
	for _, v := range rawY {
		if v < yMin {
			yMin = v
		}
		if v > yMax {
			yMax = v
		}
	}

	xs = make([]float64, len(rawX))
	ys = make([]float64, len(rawY))
	for i := range rawX {
		xs[i] = NormalizeMinMax(rawX[i], xMin, xMax)
		ys[i] = NormalizeMinMax(rawY[i], yMin, yMax)
	}

	fmt.Printf("  [Loader] Ucitano %d tacaka iz %s\n", len(xs), filepath)
	fmt.Printf("           X kolona: '%s' (min=%.2f, max=%.2f)\n", xColumn, xMin, xMax)
	fmt.Printf("           Y kolona: '%s' (min=%.2f, max=%.2f)\n", yColumn, yMin, yMax)
	return
}

// NormalizeMinMax normalizuje vrednost u opseg [0, 1] koristeci min-max skaliranje.
func NormalizeMinMax(value, min, max float64) float64 {
	if max == min {
		return 0.5
	}
	return (value - min) / (max - min)
}

// Denormalize vraca normalizovanu vrednost (0-1) u originalni opseg.
func Denormalize(normalized, min, max float64) float64 {
	return normalized*(max-min) + min
}

// LoadStudentPerformance ucitava studentski dataset i vraca normalizovane parove
// (sati spavanja, finalni rezultat) spremne za treniranje 1-1-1 mreze.
// Vraca: xs, ys (normalizovani 0-1), scoreMin, scoreMax (za de-normalizaciju).
func LoadStudentPerformance(filepath string) (xs, ys []float64, scoreMin, scoreMax float64, err error) {
	allXs, allYs, _, _, yMin, yMax, loadErr := LoadStudentData(filepath, "Sleep_Hours", "Final_Score")
	if loadErr != nil {
		err = loadErr
		return
	}
	return allXs, allYs, yMin, yMax, nil
}

// Sprecavamo upozorenje kompajlera za neiskoriscen import
var _ = math.Exp
