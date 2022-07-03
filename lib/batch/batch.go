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

var mx sync.Mutex

func getBatch(n int64, pool int64) (res []user) {
	var i int64
	var wg sync.WaitGroup
	sem := make(chan struct{}, pool)

	for i = 0; i < n; i++ {
		wg.Add(1)
		sem <- struct{}{}
		i := i
		go func() {
			defer wg.Done()
			user := getOne(i)
			<-sem
			mtxSynchronizedAppender(user, &res)
		}()
	}

	wg.Wait()

	return res
}

func mtxSynchronizedAppender(usr user, res *[]user) {
	mx.Lock()
	*res = append(*res, usr)
	mx.Unlock()
}
