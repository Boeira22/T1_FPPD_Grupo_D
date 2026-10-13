// V0 - Versão sequencial, código do moodle. Os filósofos executam um após o outro, sem goroutines. Não há disputa
// pelos garfos (espera = 0) e só um filósofo come por vez, peguei o código do moodle e acrescentei a coleta de métricas
// (espera total, espera máxima, número de esperas, máximo comendo juntos e ciclos por filósofo) e a saida bonita ali dos dados
// pra roda  go run ./v0   
package main

import (
	"fmt"
	"sync" // serve para proteger os contadores de métricas
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
	mu          sync.Mutex // <- protege os contadores
	esperaTotal time.Duration
	esperaMax   time.Duration
	numEsperas  int
	comendo     int
	maxComendo  int
	ciclos      [N]int // contador de ciclos por filósofo
)

func registrarEspera(d time.Duration) {
	mu.Lock()
	esperaTotal += d // atualiza a soma das esperas
	numEsperas++
	if d > esperaMax { // atualiza o máximo de espera
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

func jantar() time.Duration {
	inicio := time.Now()
	for c := 0; c < CICLOS; c++ {
		for id := 0; id < N; id++ {
			time.Sleep(PENSAR) // pensar
			registrarEspera(0) // adquirir: garfos sempre livres
			comer(id)          // comer
		} // liberar: nada a fazer
	}
	return time.Since(inicio) // retorna o tempo total do jantar
}

func main() {
	fmt.Println("V0 - Sequencial")
	imprimirCabecalho()

	var somaTotal, somaMedia time.Duration
	for e := 1; e <= EXECUCOES; e++ { // vai ser sempre cinco execuções
		zerarMetricas()
		total := jantar()
		media := esperaTotal / time.Duration(numEsperas) // calcula a média de espera pega a espera total e divide pelo número de esperas em ms
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
	fmt.Println("Exec  | Temp total | Espera média  | Espera máx  | Máx comendo | Ciclo por filosofos")
}

func imprimirLinha(exec string, total, media time.Duration) {
	fmt.Printf("%-5s | %8.0f ms | %10.2f ms | %8.2f ms | %11d | %v\n",
		exec, ms(total), ms(media), ms(esperaMax), maxComendo, ciclos)
}

func ms(d time.Duration) float64 { return float64(d.Microseconds()) / 1000 } // converte de nanossegundos para milissegundos

//   //think
//   for s := 1; s <= id*6; s++ {
//          fmt.Printf(" ")
//          }
//     fmt.Printf(" T%d\n", iteration) // Thinking
//     time.Sleep(time.Duration(time_to_think) * time.Second)
//
//     // eat
//
//     for s := 1; s <= id*6; s++ {
//          fmt.Printf(" ")
//          }
//     fmt.Printf(" E%d\n", iteration) // Eating
//     time.Sleep(time.Duration(time_to_eat) * time.Second)
//  }
