package barrier_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/devlights/barrier"
)

func ExampleNewCyclicBarrier() {
	const (
		WORKER_COUNT = 3
	)
	var (
		cb = barrier.NewCyclicBarrier(WORKER_COUNT)
	)
	fmt.Println(cb)

	// Output:
	// CyclicBarrier{Parties: 3, Waiting: 0}
}

func ExampleCyclicBarrier_Await() {
	const (
		WORKER_COUNT = 3
	)
	var (
		cb = barrier.NewCyclicBarrier(WORKER_COUNT)
		wg sync.WaitGroup
	)

	// 3つのワーカーを起動し、全員揃ったら先に進むを繰り返す
	for i := 0; i < WORKER_COUNT; i++ {
		wg.Add(1)
		go worker(i+1, &wg, cb)
	}

	wg.Wait()

	// Unordered Output:
	// Worker-[ 3] 準備作業  1週目
	// Worker-[ 2] 準備作業  1週目
	// Worker-[ 1] 準備作業  1週目
	// Worker-[ 1] 待機開始
	// Worker-[ 2] 待機開始
	// Worker-[ 3] 待機開始
	// Worker-[ 3] 待機解除
	// Worker-[ 3] 準備作業  2週目
	// Worker-[ 2] 待機解除
	// Worker-[ 2] 準備作業  2週目
	// Worker-[ 1] 待機解除
	// Worker-[ 1] 準備作業  2週目
	// Worker-[ 1] 待機開始
	// Worker-[ 2] 待機開始
	// Worker-[ 3] 待機開始
	// Worker-[ 3] 待機解除
	// Worker-[ 3] 準備作業  3週目
	// Worker-[ 2] 待機解除
	// Worker-[ 2] 準備作業  3週目
	// Worker-[ 1] 待機解除
	// Worker-[ 1] 準備作業  3週目
	// Worker-[ 1] 待機開始
	// Worker-[ 2] 待機開始
	// Worker-[ 3] 待機開始
	// Worker-[ 3] 待機解除
	// Worker-[ 2] 待機解除
	// Worker-[ 1] 待機解除
}

func worker(id int, wg *sync.WaitGroup, cb *barrier.CyclicBarrier) {
	defer wg.Done()

	for i := 0; i < 3; i++ {
		fmt.Printf("Worker-[%2d] 準備作業 %2d週目\n", id, i+1)
		time.Sleep(time.Duration(id) * (100 * time.Millisecond))

		fmt.Printf("Worker-[%2d] 待機開始\n", id)
		{
			if err := cb.Await(); err != nil {
				fmt.Printf("Worker-[%2d] エラー: %v\n", id, err)
				return
			}
		}
		fmt.Printf("Worker-[%2d] 待機解除\n", id)
	}
}
