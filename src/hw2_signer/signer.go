package main

import (
	"fmt"
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
    in := make(chan interface {})
    out := make(chan interface {})
	Pipeline(jobs, in , out)
}

type blankChan chan interface {}

func Pipeline(jobs []job, in blankChan, out blankChan) {
    jobsCount := len(jobs)
    chans := make([]blankChan, 0, jobsCount)
    for range jobs {
            ch1 := make(blankChan)
            chans = append(chans, ch1)
        }
    for i := range jobs {
        i := i
        if i == 0 {
            go jobs[i](in, chans[i])
        } else if i == (jobsCount - 1) {
            go jobs[i](chans[i - 1], out)
        } else {
            go jobs[i](chans[i - 1], chans[i])
        }
    }
}

var SingleHash = func(in, out chan interface{}) {
    for {
        select {
            case dataRaw, ok := <-in:
            if !ok {
                close(out)
                return
            }
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
        }

}

var MultiHash = func(in, out chan interface{}) {
    for {
        select {
            case dataRaw, ok := <-in:
                if ok == false {
                    close(out)
                    return
                }
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
        }
}

var CombineResults = func(in, out chan interface{}) {
    dataRaw := <-in
    data := dataRaw.(string)
    out <- data
}