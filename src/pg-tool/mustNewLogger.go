package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func mustNewLogger(filePath string) (err error) {

	dir := filepath.Dir(filePath)
	if err = os.MkdirAll(dir, 0766); err != nil {
		log.Fatalf("error create dir [%s] for logger err: [%v]", dir, err)
	}

	if f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		log.Fatalf("error open file [%s] for logger err: [%v]", filePath, err)
	} else {
		w := io.MultiWriter(os.Stdout, f)
		log.SetOutput(w)
	}

	log.SetFlags(log.LstdFlags)

	return
}
