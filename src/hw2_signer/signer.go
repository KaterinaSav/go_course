package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("Started")
}

func ExecutePipeline(jobs ...job) {
	ch1 := make(chan interface {}, 1)
	ch2 := make(chan interface {}, 1)
	for _, job := range jobs {
	    go job(ch1, ch2)
	}
}

var SingleHash = func(in, out chan interface{}) {
    dataRaw := <-in
    var (
    	string1 string
    	string2 string
    )
    data := dataRaw.(string)
    string1 = DataSignerCrc32(data)
    string2 = DataSignerCrc32(DataSignerMd5(data))
    out <- string1 + "~" + string2
}

var MultiHash = func(in, out chan interface{}) {
    dataRaw := <-in
    data := dataRaw.(string)
    iterationsNum := 5
    hashResults := make([]string, 0, iterationsNum + 1)
    for i := 0; i < iterationsNum; i++ {
        th = i.(string)
        crcResult = DataSignerCrc32(th + data)
        hashResults = append(hashResults, crcResult)
    }
    out <- hashResults
}

var CombineResults = func(in, out chan interface{}) {
    dataRaw := <-in
    sort.Strings(dataRaw)
    var dataString string
    dataString := strings.Join(dataRaw, "_")
    out <- dataString
}