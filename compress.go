package baidurpc

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// GZIP do gzip action by gzip package
func GZIP(b []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := gzip.NewWriter(buf)
	defer w.Close()

	_, err := w.Write(b)
	w.Flush()

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GUNZIP do unzip action by gzip package
func GUNZIP(b []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write(b)
	r, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	undatas, _ := ioutil.ReadAll(r)

	return undatas, nil
}
