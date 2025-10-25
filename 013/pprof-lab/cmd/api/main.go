package main

import (
	"example.com/pprof-lab/internal/work"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
)

var (
	globalMutex sync.Mutex
	sharedData  int
)

func main() {
	// Включаем профилирование блокировок и мьютексов
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	log.Println("Block and mutex profiling enabled")

	// Медленная версия
	http.HandleFunc("/work-slow", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("Fib(38)")()
		res := work.Fib(38)
		w.Write([]byte(fmt.Sprintf("Slow: %d\n", res)))
	})

	// Быстрая версия  
	http.HandleFunc("/work-fast", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("FibFast(38)")()
		res := work.FibFast(38)
		w.Write([]byte(fmt.Sprintf("Fast: %d\n", res)))
	})

	// Демо блокировок
	http.HandleFunc("/block-demo", func(w http.ResponseWriter, r *http.Request) {
		// Блокировка каналом
		ch := make(chan int)
		go func() {
			ch <- 42
		}()
		result := <-ch

		// Конкуренция за мьютекс
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				globalMutex.Lock()
				sharedData++
				globalMutex.Unlock()
			}()
		}
		wg.Wait()

		w.Write([]byte(fmt.Sprintf("Result: %d, Shared: %d\n", result, sharedData)))
	})

	log.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}