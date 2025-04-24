package main

import (
	"fmt"
	"s3fifo/s3fifo"
)

func main() {
	cache := s3fifo.New[string, string](2)

	cache.Set("key", "value")
	cache.Set("key", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key4", "value4")
	cache.Set("key45", "value45")

	fmt.Println(cache.Get("key"))
}
