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

var ()

type Downloader struct {
	io.Reader
	Total int64
}
type WriteCounter struct {
	Total     int64
	LastTotal int64
	LastTime  time.Time
	Speed     float64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += int64(n)
	elapsed := float64(time.Since(wc.LastTime).Milliseconds()) / 1000.0
	if elapsed >= 0.5 {
		speed := float64(wc.Total-wc.LastTotal) / elapsed / 1024 / 1024
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

	// echo speed per second
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	go func() {
		for {
			fmt.Printf("\r👀 Download Speed: %.3f MB/s", counter.Speed)
			if counter.Total == downloader.Total {
				return
			}
		}
	}()

	// main function
	start := time.Now()
	if _, err := io.Copy(ioutil.Discard, io.TeeReader(downloader, counter)); err != nil {
		log.Fatalln(err)
	}
	elapsed := time.Since(start).Milliseconds()
	fmt.Print("\r")
	log.Printf("🍌 Speedtest Finished\t[Speed=%.2f MB/s]", float64(downloader.Total)/float64(elapsed)/1024.0/1024.0*1000.0)
}

var wg sync.WaitGroup

func Download() {
	task := []string{}
	task = append(task, "http://cachefly.cachefly.net/10mb.test")
	task = append(task, "http://cachefly.cachefly.net/10mb.test")
	log.Println("🍌 Speedtest Initialized")
	for _, k := range task {
		wg.Add(1)
		go downloadFile(k)
	}
	wg.Wait()
}
