package nasa

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type respWriter struct {
	statusCode int
	out        io.Writer
	w          http.ResponseWriter
	buf        *bytes.Buffer
}

func newRespWriter(w http.ResponseWriter) *respWriter {
	buf := bytes.NewBuffer(nil)
	out := io.MultiWriter(w, buf)
	return &respWriter{
		out: out,
		w:   w,
		buf: buf,
	}
}

func (rw *respWriter) Header() http.Header {
	return rw.w.Header()
}

func (rw *respWriter) Write(b []byte) (int, error) {
	return rw.out.Write(b)
}

func (rw *respWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.w.WriteHeader(statusCode)
}

func CacheHandler(store *Store, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	// there is a risk of having multiple threads trying to download at the same time to store the same time
	return func(w http.ResponseWriter, r *http.Request) {

		item := time.Now().Format("20060102")

		// check if the image is already in the
		b, found, err := store.Get(item)
		if err != nil {
			fmt.Printf("failed to find item %s, %+v\n", item, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if found {
			fmt.Printf("found item %s in cache\n", item)
			w.Header().Set("Content-Type", "image/png")
			w.WriteHeader(http.StatusOK)
			if _, err = w.Write(b); err != nil {
				fmt.Printf("failed to send response to client, %+v\n", err)
			}
			return
		}

		fmt.Printf("item %s not found in cache\n", item)

		cw := newRespWriter(w)

		// item not cached, get it

		handler(cw, r)

		if cw.statusCode >= 200 && cw.statusCode <= 300 {
			b = cw.buf.Bytes()
			if err := store.Put(item, b); err != nil {
				fmt.Printf("failed to store iteam %s in cache, %+v\n", item, err)
			}
		}
	}
}

type Store struct {
	sync.Mutex
	Path string
}

func (s *Store) Get(key string) (b []byte, found bool, err error) {
	s.Lock()
	defer s.Unlock()

	path := filepath.Join(s.Path, key)
	b, err = os.ReadFile(path)
	if err == nil {
		found = true
	}
	if os.IsNotExist(err) {
		err = nil
		found = false
		return
	}
	found = true
	return
}

func (s *Store) Put(key string, value []byte) error {
	s.Lock()
	defer s.Unlock()

	path := filepath.Join(s.Path, key)
	return os.WriteFile(path, value, 0666)
}
