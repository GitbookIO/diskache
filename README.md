# diskache

Lightweight Golang disk cache.

## Get

```Shell
$ go get github.com/GitbookIO/diskache
```

## Use

```Go
import (
    "fmt"
    "github.com/GitbookIO/diskache"
)

// Create an instance
opts := diskache.Opts{
    Directory: "diskache_place",
}
dc := diskache.New(opts)

// Add data to cache
spelling := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
err := dc.Set("spelling", spelling)
if err != nil {
    fmt.Println("Impossible to set data in cache")
}

// Read from cache
cached, inCache := dc.Get("spelling")
if inCache {
    fmt.Println(string(cached))
}

// Read stats
stats := dc.Stats()
reflect.DeepEqual(stats, Stats{
    Directory: "diskache_place",
    Items:     1,
})

// Cleanup
err = dc.Clean()
if err != nil {
    fmt.Println("Impossible to clean cache")
}
```
