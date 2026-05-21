// ============================================================================
// webui.go — Interaktivni web interfejs za crtanje i predikciju cifara
// ============================================================================
// Pokrece lokalni HTTP server sa HTML canvas-om gde korisnik moze da
// nacrta cifru, a treniran model predvidja koja je cifra.
// ============================================================================

package helpers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// PredictFunc je tip funkcije koju web server koristi za predikciju.
// Prima 784 normalizovanih piksela (28x28), vraca predvidjenu klasu i verovatnoce.
type PredictFunc func(pixels []float64) (prediction int, probabilities []float64)

// predictRequest je JSON struktura koju salje JavaScript iz browsera.
type predictRequest struct {
	Pixels []float64 `json:"pixels"` // 784 float64 vrednosti (28x28)
}

// predictResponse je JSON odgovor koji saljemo nazad u browser.
type predictResponse struct {
	Prediction    int       `json:"prediction"`
	Probabilities []float64 `json:"probabilities"`
}

// StartDrawingServer pokrece interaktivni web server za crtanje cifara.
// predict je funkcija koja prima 784 piksela i vraca predikciju.
// Server blokira (ne vraca se) dok korisnik ne zatvori program.
func StartDrawingServer(port int, predict PredictFunc) error {
	mux := http.NewServeMux()

	// Glavna stranica sa canvas-om
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, drawingPageHTML)
	})

	// API endpoint za predikciju
	mux.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Samo POST metod je dozvoljen", http.StatusMethodNotAllowed)
			return
		}

		var req predictRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Neispravan JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if len(req.Pixels) != 784 {
			http.Error(w, fmt.Sprintf("Ocekivano 784 piksela, dobijeno %d", len(req.Pixels)), http.StatusBadRequest)
			return
		}

		prediction, probs := predict(req.Pixels)

		resp := predictResponse{
			Prediction:    prediction,
			Probabilities: probs,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Pokusamo trazeni port, a ako je zauzet koristimo :0 (automatski slobodan port)
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		// Automatski biramo slobodan port
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			return fmt.Errorf("ne mogu da pokrenem server: %w", err)
		}
	}

	actualPort := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d", actualPort)
	fmt.Println()
	fmt.Println("============================================")
	fmt.Println("  Interaktivni mod: crtanje cifara!")
	fmt.Printf("  Otvorite u browseru: %s\n", url)
	fmt.Println("  (Ctrl+C za zaustavljanje)")
	fmt.Println("============================================")

	return http.Serve(listener, mux)
}

// drawingPageHTML sadrzi kompletnu HTML stranicu sa canvas-om za crtanje.
const drawingPageHTML = `<!DOCTYPE html>
<html lang="sr">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>MNIST Prediktor — Petnica DL Workshop</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap" rel="stylesheet">
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }

  body {
    font-family: 'Inter', system-ui, sans-serif;
    background: #0f0f1a;
    color: #e0e0e8;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem;
  }

  h1 {
    font-size: 1.8rem;
    font-weight: 700;
    background: linear-gradient(135deg, #7c5cff, #00d4ff);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    margin-bottom: 0.3rem;
  }

  .subtitle {
    font-size: 0.9rem;
    color: #888;
    margin-bottom: 2rem;
  }

  .app {
    display: flex;
    gap: 2.5rem;
    align-items: flex-start;
    flex-wrap: wrap;
    justify-content: center;
  }

  .draw-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  canvas {
    border: 2px solid #333;
    border-radius: 12px;
    cursor: crosshair;
    background: #000;
    touch-action: none;
    box-shadow: 0 4px 24px rgba(124, 92, 255, 0.15);
  }

  .buttons {
    display: flex;
    gap: 0.75rem;
  }

  button {
    font-family: 'Inter', sans-serif;
    font-size: 0.95rem;
    font-weight: 600;
    padding: 0.6rem 1.5rem;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-predict {
    background: linear-gradient(135deg, #7c5cff, #5c3ddb);
    color: white;
  }
  .btn-predict:hover { transform: translateY(-1px); box-shadow: 0 4px 16px rgba(124, 92, 255, 0.4); }
  .btn-predict:active { transform: translateY(0); }

  .btn-clear {
    background: #1e1e2e;
    color: #aaa;
    border: 1px solid #333;
  }
  .btn-clear:hover { background: #2a2a3e; color: #ddd; }

  .result-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.2rem;
    min-width: 280px;
  }

  .prediction-display {
    width: 140px;
    height: 140px;
    border-radius: 20px;
    background: linear-gradient(135deg, #1a1a2e, #16213e);
    border: 2px solid #333;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 5rem;
    font-weight: 700;
    color: #7c5cff;
    transition: all 0.3s ease;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3);
  }

  .prediction-display.active {
    border-color: #7c5cff;
    box-shadow: 0 4px 32px rgba(124, 92, 255, 0.3);
  }

  .confidence-label {
    font-size: 0.85rem;
    color: #888;
  }

  .bars {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .bar-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .bar-label {
    width: 1.2rem;
    text-align: right;
    font-size: 0.85rem;
    font-weight: 600;
    color: #888;
  }

  .bar-track {
    flex: 1;
    height: 18px;
    background: #1a1a2e;
    border-radius: 4px;
    overflow: hidden;
  }

  .bar-fill {
    height: 100%;
    border-radius: 4px;
    background: linear-gradient(90deg, #5c3ddb, #7c5cff);
    transition: width 0.4s cubic-bezier(0.16, 1, 0.3, 1);
    width: 0%;
  }

  .bar-fill.top {
    background: linear-gradient(90deg, #00b4d8, #00d4ff);
  }

  .bar-pct {
    width: 3rem;
    font-size: 0.8rem;
    color: #666;
    text-align: right;
  }

  .preview-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    margin-top: 0.5rem;
  }

  .preview-label {
    font-size: 0.75rem;
    color: #555;
  }

  .preview-canvas {
    border: 1px solid #333;
    border-radius: 4px;
    image-rendering: pixelated;
  }
</style>
</head>
<body>

<h1>Nacrtaj cifru</h1>
<p class="subtitle">Petnica DL Workshop — Interaktivni MNIST prediktor</p>

<div class="app">
  <div class="draw-section">
    <canvas id="drawCanvas" width="280" height="280"></canvas>
    <div class="preview-section">
      <span class="preview-label">Preprocessed (28×28):</span>
      <canvas id="previewCanvas" class="preview-canvas" width="28" height="28" style="width:84px;height:84px;"></canvas>
    </div>
    <div class="buttons">
      <button class="btn-predict" onclick="predict()">Predvidi</button>
      <button class="btn-clear" onclick="clearCanvas()">Obrisi</button>
    </div>
  </div>

  <div class="result-section">
    <div class="prediction-display" id="predDisplay">?</div>
    <div class="confidence-label" id="confLabel">Nacrtajte cifru i kliknite "Predvidi"</div>
    <div class="bars" id="bars">
      <div class="bar-row"><span class="bar-label">0</span><div class="bar-track"><div class="bar-fill" id="bar0"></div></div><span class="bar-pct" id="pct0">—</span></div>
      <div class="bar-row"><span class="bar-label">1</span><div class="bar-track"><div class="bar-fill" id="bar1"></div></div><span class="bar-pct" id="pct1">—</span></div>
      <div class="bar-row"><span class="bar-label">2</span><div class="bar-track"><div class="bar-fill" id="bar2"></div></div><span class="bar-pct" id="pct2">—</span></div>
      <div class="bar-row"><span class="bar-label">3</span><div class="bar-track"><div class="bar-fill" id="bar3"></div></div><span class="bar-pct" id="pct3">—</span></div>
      <div class="bar-row"><span class="bar-label">4</span><div class="bar-track"><div class="bar-fill" id="bar4"></div></div><span class="bar-pct" id="pct4">—</span></div>
      <div class="bar-row"><span class="bar-label">5</span><div class="bar-track"><div class="bar-fill" id="bar5"></div></div><span class="bar-pct" id="pct5">—</span></div>
      <div class="bar-row"><span class="bar-label">6</span><div class="bar-track"><div class="bar-fill" id="bar6"></div></div><span class="bar-pct" id="pct6">—</span></div>
      <div class="bar-row"><span class="bar-label">7</span><div class="bar-track"><div class="bar-fill" id="bar7"></div></div><span class="bar-pct" id="pct7">—</span></div>
      <div class="bar-row"><span class="bar-label">8</span><div class="bar-track"><div class="bar-fill" id="bar8"></div></div><span class="bar-pct" id="pct8">—</span></div>
      <div class="bar-row"><span class="bar-label">9</span><div class="bar-track"><div class="bar-fill" id="bar9"></div></div><span class="bar-pct" id="pct9">—</span></div>
    </div>
  </div>
</div>

<script>
const canvas = document.getElementById('drawCanvas');
const ctx = canvas.getContext('2d');
const previewCanvas = document.getElementById('previewCanvas');
const previewCtx = previewCanvas.getContext('2d');

let drawing = false;

// Crni pozadina
ctx.fillStyle = '#000';
ctx.fillRect(0, 0, 280, 280);

// Podesavanja za crtanje
ctx.strokeStyle = '#fff';
ctx.lineWidth = 18;
ctx.lineCap = 'round';
ctx.lineJoin = 'round';

function getPos(e) {
  const rect = canvas.getBoundingClientRect();
  const clientX = e.touches ? e.touches[0].clientX : e.clientX;
  const clientY = e.touches ? e.touches[0].clientY : e.clientY;
  return {
    x: (clientX - rect.left) * (280 / rect.width),
    y: (clientY - rect.top) * (280 / rect.height)
  };
}

canvas.addEventListener('mousedown', e => { drawing = true; ctx.beginPath(); const p = getPos(e); ctx.moveTo(p.x, p.y); });
canvas.addEventListener('mousemove', e => { if (!drawing) return; const p = getPos(e); ctx.lineTo(p.x, p.y); ctx.stroke(); });
canvas.addEventListener('mouseup', () => { drawing = false; updatePreview(); });
canvas.addEventListener('mouseleave', () => { drawing = false; updatePreview(); });

canvas.addEventListener('touchstart', e => { e.preventDefault(); drawing = true; ctx.beginPath(); const p = getPos(e); ctx.moveTo(p.x, p.y); });
canvas.addEventListener('touchmove', e => { e.preventDefault(); if (!drawing) return; const p = getPos(e); ctx.lineTo(p.x, p.y); ctx.stroke(); });
canvas.addEventListener('touchend', e => { e.preventDefault(); drawing = false; updatePreview(); });

function updatePreview() {
  // Prikazujemo preprocessed 28x28 verziju
  previewCtx.clearRect(0, 0, 28, 28);
  previewCtx.drawImage(canvas, 0, 0, 28, 28);
}

function getPixels() {
  // Downscale 280x280 -> 28x28
  const tmpCanvas = document.createElement('canvas');
  tmpCanvas.width = 28;
  tmpCanvas.height = 28;
  const tmpCtx = tmpCanvas.getContext('2d');
  tmpCtx.drawImage(canvas, 0, 0, 28, 28);
  const imgData = tmpCtx.getImageData(0, 0, 28, 28);

  // Konvertujemo u grayscale normalizovan na [0, 1]
  // MNIST je beli tekst na crnoj pozadini, canvas je vec takav
  const pixels = [];
  for (let i = 0; i < 784; i++) {
    // Koristimo samo R kanal (svi kanali su isti za grayscale)
    pixels.push(imgData.data[i * 4] / 255.0);
  }
  return pixels;
}

async function predict() {
  const pixels = getPixels();

  try {
    const resp = await fetch('/predict', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ pixels })
    });

    if (!resp.ok) {
      const errText = await resp.text();
      console.error('Greska:', errText);
      return;
    }

    const data = await resp.json();
    showResult(data.prediction, data.probabilities);
  } catch (err) {
    console.error('Fetch greska:', err);
  }
}

function showResult(prediction, probs) {
  const display = document.getElementById('predDisplay');
  display.textContent = prediction;
  display.classList.add('active');

  const maxProb = Math.max(...probs);
  document.getElementById('confLabel').textContent =
    'Sigurnost: ' + (maxProb * 100).toFixed(1) + '%';

  for (let d = 0; d < 10; d++) {
    const pct = (probs[d] * 100);
    const bar = document.getElementById('bar' + d);
    bar.style.width = pct + '%';

    // Istakni najjacu predikciju
    if (d === prediction) {
      bar.classList.add('top');
    } else {
      bar.classList.remove('top');
    }

    document.getElementById('pct' + d).textContent = pct.toFixed(1) + '%';
  }
}

function clearCanvas() {
  ctx.fillStyle = '#000';
  ctx.fillRect(0, 0, 280, 280);
  previewCtx.clearRect(0, 0, 28, 28);

  document.getElementById('predDisplay').textContent = '?';
  document.getElementById('predDisplay').classList.remove('active');
  document.getElementById('confLabel').textContent = 'Nacrtajte cifru i kliknite "Predvidi"';

  for (let d = 0; d < 10; d++) {
    document.getElementById('bar' + d).style.width = '0%';
    document.getElementById('bar' + d).classList.remove('top');
    document.getElementById('pct' + d).textContent = '\u2014';
  }
}
</script>
</body>
</html>`
