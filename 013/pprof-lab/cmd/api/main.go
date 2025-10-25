package main

import (
	"example.com/pprof-lab/internal/work"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

var (
	globalMutex sync.Mutex
	sharedData  int
)

func enableProfiling() {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	log.Println("Block and mutex profiling enabled")
}

func main() {
	enableProfiling()

	// Старая медленная версия
	http.HandleFunc("/work-slow", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("Fib(38)")()
		n := 38
		res := work.Fib(n)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(fmt.Sprintf("Slow: %d\n", res)))
	})

	// Новая быстрая версия
	http.HandleFunc("/work-fast", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("FibFast(38)")()
		n := 38
		res := work.FibFast(n)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(fmt.Sprintf("Fast: %d\n", res)))
	})

	// Эндпоинт для демонстрации блокировок
	http.HandleFunc("/block-demo", func(w http.ResponseWriter, r *http.Request) {
		// Демонстрация блокировок на каналах
		ch := make(chan int)

		go func() {
			time.Sleep(100 * time.Millisecond)
			ch <- 42
		}()

		result := <-ch

		// Демонстрация конкуренции за мьютекс
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				globalMutex.Lock()
				sharedData++
				time.Sleep(10 * time.Millisecond)
				globalMutex.Unlock()
			}(i)
		}
		wg.Wait()

		w.Write([]byte(fmt.Sprintf("Result: %d, Shared: %d\n", result, sharedData)))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is working! Endpoints: /work-slow, /work-fast, /block-demo, /debug/pprof/"))
	})

	log.Println("Server on :8080; pprof on /debug/pprof/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
