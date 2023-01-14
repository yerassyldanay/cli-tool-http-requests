package main

import "fmt"

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func printInt(i *int64) string {
	if i == nil {
		return "[no-resp]"
	}
	return fmt.Sprint(*i)
}
