package main

import (
	"fmt"
	"sync"
	"runtime"
	"strconv"
)

func main() {
	fmt.Println("Started")
	inputData := []int{0,1}

    hashSignJobs := []job{
        job(func(in, out chan interface{}) {
            for _, fibNum := range inputData {
                fmt.Println("save to out")
                fmt.Println(fibNum)
                out <- fibNum
            }
        }),
        job(SingleHash),
        job(MultiHash),
        job(func(in, out chan interface{}) {
            dataRaw := <-in
            data := fmt.Sprintf("%v", dataRaw)
            fmt.Println("data result" + data)
        }),
    }

    ExecutePipeline(hashSignJobs...)
    fmt.Scanln()
}

func ExecutePipeline(jobs ...job) {
	in := make(chan interface {}, 3)
    out := make(chan interface {}, 3)
    wg := &sync.WaitGroup{}
	for _, job := range jobs {
        wg.Add(1)
	    go Pipeline(wg, job, in, out)
	}

	go func(in, out chan interface {}) {
        for {
            select {
                case <-out:
                    fmt.Println("save out to in")
                    data := fmt.Sprintf("%v", <-out)
                    fmt.Printf(data)
                    in <- data
                case <-in:
                    data := fmt.Sprintf("%v", <-in)
                                fmt.Printf(data)
                 fmt.Println("something saved to in")
            }
        }
    }(in, out)
	wg.Wait()
}

func Pipeline(wg *sync.WaitGroup, currentJob job, in chan interface {}, out chan interface {}) {
    defer wg.Done()
    currentJob(in, out)
    runtime.Gosched()
}

var SingleHash = func(in, out chan interface{}) {
    dataRaw := <-in
    data := fmt.Sprintf("%v", dataRaw)
    fmt.Println(data + " SingleHash " + "data " + data)
    md5Data := DataSignerMd5(data)
    fmt.Println(data + " SingleHash " + "md5(data) " + md5Data)
    crc32md5Data := DataSignerCrc32(md5Data)
    fmt.Println(data + " SingleHash " + "crc32(md5(data)) " + crc32md5Data)
    crc32Data := DataSignerCrc32(data)
    fmt.Println(data + " SingleHash " + "crc32(data) " + crc32Data)
    result := crc32Data + "~" + crc32md5Data
    fmt.Println(data + " SingleHash " + "result " + result)
    out <- result
}

var MultiHash = func(in, out chan interface{}) {
    dataRaw := <-in
    data := fmt.Sprintf("%v", dataRaw)
    iterationsNum := 6
    result := ""
    for i := 0; i < iterationsNum; i++ {
        th := strconv.Itoa(i)
        crcResult := DataSignerCrc32(th + data)
        fmt.Println(data + " MultiHash " + "crc32(th+step1)) " + th + " " + crcResult)
        result = result + crcResult
    }
    fmt.Println(data + " MultiHash " + "result " + result)
    out <- result
}

var CombineResults = func(in, out chan interface{}) {
    dataRaw := <-in
    data := dataRaw.(string)
    out <- data
}