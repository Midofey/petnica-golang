// ============================================================================
// Playground 3: Autodiff Engine (Micrograd u Go-u)
// ============================================================================
// Ovaj fajl definise "Value" strukturu - srce automatske diferencijacije.
// Svaki Value cuva: podatak, gradijent, decu, i backward funkciju.
//
// Ucenici implementiraju: Add(), Mul(), Sigmoid()
// Vec je implementirano: Value struct, NewValue() i Backward() sa topoloskim sortiranjem
// ============================================================================

package main

import "math"

// Value je osnovni gradivni blok naseg autodiff engine-a.
// Svaki Value predstavlja JEDAN BROJ u racunskom grafu.
//   - Data:     vrednost (rezultat operacije)
//   - Grad:     gradijent (izracunat tokom backward pass-a)
//   - backward: funkcija koja racuna gradijente roditelja
//   - children: pokazivaci na decu (ulaze u operaciju)
//   - label:    opcioni naziv za debagovanje
type Value struct {
	Data     float64
	Grad     float64
	backward func()
	children []*Value
	label    string
}

// NewValue pravi novi Value sa zadatom vrednoscu i labelom.
func NewValue(data float64, label string) *Value {
	return &Value{
		Data:     data,
		Grad:     0.0,
		backward: func() {}, // Podrazumevano: nista ne radi
		children: nil,
		label:    label,
	}
}

// ==========================================================
// TODO 2: Implementiraj Add metodu.
// Sabira dva Value-a i vraca novi Value.
//
// **Ovde samo treba da prekucate kod ispod:**
//
// Forward: out.Data = a.Data + b.Data
//
// Backward: Gradijent sabiranja se samo PRENOSI.
//
//	a.Grad += out.Grad    (jer d(a+b)/da = 1)
//	b.Grad += out.Grad    (jer d(a+b)/db = 1)
//
// ==========================================================
func (a *Value) Add(b *Value) *Value {
	out := &Value{
		Data:     0.0, // Ovde ide sabiranje Data polja od a i b (forward pass)
		children: []*Value{a, b},
		label:    "(" + a.label + "+" + b.label + ")",
	}

	out.backward = func() {
		// ovde ide racunanje gradijenta i dodavanje istog na a i b (backward pass)
		// a.Grad += ...
		// b.Grad += ...
	}

	return out
}

// ==========================================================
// TODO 3: Implementiraj Mul metodu.
// Mnozi dva Value-a i vraca novi Value.
//
// Forward: out = a * b
//
// **Hint (Backward)**: d(a*b)/da = b,
//
//	d(a*b)/db = a
//
// ==========================================================
func (a *Value) Mul(b *Value) *Value {
	out := &Value{
		Data:     0.0, // Ovde implementirajte mnozenje a i b (forward pass)
		children: []*Value{a, b},
		label:    "(" + a.label + "*" + b.label + ")",
	}

	out.backward = func() {
		// ovde ide racunanje gradijenta i dodavanje istog na a i b (backward pass)
		// a.Grad +=
		// b.Grad +=
	}

	return out
}

// ==========================================================
// TODO 4: Implementiraj Sigmoid metodu.
// Primenjuje sigmoid aktivaciju na Value.
//
// Forward: out = sigmoid(a)
//
// **Hint (Backward)**:
//
//	(d_sigmoid(x) = sigmoid(x) * (1 - sigmoid(x)))
//
// ==========================================================
func (a *Value) Sigmoid() *Value {
	out := &Value{
		Data:     0.0, // Ovde ide implementacija sigmoida nad a (forward pass)
		children: []*Value{a},
		label:    "sig(" + a.label + ")",
	}

	out.backward = func() {
		// ovde ide racunanje gradijenta i dodavanje istog na a (backward pass)
		// a.Grad +=
	}

	return out
}

// Sub oduzima dva Value-a: a - b = a + (-1 * b)
// Ova metoda je vec implementirana koristeci Add i Mul.
func (a *Value) Sub(b *Value) *Value {
	negOne := NewValue(-1.0, "-1")
	negB := negOne.Mul(b)
	return a.Add(negB)
}

// ============================================================================
// Backward() - Automatska diferencijacija putem topoloskog sortiranja
// ============================================================================
//
// Algoritam:
// 1. Postavi gradijent izlaza na 1.0 (dL/dL = 1)
// 2. Topoloski sortiraj sve cvorove u grafu (od listova ka korenu)
// 3. Prodji kroz cvorove obrnutim redosledom i pozovi backward() za svaki
// ============================================================================
func (v *Value) Backward() {
	// Topolosko sortiranje koristeci DFS (pretraga u dubinu)
	visited := make(map[*Value]bool)
	sorted := make([]*Value, 0)

	var topoSort func(node *Value)
	topoSort = func(node *Value) {
		if visited[node] {
			return
		}
		visited[node] = true
		for _, child := range node.children {
			topoSort(child)
		}
		sorted = append(sorted, node)
	}

	topoSort(v)

	// Resetujemo sve gradijente na 0
	for _, node := range sorted {
		node.Grad = 0.0
	}

	// Gradijent izlaznog cvora je 1.0 (polazna tacka)
	v.Grad = 1.0

	// Prolazimo obrnutim redosledom (od izlaza ka ulazima)
	for i := len(sorted) - 1; i >= 0; i-- {
		sorted[i].backward()
	}
}

// Sprecavamo upozorenje kompajlera za neiskoriscen import
var _ = math.Exp
