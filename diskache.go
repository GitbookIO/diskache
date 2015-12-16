package diskache

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
)

type Diskache struct {
	directory string
	items     int
	sync.RWMutex
}

type Opts struct {
	Directory string
}

type Stats struct {
	Directory string
	Items     int
}

func New(opts *Opts) (*Diskache, error) {
	// Create Diskache directory
	err := os.MkdirAll(opts.Directory, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Create Diskache instance
	dc := &Diskache{}
	dc.directory = opts.Directory

	return dc, nil
}

func (dc *Diskache) Set(key string, data []byte) error {
	// Lock for writing
	dc.Lock()
	defer dc.Unlock()

	// Open file
	file, err := os.Create(dc.buildFilename(key))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write data
	if _, err = file.Write(data); err == nil {
		// Increment items
		dc.items += 1
	}

	return err
}

func (dc *Diskache) Get(key string) ([]byte, bool) {
	// Lock for reading
	dc.RLock()
	defer dc.RUnlock()

	// Open file
	file, err := os.Open(dc.buildFilename(key))
	if err != nil {
		return nil, false
	}
	defer file.Close()

	// Read file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Diskache: Error reading from file %s\n", key)
		return nil, false
	}

	return data, true
}

func (dc *Diskache) Clean() error {
	// Delete directory
	if err := os.RemoveAll(dc.directory); err != nil {
		return err
	}
	// Recreate directory
	return os.MkdirAll(dc.directory, os.ModePerm)
}

func (dc *Diskache) Stats() Stats {
	return Stats{
		Directory: dc.directory,
		Items:     dc.items,
	}
}

func (dc *Diskache) buildFilename(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return path.Join(dc.directory, hex.EncodeToString(hasher.Sum(nil)))
}
