package pkg

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

var (
	totalElapsed int64
	totalBytes   int64
	mu           sync.Mutex
)

type Downloader struct {
	io.Reader
	Total int64
}
type WriteCounter struct {
	Total     int64
	LastTotal int64
	LastTime  time.Time
	elapsed   int64
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

func downloadFile(url string, i int, p *mpb.Progress, proxy string) {
	defer wg.Done()
	var resp *http.Response
	var err error
	if proxy != "" {
		resp, err = ProxyGet(url, proxy)
	} else {
		resp, err = http.Get(url)
	}
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

	// progress bar init
	name := fmt.Sprintf("🔗 Thread-%d ", i)
	bar := p.AddBar(downloader.Total,
		mpb.PrependDecorators(
			// simple name decorator
			decor.Name(name),
			// decor.DSyncWidth bit enables column width synchronization
			decor.CountersNoUnit("%d/%d Bytes", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(

			decor.Any(func(s decor.Statistics) string {
				return fmt.Sprintf("Speed: %.2f MB/s", counter.Speed)
			}, decor.WCSyncSpaceR),
			decor.OnComplete(decor.Percentage(decor.WC{W: 5}), "done"),
		),
	)
	start := time.Now()
	// echo speed per second
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	go func() {
		for {
			bar.SetCurrent(counter.Total)
		}
	}()

	// main function

	if _, err := io.Copy(ioutil.Discard, io.TeeReader(downloader, counter)); err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start).Milliseconds()
	mu.Lock()
	if totalElapsed < elapsed {
		totalElapsed = elapsed
	}
	totalBytes += downloader.Total
	mu.Unlock()
}

var wg sync.WaitGroup

func Download(threadsCount int, testFileSize string, proxy string) {
	p := mpb.New(mpb.WithWaitGroup(&wg))
	var task string
	switch testFileSize {
	case "1":
		task = "http://cachefly.cachefly.net/1mb.test"
	case "5":
		task = "http://cachefly.cachefly.net/5mb.test"
	case "10":
		task = "http://cachefly.cachefly.net/10mb.test"
	case "100":
		task = "http://cachefly.cachefly.net/100mb.test"
	default:
		task = "http://cachefly.cachefly.net/10mb.test"
	}

	log.Println("🍌 Speedtest Initialized")
	for i := 0; i < threadsCount; i++ {
		wg.Add(1)
		go downloadFile(task, i, p, proxy)
	}
	p.Wait()
	log.Printf("🍌 Speedtest Finished in %d ms. Average Speed %.2f MB/s", totalElapsed, float64(totalBytes)/float64(totalElapsed)/1024/1024*1000)
}
