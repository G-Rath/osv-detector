package database

import (
	"sort"
)

// a struct to hold the result from each request including an index
// which will be used for sorting the results after they come in
type result struct {
	index int
	res   OSV
	err   error
}

func (db APIDB) FetchAll(ids []string) Vulnerabilities {
	conLimit := 200

	var osvs Vulnerabilities

	if len(ids) == 0 {
		return osvs
	}

	// buffered channel which controls the number of concurrent operations
	semaphoreChan := make(chan struct{}, conLimit)
	resultsChan := make(chan *result)

	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	for i, id := range ids {
		go func(i int, id string) {
			// read from the buffered semaphore channel, which will block if we're
			// already got as many goroutines as our concurrency limit allows
			//
			// when one of those routines finish they'll read from this channel,
			// freeing up a slot to unblock this send
			semaphoreChan <- struct{}{}

			// if we error, still report the vulnerability as hopefully the ID should be
			// enough to manually look up the details - in future we should ideally warn
			// the user too, but for now we just silently eat the error
			osv, _ := db.Fetch(id)
			result := &result{i, osv, nil}

			resultsChan <- result

			// read from the buffered semaphore to free up a slot to allow
			// another goroutine to start, since this one is wrapping up
			<-semaphoreChan
		}(i, id)
	}

	for {
		result := <-resultsChan
		osvs = append(osvs, result.res)

		if len(osvs) == len(ids) {
			break
		}
	}

	sort.Slice(osvs, func(i, j int) bool {
		return osvs[i].ID < osvs[j].ID
	})

	return osvs
}
