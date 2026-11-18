# Primeri programerskih nalog

## Naloga 1

V jeziku Go želimo napisati program, v katerem glavna gorutina zažene *N* dodatnih gorutin, ki se izvajajo sočasno. Vsaka dodatna gorutina izpiše svojo oznako in zaključi. Program popravite tako, da se bodo vse oznake zagotovo pojavile na standardnem izhodu. Za sinhronizacijo gorutin uporabite kanal `ch`.

```go
package main

import "fmt"

const N int = 5

var ch chan struct{} = make(chan struct{})

func worker(id int) {
	fmt.Println(id)
}

func main() {
	for i := range N {
		go worker(i)
	}
}
```

## Naloga 2

Spodnji program v Go dopolnite tako, da bo nemudoma izpisal sporočilo, poslano v katerikoli kanal.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		for {
			time.Sleep(300 * time.Millisecond)
			ch1 <- "VPSA"
		}
	}()

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			ch2 <- "vpsa"
		}
	}()

	// manjkajoča koda
}
```

## Naloga 3

Program zažene deset gorutin, vsaka izračuna kvadrat svojega števila in rezultat pošlje v kanal `results`. Glavna gorutina rezultate sešteje in izpiše vsoto. Program se ne zaključi pravilno. Pojasnite, kaj gre narobe, in program popravite tako, da bo izpisal vsoto kvadratov. Za čakanje na gorutine uporabite čakalno skupino `sync.WaitGroup`.

```go
package main

import "fmt"

func square(x int, results chan<- int) {
	results <- x * x
}

func main() {
	results := make(chan int)
	for i := 1; i <= 10; i++ {
		go square(i, results)
	}
	sum := 0
	for r := range results {
		sum += r
	}
	fmt.Println(sum)
}
```

## Naloga 4

Napisali smo števec, ki ga lahko sočasno povečuje več gorutin. Kaj izpiše spodnji program? Pojasnite zakaj in program popravite tako, da bo izpisal 100.

```go
package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	m sync.Mutex
	n int
}

func (c Counter) Inc() {
	c.m.Lock()
	c.n++
	c.m.Unlock()
}

func main() {
	var c Counter
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			c.Inc()
		})
	}
	wg.Wait()
	fmt.Println(c.n)
}
```

## Naloga 5

Števni semafor smo namesto s kanalom implementirali s ključavnico in pogojno spremenljivko. Program zažene `N` gorutin, hkrati pa sme delati največ `K` gorutin. Program izpiše nekaj vrstic, nato pa se obesi. Pojasnite, zakaj se to zgodi, in program popravite. V metodi `Acquire` je še ena napaka, ki se pokaže le pri določenem prepletanju gorutin. Poiščite jo, opišite zaporedje dogodkov, ki privede do napačnega delovanja semaforja, in jo odpravite.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

const N int = 10 // število gorutin
const K int = 3  // največje število hkrati delujočih gorutin

type Semaphore struct {
	lock  sync.Mutex
	cond  *sync.Cond
	value int
}

func NewSemaphore(value int) *Semaphore {
	s := &Semaphore{value: value}
	s.cond = sync.NewCond(&s.lock)
	return s
}

func (s *Semaphore) Acquire() {
	s.lock.Lock()
	if s.value == 0 {
		s.cond.Wait()
	}
	s.value--
	s.lock.Unlock()
}

func (s *Semaphore) Release() {
	s.lock.Lock()
	s.value++
	s.lock.Unlock()
}

func main() {
	sem := NewSemaphore(K)
	var wg sync.WaitGroup
	for i := range N {
		wg.Go(func() {
			sem.Acquire()
			fmt.Println("gorutina", i, "dela")
			time.Sleep(100 * time.Millisecond)
			sem.Release()
		})
	}
	wg.Wait()
}
```
