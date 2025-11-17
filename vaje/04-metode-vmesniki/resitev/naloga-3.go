package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ===== GLOBALNE SPREMENLJIVKE =====
var (
	promet     float64
	stNarocil  int
	mutex      sync.Mutex // za promet in stNarocil
	printMutex sync.Mutex // za izpis, da se ne meša med gorutinami
)

// ===== VMESNIK =====
type narocilo interface {
	obdelaj()
}

// ===== STRUKTURE =====

type izdelek struct {
	imeIzdelka string
	cena       float64
	teza       float64
}

type eknjiga struct {
	naslovKnjige string
	cena         float64
}

type spletniTecaj struct {
	imeTecaja   string
	trajanjeUre int
	cenaUre     float64
}

// ===== METODA obdelaj() ZA VSAK TIP =====

func (i izdelek) obdelaj() {
	// Zaklenemo izpis
	printMutex.Lock()
	stNarocil++
	stevilka := stNarocil
	fmt.Printf("Številka naročila: %d\n", stevilka)
	fmt.Println("Ime izdelka:", i.imeIzdelka)
	fmt.Printf("Cena: %.2f €\n", i.cena)
	fmt.Printf("Teža: %.2f kg\n\n", i.teza)
	printMutex.Unlock()

	// Posodobitev prometa
	mutex.Lock()
	promet += i.cena
	mutex.Unlock()
}

func (e eknjiga) obdelaj() {
	printMutex.Lock()
	stNarocil++
	stevilka := stNarocil
	fmt.Printf("Številka naročila: %d\n", stevilka)
	fmt.Println("Naslov e-knjige:", e.naslovKnjige)
	fmt.Printf("Cena: %.2f €\n\n", e.cena)
	printMutex.Unlock()

	mutex.Lock()
	promet += e.cena
	mutex.Unlock()
}

func (s spletniTecaj) obdelaj() {
	cena := float64(s.trajanjeUre) * s.cenaUre

	printMutex.Lock()
	stNarocil++
	stevilka := stNarocil
	fmt.Printf("Številka naročila: %d\n", stevilka)
	fmt.Println("Ime tečaja:", s.imeTecaja)
	fmt.Printf("Trajanje: %d ur\n", s.trajanjeUre)
	fmt.Printf("Cena na uro: %.2f €\n", s.cenaUre)
	fmt.Printf("Skupna cena: %.2f €\n\n", cena)
	printMutex.Unlock()

	mutex.Lock()
	promet += cena
	mutex.Unlock()
}

// ===== GLAVNI PROGRAM =====

func main() {

	rand.Seed(time.Now().UnixNano())

	// Rezina naročil
	var narocila []narocilo

	// Dodamo 10 naključnih naročil
	for i := 0; i < 10; i++ {
		t := rand.Intn(3)
		switch t {
		case 0:
			narocila = append(narocila, izdelek{
				imeIzdelka: fmt.Sprintf("Izdelek %d", i),
				cena:       float64(10 + rand.Intn(200)),
				teza:       rand.Float64()*5 + 0.5,
			})
		case 1:
			narocila = append(narocila, eknjiga{
				naslovKnjige: fmt.Sprintf("E-knjiga %d", i),
				cena:         float64(5 + rand.Intn(50)),
			})
		case 2:
			narocila = append(narocila, spletniTecaj{
				imeTecaja:   fmt.Sprintf("Tečaj %d", i),
				trajanjeUre: rand.Intn(10) + 1,
				cenaUre:     float64(10 + rand.Intn(30)),
			})
		}
	}

	// WaitGroup za gorutine
	var wg sync.WaitGroup

	// Za vsako naročilo zaženemo gorutino
	for _, n := range narocila {
		wg.Add(1)
		go func(n narocilo) {
			defer wg.Done()
			n.obdelaj()
		}(n)
	}

	// Počakamo na konec
	wg.Wait()

	// Končni izpis
	fmt.Println("=======================================")
	fmt.Println("Skupno število naročil:", stNarocil)
	fmt.Printf("Skupni promet: %.2f €\n", promet)
}