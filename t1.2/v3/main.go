// V3 ordenação na aquisição dos garfos
// todo filósofo pega primeiro o garfo de MENOR número e depois o de MAIOR
// isso impede a espera circular, então não há deadlock, e filósofos
// não vizinhos podem comer ao mesmo tempo
// pra rodar go run ./v3
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

func registrarEspera(d time.Duration) {
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

func filosofo(id int, garfos []sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	primeiro, segundo := id, (id+1)%N
	if primeiro > segundo {
		primeiro, segundo = segundo, primeiro // sempre o menor primeiro
	}

	for c := 0; c < CICLOS; c++ {
		time.Sleep(PENSAR) // pensar

		inicio := time.Now() // adquirir recursos
		garfos[primeiro].Lock()
		garfos[segundo].Lock()
		registrarEspera(time.Since(inicio))

		comer(id) // comer

		garfos[segundo].Unlock() // liberar recursos
		garfos[primeiro].Unlock()
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
	fmt.Println("V3 - Ordenação dos garfos")
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
