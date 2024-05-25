package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key)) // [20]byte ->  [:]
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		Filename: hashStr,
	}

}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	Filename string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
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

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()
	return os.RemoveAll(pathKey.FullPath())
	//err != nil {
	//	return err
	//}
	// fi,err := os.Stat
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	// io.Copy(buf, f)
	//f.Close()
	return buf, err

}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)

	return os.Open(pathKey.FullPath())

}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	//	buf := new(bytes.Buffer)
	//	io.Copy(buf, r)
	//
	//	filenameBytes := md5.Sum(buf.Bytes())
	//	filename := hex.EncodeToString(filenameBytes[:])
	//	pathAndFilename := pathKey.PathName + "/" + filename

	//	f , err := os.Create(pathAndFilename)
	//	if err != nil {
	//		return err
	//	}
	//
	//	filename := "somefilename"
	//
	//	pathAndFilename := pathName + "/" + filename

	// Create the file
	//pathAndFilename := pathKey.Filename()

	fullPath := pathKey.FullPath()

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	// Copy data from the reader to the file
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("Written {%d} bytes to disk: %s", n, fullPath)
	return nil
}
