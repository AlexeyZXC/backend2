package l5

import (
	"fmt"
	"sync"
)

func f1() {
	var (
		done = make(chan struct{})
		res  = make(chan int)
	)

	go func() {
		for userId := range res {
			fmt.Println(userId)
		}
		close(done)
	}()

	wg := &sync.WaitGroup{}

	for _, shard := range shards {
		wg.Add(1)
		go func(shard *sharding.Shard) {
			defer wg.Done()
			rows, err := shard.Query(`SELECT "user_id" FROM "users" WHERE "age" BETWEEN	30 AND 39`)
			if err != nil {
				//@todo process error
				return
			}
			defer rows.Close()
			for rows.Next() {
				var userId int
				if err := rows.Scan(&userId); err != nil {
					//@todo process error
					return
				}
				res <- userId
			}
		}(shard)
	}
	wg.Wait()
	close(res)
	<-done
}
