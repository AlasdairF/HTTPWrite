package httpwrite

import (
	"github.com/AlasdairF/Pool"
	"github.com/AlasdairF/Conv"
	"net/http"
	"compress/gzip"
	"reflect"	
)

const (
	bestLength = pool.Size
)

type ResponseWriter interface {
	Header() http.Header
	WriteHeader(int) 
	Write([]byte) (int, error)
	WriteString(string) (int, error)
	WriteByte(byte) error
	WriteAll(...interface{})
	Close() error
}

type Plain struct {
	ResponseWriter http.ResponseWriter
	data []byte
	cursor int
}

type Gzip struct {
	ResponseWriter http.ResponseWriter
	gz *gzip.Writer
	data []byte
	cursor int
}

func New(w http.ResponseWriter) *Plain {
	return &Plain{ResponseWriter: w, data: pool.Get(bufferLen)}
}

func (b *Plain) Header() http.Header {
	return b.ResponseWriter.Header()
}

func (b *Plain) WriteHeader(header int) {
	b.ResponseWriter.WriteHeader(header)
}

func (b *Plain) Write(p []byte) (int, error) {
	l := len(p)
	if b.cursor + l > bestLength {
		var err error
		if b.cursor > 0 {
			_, err = b.ResponseWriter.Write(b.data[0:b.cursor]) // flush
		}
		if l > bestLength { // data to write is longer than the length of the Plain
			b.cursor = 0
			return b.ResponseWriter.Write(p)
		}
		copy(b.data[0:l], p)
		b.cursor = l
		return l, err
	}
	copy(b.data[b.cursor:], p)
	b.cursor += l
	return l, nil
}

func (b *Plain) WriteString(p string) (int, error) {
	l := len(p)
	if b.cursor + l > bestLength {
		var err error
		if b.cursor > 0 {
			_, err = b.ResponseWriter.Write(b.data[0:b.cursor]) // flush
		}
		if l > bestLength { // data to write is longer than the length of the Plain
			b.cursor = 0
			return b.ResponseWriter.Write([]byte(p))
		}
		copy(b.data[0:l], p)
		b.cursor = l
		return l, err
	}
	copy(b.data[b.cursor:], p)
	b.cursor += l
	return l, nil
}

func (b *Plain) WriteByte(p byte) error {
	if b.cursor < bestLength {
		b.data[b.cursor] = p
		b.cursor++
		return nil
	}
	var err error
	if b.cursor > 0 {
		_, err = b.ResponseWriter.Write(b.data[0:b.cursor]) // flush
	}
	b.data[0] = p
	b.cursor = 1
	return err
}

func (b *Plain) WriteAll(a ...interface{}) {
	for _, p := range a {
		switch reflect.TypeOf(p).Kind() {
			case reflect.String:
				b.WriteString(reflect.ValueOf(p).String())
			case reflect.Slice:
				b.Write(reflect.ValueOf(p).Bytes())
			case reflect.Int: case reflect.Int8: case reflect.Int16: case reflect.Int32: case reflect.Int64:
				b.Write(conv.Bytes(int(reflect.ValueOf(p).Int())))
			case reflect.Uint: case reflect.Uint8: case reflect.Uint16: case reflect.Uint32: case reflect.Uint64:
				b.Write(conv.Bytes(int(reflect.ValueOf(p).Uint())))
		}
	}
}

func (b *Plain) Close() (err error) {
	if b.cursor > 0 {
		_, err = b.ResponseWriter.Write(b.data[0:b.cursor])
		b.cursor = 0
	}
	b.ResponseWriter = nil
	pool.Return(b.data)
	return
}

func NewGzip(w http.ResponseWriter) *Gzip {
	return &Gzip{ResponseWriter: w, gz: gzip.NewWriter(w), data: pool.Get(bufferLen)}
}

func (b *Gzip) Header() http.Header {
	return b.ResponseWriter.Header()
}

func (b *Gzip) WriteHeader(header int) {
	b.ResponseWriter.WriteHeader(header)
}

func (b *Gzip) Write(p []byte) (int, error) {
	l := len(p)
	if b.cursor + l > bestLength {
		var err error
		if b.cursor > 0 {
			_, err = b.gz.Write(b.data[0:b.cursor]) // flush
		}
		if l > bestLength { // data to write is longer than the length of the Gzip
			b.cursor = 0
			return b.gz.Write(p)
		}
		copy(b.data[0:l], p)
		b.cursor = l
		return l, err
	}
	copy(b.data[b.cursor:], p)
	b.cursor += l
	return l, nil
}

func (b *Gzip) WriteString(p string) (int, error) {
	l := len(p)
	if b.cursor + l > bestLength {
		var err error
		if b.cursor > 0 {
			_, err = b.gz.Write(b.data[0:b.cursor]) // flush
		}
		if l > bestLength { // data to write is longer than the length of the Gzip
			b.cursor = 0
			return b.gz.Write([]byte(p))
		}
		copy(b.data[0:l], p)
		b.cursor = l
		return l, err
	}
	copy(b.data[b.cursor:], p)
	b.cursor += l
	return l, nil
}

func (b *Gzip) WriteByte(p byte) error {
	if b.cursor < bestLength {
		b.data[b.cursor] = p
		b.cursor++
		return nil
	}
	var err error
	if b.cursor > 0 {
		_, err = b.gz.Write(b.data[0:b.cursor]) // flush
	}
	b.data[0] = p
	b.cursor = 1
	return err
}

func (b *Gzip) WriteAll(a ...interface{}) {
	for _, p := range a {
		switch reflect.TypeOf(p).Kind() {
			case reflect.String:
				b.WriteString(reflect.ValueOf(p).String())
			case reflect.Slice:
				b.Write(reflect.ValueOf(p).Bytes())
			case reflect.Int: case reflect.Int8: case reflect.Int16: case reflect.Int32: case reflect.Int64:
				b.Write(conv.Bytes(int(reflect.ValueOf(p).Int())))
			case reflect.Uint: case reflect.Uint8: case reflect.Uint16: case reflect.Uint32: case reflect.Uint64:
				b.Write(conv.Bytes(int(reflect.ValueOf(p).Uint())))
		}
	}
}

func (b *Gzip) Close() (err error) {
	if b.cursor > 0 {
		_, err = b.gz.Write(b.data[0:b.cursor])
		b.cursor = 0
	}
	b.gz.Close()
	b.gz = nil
	b.ResponseWriter = nil
	pool.Return(b.data)
	return
}
