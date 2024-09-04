package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second

	timeout = 100 * time.Millisecond
)

func main() {
	rmt := &runtime.MemStats{}

	go func() {
		for {
			runtime.ReadMemStats(rmt)
			time.Sleep(pollInterval)
		}
	}()

	client := &http.Client{}
	client.Timeout = timeout

	log.Println("Starting agent")

	for {
		time.Sleep(reportInterval)

		url := fmt.Sprintf("http://127.0.0.1:8080/update/%s/%s/%d", "gauge", "Alloc", rmt.Alloc)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			log.Printf("Error while creating request to '%s', error %v", url, err)
			continue
		}

		response, err := client.Do(request)
		if err != nil {
			log.Printf("Error while sending request to '%s', error %v", url, err)
			continue
		}

		log.Printf("Sent request to %s, status %d", url, response.StatusCode)
	}

	//runtime.ReadMemStats(rmt)
	//
	//fmt.Println(rmt.Alloc)
	//fmt.Println(rmt.TotalAlloc)
	//fmt.Println(rmt.BuckHashSys)
	//fmt.Println(rmt.Frees)
	//fmt.Println(rmt.GCCPUFraction)
	//fmt.Println(rmt.GCSys)
	//fmt.Println(rmt.HeapAlloc)
	//fmt.Println(rmt.HeapIdle)
	//fmt.Println(rmt.HeapInuse)
	//fmt.Println(rmt.HeapObjects)
	//fmt.Println(rmt.HeapReleased)
	//fmt.Println(rmt.HeapSys)
	//fmt.Println(rmt.LastGC)
	//fmt.Println(rmt.Lookups)
	//fmt.Println(rmt.MCacheInuse)
	//fmt.Println(rmt.MCacheSys)
	//fmt.Println(rmt.MSpanInuse)
	//fmt.Println(rmt.MSpanSys)
	//fmt.Println(rmt.Mallocs)
	//fmt.Println(rmt.NextGC)
	//fmt.Println(rmt.NumForcedGC)
	//fmt.Println(rmt.NumGC)
	//fmt.Println(rmt.OtherSys)
	//fmt.Println(rmt.PauseTotalNs)
	//fmt.Println(rmt.StackInuse)
	//fmt.Println(rmt.StackSys)
	//fmt.Println(rmt.Sys)
}
