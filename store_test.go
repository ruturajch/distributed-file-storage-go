package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "MOMSBESTPICTURE"
	pathKey := CASPathTransformFunc(key)
	//fmt.Println(pathname)
	expectedOriginalKey := "bfa0a6536ae05408417fd59ce2e3dfd6dac13aba"
	expectedPathName := "bfa0a/6536a/e0540/8417f/d59ce/2e3df/d6dac/13aba"
	if pathKey.PathName != expectedPathName {
		t.Errorf("expected path %s, got %s", expectedPathName, pathKey.PathName)
		return
	}
	if pathKey.Original != expectedOriginalKey {
		t.Errorf("expected original %s, got %s", expectedOriginalKey, pathKey.Original)
		return
	}
}
func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	data := bytes.NewReader([]byte("some jpg bytes"))
	err := s.writeStream("myspecialpicture", data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("test")
}
