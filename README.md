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
- **Backend Loading Support**: Optional loader function to fetch values from a backend when not found in cache

## Installation

```bash
go get github.com/aryehlev/cache-go
```

## Usage

### Basic Cache

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

### Cache with Backend Loader

```go
package main

import (
    "fmt"
    "database/sql"

    "github.com/aryehlev/cache-go"
)

func main() {
    // Example database connection
    db, err := sql.Open("mysql", "user:password@/dbname")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Create a loader function that fetches values from the database
    loader := func(key string) (int, bool) {
        var value int
        err := db.QueryRow("SELECT value FROM data WHERE key = ?", key).Scan(&value)
        if err != nil {
            return 0, false // Not found or error
        }
        return value, true
    }

    // Create a new cache with size 1000 and the loader function
    cache, err := cache_go.NewWithLoader[string, int](1000, loader)
    if err != nil {
        panic(err)
    }

    // Get a value - will fetch from database if not in cache
    if value, ok := cache.Get("key1"); ok {
        fmt.Printf("Value: %d\n", value)
    }
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
