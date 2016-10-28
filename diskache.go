package diskache

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/GitbookIO/syncgroup"
)

type Diskache struct {
	directory string
	items     int
	lock      *syncgroup.MutexGroup
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
	if err := os.MkdirAll(opts.Directory, os.ModePerm); err != nil {
		return nil, err
	}

	// Create Diskache instance
	dc := &Diskache{
		directory: opts.Directory,
		lock:      syncgroup.NewMutexGroup(),
	}

	return dc, nil
}

func (dc *Diskache) Set(key string, data []byte) error {
	// Get encoded key
	filename := dc.buildFilename(key)

	// Lock for writing
	dc.lock.Lock(filename)
	defer dc.lock.Unlock(filename)

	// Open file
	file, err := os.Create(filename)
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
	// Get encoded key
	filename := dc.buildFilename(key)

	// Lock for reading
	dc.lock.RLock(filename)
	defer dc.lock.RUnlock(filename)

	// Open file
	file, err := os.Open(filename)
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
