// V2 - Sem deadlock, mas com sincronização excessiva.
// Um "garçom" (um único Mutex para a mesa toda) só deixa um unico filósofo
// pegar garfos e comer por vez. Não há deadlock, pois nunca dois filósofos
// disputam garfos ao mesmo tempo, mas a concorrência cai para 1. Com isso vai demorar mais os 5 rounds de jantar
// usamo mutex
//
//	pra rodar go run ./v2
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	N         = 5                     // filósofos (e garfos)
	CICLOS    = 20                    // ciclos por filósofo
	PENSAR    = 50 * time.Millisecond // tempo pensando
	COMER     = 50 * time.Millisecond // tempo comendo
	EXECUCOES = 5                     // quantas vezes o jantar é repetido
)

var (
	mu          sync.Mutex
	esperaTotal time.Duration
	esperaMax   time.Duration
	numEsperas  int
	comendo     int
	maxComendo  int
	ciclos      [N]int
)

func registrarEspera(d time.Duration) { //registra o tempo de espera para pegar os garfos
	mu.Lock()
	esperaTotal += d
	numEsperas++
	if d > esperaMax {
		esperaMax = d
	}
	mu.Unlock()
}

func comer(id int) {
	mu.Lock()
	comendo++
	if comendo > maxComendo {
		maxComendo = comendo
	}
	mu.Unlock()

	time.Sleep(COMER)

	mu.Lock()
	comendo--
	ciclos[id]++
	mu.Unlock()
}

var garcom sync.Mutex // trava a mesa inteira

func filosofo(id int, garfos []sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	esquerdo, direito := id, (id+1)%N

	for c := 0; c < CICLOS; c++ {
		time.Sleep(PENSAR) // pensar

		inicio := time.Now() // adquirir recursos
		garcom.Lock()        // pede permissão ao garçom
		garfos[esquerdo].Lock()
		garfos[direito].Lock()
		registrarEspera(time.Since(inicio))

		comer(id) // comer

		garfos[direito].Unlock() // liberar recursos
		garfos[esquerdo].Unlock()
		garcom.Unlock()
	}
}

func jantar() time.Duration {
	garfos := make([]sync.Mutex, N)
	var wg sync.WaitGroup

	inicio := time.Now()
	for id := 0; id < N; id++ {
		wg.Add(1)
		go filosofo(id, garfos, &wg)
	}
	wg.Wait()
	return time.Since(inicio)
}

func main() {
	fmt.Println("V2 - Garçom (sincronização excessiva)")
	imprimirCabecalho()

	var somaTotal, somaMedia time.Duration
	for e := 1; e <= EXECUCOES; e++ {
		zerarMetricas()
		total := jantar()
		media := esperaTotal / time.Duration(numEsperas)
		somaTotal += total
		somaMedia += media
		imprimirLinha(fmt.Sprint(e), total, media)
	}
	fmt.Printf("Média | %8.0f ms | %10.2f ms |\n", ms(somaTotal/EXECUCOES), ms(somaMedia/EXECUCOES))
}

func zerarMetricas() {
	esperaTotal, esperaMax, numEsperas, comendo, maxComendo = 0, 0, 0, 0, 0
	ciclos = [N]int{}
}

func imprimirCabecalho() {
	fmt.Printf("%d filósofos, %d ciclos, pensar %v, comer %v\n\n", N, CICLOS, PENSAR, COMER)
	fmt.Println("Exec  | Tempo total | Espera média  | Espera máx  | Máx comendo | Ciclos por filósofo")
}

func imprimirLinha(exec string, total, media time.Duration) {
	fmt.Printf("%-5s | %8.0f ms | %10.2f ms | %8.2f ms | %11d | %v\n",
		exec, ms(total), ms(media), ms(esperaMax), maxComendo, ciclos)
}

func ms(d time.Duration) float64 { return float64(d.Microseconds()) / 1000 }
