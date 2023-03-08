package pkg

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Downloader struct {
	io.Reader
	Total int64
}
type WriteCounter struct {
	Total     uint64
	LastTotal uint64
	LastTime  time.Time
	Speed     uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	elapsed := time.Since(wc.LastTime).Seconds()
	if elapsed >= 1 {
		speed := uint64(float64(wc.Total-wc.LastTotal) / elapsed)
		wc.LastTotal = wc.Total
		wc.LastTime = time.Now()
		wc.Speed = speed
	}
	return n, nil
}

func downloadFile(url string) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	downloader := &Downloader{
		Reader: resp.Body,
		Total:  resp.ContentLength,
	}
	counter := &WriteCounter{}
	if _, err := io.Copy(ioutil.Discard, io.TeeReader(downloader, counter)); err != nil {
		log.Fatalln(err)
	}
}

var wg sync.WaitGroup

func Download() {
	task := []string{}
	task = append(task, "http://cachefly.cachefly.net/100mb.test")
	for _, k := range task {
		wg.Add(1)
		downloadFile(k)
	}
	wg.Wait()
}
