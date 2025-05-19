# cache-go

A simple and efficient cache implementation in Go that implements the S3FIFO (Simple, Scalable, and Stable FIFO) algorithm. The cache is designed to be thread-safe and easy to use.

## Features

- **S3FIFO Algorithm**: Implements the S3FIFO caching algorithm, which provides better hit rates than traditional LRU caches
- **Thread Safety**: All operations are thread-safe using RWMutex and atomic.
- **Simple Design**:
    - Clean and straightforward implementation
    - Easy to understand and maintain
    - Minimal dependencies
- **Generic Implementation**: Written in Go with generics support for type-safe caching of any comparable key type

## Installation

```bash
go get github.com/aryehlev/cache-go
```

## Usage

```go
package main

import (
    "fmt"
	
    "github.com/aryehlev/cache-go"
)

func main() {
    // Create a new cache with size 1000
    cache, err := cache_go.New[string, int](1000)
    if err != nil {
        panic(err)
    }

    // Set a value
    cache.Set("key1", 42)

    // Get a value
    if value, ok := cache.Get("key1"); ok {
        fmt.Printf("Value: %d\n", value)
    }

    // Delete a value
    cache.Delete("key1")

    // Clear the cache
    cache.Clear()
}
```

## How S3FIFO Works

The S3FIFO algorithm divides the cache into three segments:

1. **Small Queue**: For new entries
2. **Main Queue**: For frequently accessed entries
3. **Ghost Queue**: Tracks recently evicted entries

This design provides several benefits:
- Better hit rates than traditional LRU caches
- Simple and effective caching strategy
- Easy to understand and implement
- Good performance for most use cases

visit the [official S3FIFO website](https://s3fifo.com/).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

