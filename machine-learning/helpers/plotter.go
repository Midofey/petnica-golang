// Pomocne funkcije za vizuelizaciju: grafikon gubitka i animacija fitovanja.
// Ucenici NE menjaju ovaj fajl.

package helpers

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"math"
	"os"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// -----------------------------------------------------------------------
// Grafikon gubitka (Loss Plot)
// -----------------------------------------------------------------------

// PlotLoss prima istoriju gubitka (loss) i cuva PNG grafikon.
// X osa = epoha, Y osa = MSE gubitak.
func PlotLoss(history []float64, path string) error {
	p := plot.New()
	p.Title.Text = "Gubitak tokom treniranja (Training Loss)"
	p.X.Label.Text = "Epoha"
	p.Y.Label.Text = "MSE gubitak"

	// Pretvaramo istoriju u plotter tacke
	pts := make(plotter.XYs, len(history))
	for i, loss := range history {
		pts[i].X = float64(i + 1)
		pts[i].Y = loss
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		return fmt.Errorf("greska pri pravljenju linije: %w", err)
	}
	line.Color = plotutil.Color(0)
	line.Width = vg.Points(2)

	p.Add(line)
	p.Add(plotter.NewGrid())

	if err := p.Save(8*vg.Inch, 5*vg.Inch, path); err != nil {
		return fmt.Errorf("greska pri cuvanju grafikona: %w", err)
	}

	fmt.Printf("  [Plotter] Grafikon gubitka sacuvan: %s\n", path)
	return nil
}

// -----------------------------------------------------------------------
// Animacija fitovanja krive (Curve Fitting Animation)
// -----------------------------------------------------------------------

// AnimationFrame cuva stanje jednog kadra animacije.
type AnimationFrame struct {
	Epoch       int
	Predictions []float64
}

// Animator prikuplja kadrove tokom treniranja i generise GIF.
type Animator struct {
	OriginalX []float64
	OriginalY []float64
	Frames    []AnimationFrame
}

// NewAnimator pravi novi Animator sa originalnim podacima.
func NewAnimator(x, y []float64) *Animator {
	return &Animator{
		OriginalX: x,
		OriginalY: y,
		Frames:    make([]AnimationFrame, 0),
	}
}

// CaptureFrame snima trenutne predikcije mreze kao jedan kadar.
// Poziva se periodicno tokom treniranja (npr. svakih 100 epoha).
func (a *Animator) CaptureFrame(epoch int, predictions []float64) {
	predCopy := make([]float64, len(predictions))
	copy(predCopy, predictions)
	a.Frames = append(a.Frames, AnimationFrame{
		Epoch:       epoch,
		Predictions: predCopy,
	})
}

// SaveAnimation generise animirani GIF od snimljenih kadrova.
// Svaki kadar prikazuje originalne tacke (crvene) i trenutnu predikciju (plava kriva).
func (a *Animator) SaveAnimation(path string) error {
	if len(a.Frames) == 0 {
		return fmt.Errorf("nema snimljenih kadrova za animaciju")
	}

	const width = 640
	const height = 480

	// Pronalazimo opseg podataka za skaliranje
	minX, maxX := a.OriginalX[0], a.OriginalX[0]
	minY, maxY := 0.0, 1.0
	for _, x := range a.OriginalX {
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
	}
	for _, y := range a.OriginalY {
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	padX := (maxX - minX) * 0.05
	padY := (maxY - minY) * 0.1
	minX -= padX
	maxX += padX
	minY -= padY
	maxY += padY

	var images []*image.Paletted
	var delays []int

	gifPalette := buildPalette()

	for _, frame := range a.Frames {
		img := image.NewPaletted(image.Rect(0, 0, width, height), gifPalette)

		// Bela pozadina
		bgColor := uint8(findClosestPaletteIndex(gifPalette, color.RGBA{255, 255, 255, 255}))
		for py := 0; py < height; py++ {
			for px := 0; px < width; px++ {
				img.SetColorIndex(px, py, bgColor)
			}
		}

		// Crtamo originalne tacke (crvene)
		redIdx := uint8(findClosestPaletteIndex(gifPalette, color.RGBA{220, 50, 50, 255}))
		for i := range a.OriginalX {
			px := int(mapRange(a.OriginalX[i], minX, maxX, 40, float64(width-20)))
			py := int(mapRange(a.OriginalY[i], minY, maxY, float64(height-30), 20))
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					nx, ny := px+dx, py+dy
					if nx >= 0 && nx < width && ny >= 0 && ny < height {
						img.SetColorIndex(nx, ny, redIdx)
					}
				}
			}
		}

		// Crtamo predikciju mreze (plava linija)
		blueIdx := uint8(findClosestPaletteIndex(gifPalette, color.RGBA{30, 100, 220, 255}))

		type xyPair struct{ x, y float64 }
		pairs := make([]xyPair, len(a.OriginalX))
		for i := range a.OriginalX {
			predVal := 0.5
			if i < len(frame.Predictions) {
				predVal = frame.Predictions[i]
			}
			pairs[i] = xyPair{a.OriginalX[i], predVal}
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].x < pairs[j].x })

		for i := 1; i < len(pairs); i++ {
			x0 := int(mapRange(pairs[i-1].x, minX, maxX, 40, float64(width-20)))
			y0 := int(mapRange(pairs[i-1].y, minY, maxY, float64(height-30), 20))
			x1 := int(mapRange(pairs[i].x, minX, maxX, 40, float64(width-20)))
			y1 := int(mapRange(pairs[i].y, minY, maxY, float64(height-30), 20))
			drawLine(img, x0, y0, x1, y1, blueIdx, width, height)
		}

		// Broj epohe u gornjem levom uglu
		drawEpochLabel(img, frame.Epoch, gifPalette, width)

		images = append(images, img)
		delays = append(delays, 15)
	}

	// Zadrzavanje na poslednjem kadru
	if len(delays) > 0 {
		delays[len(delays)-1] = 200
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("greska pri kreiranju GIF fajla: %w", err)
	}
	defer f.Close()

	if err := gif.EncodeAll(f, &gif.GIF{
		Image:     images,
		Delay:     delays,
		LoopCount: 0,
	}); err != nil {
		return fmt.Errorf("greska pri enkodiranju GIF-a: %w", err)
	}

	fmt.Printf("  [Plotter] Animacija sacuvana: %s (%d kadrova)\n", path, len(a.Frames))
	return nil
}

// -----------------------------------------------------------------------
// Pomocne funkcije za crtanje
// -----------------------------------------------------------------------

func buildPalette() color.Palette {
	p := make(color.Palette, 0, 256)
	p = append(p, palette.Plan9...)
	p = append(p, color.RGBA{255, 255, 255, 255})
	p = append(p, color.RGBA{220, 50, 50, 255})
	p = append(p, color.RGBA{30, 100, 220, 255})
	p = append(p, color.RGBA{50, 50, 50, 255})
	if len(p) > 256 {
		p = p[:256]
	}
	return p
}

func findClosestPaletteIndex(p color.Palette, target color.Color) int {
	bestIdx := 0
	bestDist := math.MaxFloat64
	tr, tg, tb, _ := target.RGBA()
	for i, c := range p {
		cr, cg, cb, _ := c.RGBA()
		dr := float64(tr) - float64(cr)
		dg := float64(tg) - float64(cg)
		db := float64(tb) - float64(cb)
		dist := dr*dr + dg*dg + db*db
		if dist < bestDist {
			bestDist = dist
			bestIdx = i
		}
	}
	return bestIdx
}

func mapRange(value, inMin, inMax, outMin, outMax float64) float64 {
	if inMax == inMin {
		return (outMin + outMax) / 2
	}
	return outMin + (value-inMin)*(outMax-outMin)/(inMax-inMin)
}

// drawLine crta liniju (Bresenham algoritam) na paletizovanoj slici.
func drawLine(img *image.Paletted, x0, y0, x1, y1 int, colorIdx uint8, w, h int) {
	dx := intAbs(x1 - x0)
	dy := intAbs(y1 - y0)
	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	sy := 1
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	for {
		if x0 >= 0 && x0 < w && y0 >= 0 && y0 < h {
			img.SetColorIndex(x0, y0, colorIdx)
			if x0+1 < w {
				img.SetColorIndex(x0+1, y0, colorIdx)
			}
			if y0+1 < h {
				img.SetColorIndex(x0, y0+1, colorIdx)
			}
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// drawEpochLabel iscrtava broj epohe u gornjem levom uglu slike.
func drawEpochLabel(img *image.Paletted, epoch int, pal color.Palette, width int) {
	darkIdx := uint8(findClosestPaletteIndex(pal, color.RGBA{50, 50, 50, 255}))
	label := fmt.Sprintf("E: %d", epoch)
	startX := 5
	startY := 5

	bgIdx := uint8(findClosestPaletteIndex(pal, color.RGBA{240, 240, 240, 255}))
	for py := startY; py < startY+12; py++ {
		for px := startX; px < startX+len(label)*6+4; px++ {
			if px < width && py < 480 {
				img.SetColorIndex(px, py, bgIdx)
			}
		}
	}

	drawSimpleText(img, label, startX+2, startY+2, darkIdx, width, 480)
}

// drawSimpleText crta tekst koristeci minimalan 5x7 bitmap font.
func drawSimpleText(img *image.Paletted, text string, x, y int, colorIdx uint8, w, h int) {
	fontMap := map[byte][7]uint8{
		'0': {0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E},
		'1': {0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E},
		'2': {0x0E, 0x11, 0x01, 0x06, 0x08, 0x10, 0x1F},
		'3': {0x0E, 0x11, 0x01, 0x06, 0x01, 0x11, 0x0E},
		'4': {0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02},
		'5': {0x1F, 0x10, 0x1E, 0x01, 0x01, 0x11, 0x0E},
		'6': {0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E},
		'7': {0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08},
		'8': {0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E},
		'9': {0x0E, 0x11, 0x11, 0x0F, 0x01, 0x02, 0x0C},
		'E': {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x1F},
		':': {0x00, 0x04, 0x04, 0x00, 0x04, 0x04, 0x00},
		' ': {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}

	cx := x
	for i := 0; i < len(text); i++ {
		ch := text[i]
		glyph, ok := fontMap[ch]
		if !ok {
			cx += 6
			continue
		}
		for row := 0; row < 7; row++ {
			for col := 0; col < 5; col++ {
				if glyph[row]&(1<<uint(4-col)) != 0 {
					px, py := cx+col, y+row
					if px >= 0 && px < w && py >= 0 && py < h {
						img.SetColorIndex(px, py, colorIdx)
					}
				}
			}
		}
		cx += 6
	}
}
