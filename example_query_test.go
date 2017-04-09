// Copyright 2017 Cristian Greco <cristian@regolo.cc>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collapser_test

import (
	"fmt"
	"time"

	"github.com/TypeTwo/go-collapser"
)

// Consider a server application accessing a database.
// Running the query execution function within a Collapser will alleviate
// database load by reducing many simultaneous requests to a single access.
func Example_databaseQuery() {

	query := "select count(*) from books"

	dbExec := func(q string) int {
		fmt.Println("Query hit database!")
		time.Sleep(1 * time.Second)
		return 42 // After much thought.
	}

	out := make(chan string)

	c := collapser.NewCollapser()

	for i := 0; i < 3; i++ {
		// Launch a goroutine to execute the query.
		go func(i int) {
			res := c.Do(query, func() interface{} {
				return dbExec(query) // Only one query will hit the database.
			})
			out <- fmt.Sprintf("Query #%d: %d", i, res.Get().(int))
		}(i)
	}

	// Wait for all goroutines to complete.
	fmt.Println(<-out)
	fmt.Println(<-out)
	fmt.Println(<-out)

	// Unordered output:
	// Query hit database!
	// Query #0: 42
	// Query #1: 42
	// Query #2: 42
}

