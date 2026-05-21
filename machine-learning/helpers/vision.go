// Pomocne funkcije za vizuelizaciju MNIST slika.
// Ucenici NE menjaju ovaj fajl.

package helpers

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// SaveDigitImage konvertuje niz od 784 float64 vrednosti (0.0-1.0) u grayscale PNG sliku.
// Iscrtava predvidjeni broj u gornjem desnom uglu slike.
// pixelArray: 784 vrednosti (28x28 piksela), normalizovane 0.0-1.0
// prediction: predvidjeni broj (0-9)
// path: putanja do izlaznog PNG fajla
func SaveDigitImage(pixelArray []float64, prediction int, path string) error {
	if len(pixelArray) != 784 {
		return fmt.Errorf("ocekivano 784 piksela, dobijeno %d", len(pixelArray))
	}

	// Skaliramo sliku 4x da bude bolje vidljiva (28*4 = 112 piksela)
	const scale = 4
	const origSize = 28
	const imgSize = origSize * scale

	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))

	// Crtamo piksele
	for row := 0; row < origSize; row++ {
		for col := 0; col < origSize; col++ {
			idx := row*origSize + col
			// Invertujemo: 0.0 = belo, 1.0 = crno (standard za MNIST prikaz)
			grayVal := uint8((1.0 - pixelArray[idx]) * 255.0)
			c := color.RGBA{grayVal, grayVal, grayVal, 255}

			// Popunjavamo scale x scale blok
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.Set(col*scale+dx, row*scale+dy, c)
				}
			}
		}
	}

	// Crtamo predvidjeni broj u gornjem desnom uglu
	drawPredictionLabel(img, prediction, imgSize)

	// Cuvamo PNG fajl
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("greska pri kreiranju fajla %s: %w", path, err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("greska pri enkodiranju PNG slike: %w", err)
	}

	fmt.Printf("  [Vision] Slika sacuvana: %s (predikcija: %d)\n", path, prediction)
	return nil
}

// drawPredictionLabel iscrtava predvidjeni broj u gornjem desnom uglu slike.
func drawPredictionLabel(img *image.RGBA, prediction int, imgSize int) {
	// Pozadina za labelu (crveni pravougaonik)
	labelBg := color.RGBA{220, 50, 50, 255}
	textColor := color.RGBA{255, 255, 255, 255}

	labelW := 24
	labelH := 20
	startX := imgSize - labelW - 4
	startY := 4

	// Crtamo pozadinu labele
	for py := startY; py < startY+labelH; py++ {
		for px := startX; px < startX+labelW; px++ {
			if px >= 0 && px < imgSize && py >= 0 && py < imgSize {
				img.Set(px, py, labelBg)
			}
		}
	}

	// Crtamo cifru koristeci bitmap font (skalirane na 2x)
	drawDigitGlyph(img, prediction, startX+6, startY+2, textColor, imgSize)
}

// drawDigitGlyph iscrtava jednu cifru (0-9) koristeci 5x7 bitmap font, skaliran 2x.
func drawDigitGlyph(img *image.RGBA, digit int, x, y int, c color.RGBA, imgSize int) {
	digitFont := [10][7]uint8{
		{0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E}, // 0
		{0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E}, // 1
		{0x0E, 0x11, 0x01, 0x06, 0x08, 0x10, 0x1F}, // 2
		{0x0E, 0x11, 0x01, 0x06, 0x01, 0x11, 0x0E}, // 3
		{0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02}, // 4
		{0x1F, 0x10, 0x1E, 0x01, 0x01, 0x11, 0x0E}, // 5
		{0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E}, // 6
		{0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08}, // 7
		{0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E}, // 8
		{0x0E, 0x11, 0x11, 0x0F, 0x01, 0x02, 0x0C}, // 9
	}

	if digit < 0 || digit > 9 {
		return
	}

	glyph := digitFont[digit]
	const scale = 2

	for row := 0; row < 7; row++ {
		for col := 0; col < 5; col++ {
			if glyph[row]&(1<<uint(4-col)) != 0 {
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						px := x + col*scale + dx
						py := y + row*scale + dy
						if px >= 0 && px < imgSize && py >= 0 && py < imgSize {
							img.Set(px, py, c)
						}
					}
				}
			}
		}
	}
}

// SaveMultipleDigits cuva vise MNIST slika u jednom pozivu.
// Koristi se na kraju treniranja za brzu vizuelizaciju rezultata.
func SaveMultipleDigits(inputs [][]float64, predictions []int, dirPath string) error {
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("greska pri kreiranju direktorijuma %s: %w", dirPath, err)
	}

	count := len(inputs)
	if count > len(predictions) {
		count = len(predictions)
	}

	for i := 0; i < count; i++ {
		path := fmt.Sprintf("%s/digit_%02d_pred_%d.png", dirPath, i, predictions[i])
		if err := SaveDigitImage(inputs[i], predictions[i], path); err != nil {
			return err
		}
	}

	fmt.Printf("  [Vision] Sacuvano %d slika u %s\n", count, dirPath)
	return nil
}
