// V1 - Versão concorrente Deadlock!!.
// Cada garfo é um Mutex e TODO filósofo pega primeiro o garfo da esquerda
// e depois o da direita. Se os 5 pegarem o esquerdo ao mesmo tempo, cada um
// espera o garfo que está com o vizinho: espera circular -> deadlock.
// Quando todas as goroutines ficam bloqueadas, o próprio Go encerra o
// programa com: "fatal error: all goroutines are asleep - deadlock!"
// pra roda go run ./v1
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	N      = 5
	CICLOS = 20
	PENSAR = 50 * time.Millisecond
	COMER  = 50 * time.Millisecond

	// coloquei uma pauda de 10 ms entre pegar o garfo da esquerda e o da direita, pra ter certeza do deadlock em 0 tem chance de nao acontecer
	//isso acontece pois o tempo de execução das goroutines é imprevisível, então pode acontecer de um filósofo pegar o garfo da esquerda e outro pegar o da direita
	// antes do primeiro tentar pegar o da direita, evitando o deadlock
	ATRASO = 10 * time.Millisecond
)

func filosofo(id int, garfos []sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	esquerdo, direito := id, (id+1)%N

	for c := 1; c <= CICLOS; c++ {
		time.Sleep(PENSAR) // pensar

		garfos[esquerdo].Lock() // adquirir recursos
		fmt.Printf("P%d pegou o garfo %d e espera o garfo %d\n", id, esquerdo, direito)
		time.Sleep(ATRASO)
		garfos[direito].Lock()

		time.Sleep(COMER) // comer

		garfos[direito].Unlock() // liberar recursos
		garfos[esquerdo].Unlock()
	}
}

func main() {
	fmt.Println("V1 - Concorrente e deadlock")
	fmt.Printf("%d filósofos, %d ciclos, pensar %v, comer %v\n\n", N, CICLOS, PENSAR, COMER)
	garfos := make([]sync.Mutex, N) // cada garfo é um Mutex
	var wg sync.WaitGroup

	inicio := time.Now()
	for id := 0; id < N; id++ {
		wg.Add(1)
		go filosofo(id, garfos, &wg)
	}
	wg.Wait()

	fmt.Printf("\nTerminou sem deadlock em %v\n", time.Since(inicio)) //caso o programa incrivelmente consiga terminar sem o deadlock vai sair o tempo de exec
}
