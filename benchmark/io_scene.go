package benchmark

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var devNull *os.File
var tmpFile *os.File

func init() {
	f, err := os.Open(os.DevNull)
	if err != nil {
		panic(err)
	}
	devNull = f
	f, err = os.OpenFile(filepath.Join(os.TempDir(), "bench.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	tmpFile = f
}

func Write2DevNull(content string) {
	_, err := devNull.Write([]byte(content))
	if err != nil {
		_ = fmt.Errorf("write failed. %v", err)
	}
}

func Write2TmpFile(content string) {
	_, err := tmpFile.Write([]byte(content))
	if err != nil {
		_ = fmt.Errorf("write failed. %v", err)
	}
}

func LocalPipeReadWrite(content string) {
	r, w := net.Pipe()
	go func() {
		_, err := w.Write([]byte(content))
		if err != nil {
			_ = fmt.Errorf("write failed. %v", err)
			return
		}
		_ = w.Close()
	}()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		_ = fmt.Errorf("write failed. %v", err)
		return
	}
	_ = fmt.Sprintf("data:%s", string(data))
}

func generateLinks(url string) []string {
	urls := make([]string, 10)
	for i := 0; i < 10; i++ {
		urls[i] = url
	}
	return urls
}

func syncCrawl(urls []string) {
	for _, url := range urls {
		http.Get(url)
	}
}

func concurrentCrawl(urls []string) {
	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			http.Get(url)
		}(url)
	}
	wg.Wait()
}

func mockHttpServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond)
		w.WriteHeader(200)
	}))
}
