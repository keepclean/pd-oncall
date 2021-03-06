package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jedib0t/go-pretty/text"
)

type FileAsset struct {
	Path string
}

func (f FileAsset) String() string {
	return f.Path
}

func (f FileAsset) ExpandPath() string {
	return os.ExpandEnv(f.Path)
}

func (f FileAsset) DirName() string {
	return filepath.Dir(f.ExpandPath())
}

func (f FileAsset) Exist() bool {
	if _, err := os.Stat(f.ExpandPath()); err != nil && os.IsNotExist(err) {
		log.Printf("file %q doesn't exist\n", f.ExpandPath())
		return false
	} else if err != nil {
		log.Println("non-IsNotExist error upon calling os.Stat:", err)
		return false
	}

	return true
}

func (f FileAsset) Remove() error {
	if err := os.Remove(f.ExpandPath()); err != nil {
		return fmt.Errorf("can not remove %q file: %v", f, err)
	}

	log.Printf("file %q has been removed\n", f.ExpandPath())
	return nil
}

func (f FileAsset) CreateDirs() error {
	if err := os.MkdirAll(f.DirName(), 0755); err != nil {
		return fmt.Errorf("can't create directory chain for %q file: %v", f.ExpandPath(), err)
	}

	return nil
}

func (f FileAsset) Read() (*Schedules, error) {
	fd, err := os.Open(f.ExpandPath())
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var s *Schedules
	if err = json.NewDecoder(fd).Decode(&s); err != nil {
		return nil, err
	}

	return s, nil
}

func (f FileAsset) Show() {
	data, err := f.Read()
	if err != nil {
		log.Fatalln("can not read the cache file", f.ExpandPath(), err)
	}

	jsonPrettyPrinter := text.NewJSONTransformer("", "  ")
	fmt.Println(jsonPrettyPrinter(data))
}

func (f FileAsset) Write(s *Schedules) error {
	fd, err := os.Create(f.ExpandPath())
	if err != nil {
		return fmt.Errorf("can't create file: %v", err)
	}
	defer fd.Close()

	if err = json.NewEncoder(fd).Encode(s); err != nil {
		return fmt.Errorf("can't write json: %v", err)
	}

	return nil
}
