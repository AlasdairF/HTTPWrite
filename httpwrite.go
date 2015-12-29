package httpwrite

import (
	"github.com/AlasdairF/Custom"
	"github.com/AlasdairF/Conv"
	"net/http"
	"compress/gzip"
	"reflect"	
)

type Response struct {
	ResponseWriter http.ResponseWriter
	Writer *custom.Writer
}

func New(w http.ResponseWriter) Response {
	return Response{ResponseWriter: w, Writer: custom.NewWriter(w)}
}

func NewGzip(w http.ResponseWriter) Response {
	return Response{ResponseWriter: w, Writer: custom.NewWriter(gzip.NewWriter(w))}
}

func (b Response) WriteAll(a ...interface{}) {
	for _, p := range a {
		switch reflect.TypeOf(p).Kind() {
			case reflect.String:
				b.Writer.WriteString(reflect.ValueOf(p).String())
			case reflect.Slice: // all slices are assumed to be slices of bytes
				b.Writer.Write(reflect.ValueOf(p).Bytes())
			case reflect.Uint8: // byte
				b.Writer.WriteByte(byte(reflect.ValueOf(p).Uint()))
			case reflect.Int: case reflect.Int8: case reflect.Int16: case reflect.Int32: case reflect.Int64:
				conv.Write(b.Writer, int(reflect.ValueOf(p).Int()), 0)
			case reflect.Uint: case reflect.Uint16: case reflect.Uint32: case reflect.Uint64:
				conv.Write(b.Writer, int(reflect.ValueOf(p).Uint()), 0)
		}
	}
}

func (b Response) Close() (err error) {
	return b.Writer.Close()
}

func (b Response) Header() http.Header {
	return b.ResponseWriter.Header()
}

func (b Response) WriteHeader(header int) {
	b.ResponseWriter.WriteHeader(header)
}
