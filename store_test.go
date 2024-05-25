package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "applePicture"
	pathKey := CASPathTransformFunc(key)

	expectedOriginalKey := "bfb934a501cf493928efd64d8907ba7b1ff97249"
	expectedPathName := "bfb93/4a501/cf493/928ef/d64d8/907ba/7b1ff/97249"

	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s want %s", pathKey.PathName, expectedPathName)
	}
	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("have %s want %s", pathKey.Filename, expectedOriginalKey)
	}
}

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "applePicture"
	data := []byte("some jpg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	// Define test options
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	// Create a new Store instance
	s := NewStore(opts)

	key := "applePicture"

	// Create a bytes.Reader with some test data
	data := []byte("Some JPG bytes")

	// Call the writeStream method of the Store
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		// If there is an error, fail the test
		t.Error(err)
	}
	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := ioutil.ReadAll(r)

	fmt.Println(string(b))

	if string(b) != string(data) {
		t.Errorf("want %s have %s", data, b)
	}

}
