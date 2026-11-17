# Primeri programerskih nalog

## Naloga 1

V jeziku Go želimo napisati program, v katerem glavna gorutina zažene *N* dodatnih gorutin, ki se izvajajo sočasno. Vsaka dodatna gorutina izpiše svojo oznako in zaključi. Program popravite tako, da se bodo vse oznake zagotovo pojavile na standardnem izhodu. Za sinhronizacijo gorutin uporabite kanal `ch`.

```go
package main

import "fmt"

const N int = 5

var ch chan bool = make(chan bool)

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
