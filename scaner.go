package baiduai

import (
	"bytes"
	"encoding/json"
	"io"
)

// Scaner stores the result in the value pointed to by v
type Scaner interface {
	Scan(v interface{}) error
}

type bytesScaner []byte

func (b bytesScaner) Scan(v interface{}) error {
	return json.Unmarshal(b, v)
}

func (b bytesScaner) RendReader() io.Reader {
	return bytes.NewBuffer(b)
}
