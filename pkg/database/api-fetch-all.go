package database

import (
	"sort"

	"golang.org/x/sync/errgroup"
)

func (db APIDB) FetchAll(ids []string) Vulnerabilities {
	var eg errgroup.Group

	eg.SetLimit(200)

	var osvs Vulnerabilities

	for _, id := range ids {
		eg.Go(func() error {
			// if we error, still report the vulnerability as hopefully the ID should be
			// enough to manually look up the details - in future we should ideally warn
			// the user too, but for now we just silently eat the error
			osv, _ := db.Fetch(id)

			osvs = append(osvs, osv)

			return nil
		})
	}

	// errors are handled within the go routines
	_ = eg.Wait()

	sort.Slice(osvs, func(i, j int) bool {
		return osvs[i].ID < osvs[j].ID
	})

	return osvs
}
