package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

// Semaphore, arr
// cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
// Benchmark_getButch/#00-8                1000000000              10.91 ns/op
func getBatch(n int64, pool int64) (res []user) {
	var wg sync.WaitGroup
	r := make([]user, n)
	s := make(chan struct{}, pool)
	for i := int64(0); i < n; i++ {
		wg.Add(1)
		s <- struct{}{}
		go func(k int64) {
			r[k] = getOne(k)
			<-s
			wg.Done()
		}(i)
	}
	wg.Wait()
	return r
}

// Semaphore, chan sync1
// cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
// Benchmark_getButch/#00-8                1000000000              10.96 ns/op
// func getBatch(n int64, pool int64) (res []user) {
// 	var wg sync.WaitGroup
// 	r := make(chan user, n)
// 	s := make(chan struct{}, pool)
// 	for i := int64(0); i < n; i++ {
// 		wg.Add(1)
// 		s <- struct{}{}
// 		go func(k int64) {
// 			r <- getOne(k)
// 			<-s
// 			wg.Done()
// 		}(i)
// 	}
// 	wg.Wait()
// 	close(r)
// 	for i := range r {
// 		res = append(res, i)
// 	}
// 	return res
// }

// Semaphore, chan sync2
// cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
// Benchmark_getButch/#00-8                1000000000              11.03 ns/op
// PASS
// func getBatch(n int64, pool int64) (res []user) {
// 	var wg sync.WaitGroup
// 	r := make(chan user)
// 	s := make(chan struct{}, pool)
// 	wg.Add(1)
// 	go func() {
// 		for int64(len(res)) < n {
// 			res = append(res, <-r)
// 		}
// 		close(r)
// 		wg.Done()
// 	}()
// 	for i := int64(0); i < n; i++ {
// 		wg.Add(1)
// 		s <- struct{}{}
// 		go func(k int64) {
// 			r <- getOne(k)
// 			<-s
// 			wg.Done()
// 		}(i)
// 	}
// 	wg.Wait()
// 	return res
// }

// Semaphore, chan sync3
// cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
// Benchmark_getButch/#00-8                1000000000              11.02 ns/op
// PASS
// func getBatch(n int64, pool int64) (res []user) {
// 	var wgr, wgw sync.WaitGroup
// 	rr := make([]user, n)
// 	r := make(chan user)
// 	s := make(chan struct{}, pool)
// 	quit := make(chan struct{})
// 	wgr.Add(1)
// 	go func() {
// 		i := 0
// 		for {
// 			select {
// 			case rr[i] = <-r:
// 				i++
// 			case <-quit:
// 				wgr.Done()
// 			}
// 		}
// 	}()
// 	for i := int64(0); i < n; i++ {
// 		wgw.Add(1)
// 		s <- struct{}{}
// 		go func(k int64) {
// 			r <- getOne(k)
// 			<-s
// 			wgw.Done()
// 		}(i)
// 	}
// 	wgw.Wait()
// 	quit <- struct{}{}
// 	wgr.Wait()
// 	return rr
// }

// ErrGroup, chan sync
// cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
// Benchmark_getButch/#00-8                1000000000              10.92 ns/op
// PASS
// func getBatch(n int64, pool int64) (res []user) {
// 	errG, _ := errgroup.WithContext(context.Background())
// 	r := make(chan user, n)
// 	errG.SetLimit(int(pool))
// 	for i := int64(0); i < n; i++ {
// 		k := i
// 		errG.Go(func() error {
// 			r <- getOne(k)
// 			return nil
// 		})
// 	}
// 	err := errG.Wait()
// 	if err != nil {
// 		panic(err)
// 	}
// 	close(r)
// 	for i := range r {
// 		res = append(res, i)
// 	}
// 	return res
// }

// Worker group
// cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
// Benchmark_getButch/#00-8                1000000000              11.05 ns/op
// PASS
// func getBatch(n int64, pool int64) (res []user) {
// 	r := make(chan user, n)
// 	j := make(chan int64, n)
// 	for i := int64(0); i < pool; i++ {
// 		go func(j <-chan int64, r chan<- user) {
// 			for x := range j {
// 				r <- getOne(x)
// 			}
// 		}(j, r)
// 	}

// 	for i := int64(0); i < n; i++ {
// 		j <- i
// 	}

// 	close(j)

// 	for i := int64(0); i < n; i++ {
// 		res = append(res, <-r)
// 	}

// 	return res
// }
