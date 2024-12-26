package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashstr := hex.EncodeToString(hash[:])

	blocksize := 5
	slicelength := len(hashstr) / blocksize
	paths := make([]string, slicelength)

	for i := 0; i < slicelength; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashstr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		Original: key,
	}
}

type PathTransformFunc func(path string) PathKey
type PathKey struct {
	PathName string
	Original string
}

func (p PathKey) Filename() string {
	return fmt.Sprintf(p.PathName + "/" + p.Original)
}

type StoreOpts struct {
	PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathName.PathName, os.ModePerm); err != nil {
		return err
	}
	pathOfFile := pathName.Filename()

	f, err := os.Create(pathOfFile)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, pathOfFile)
	return nil
}
