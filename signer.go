package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})

	wg := sync.WaitGroup{}
	for _, j := range jobs {
		out := make(chan interface{})
		wg.Add(1)
		go func(j job, in chan interface{}) {
			defer wg.Done()
			defer close(out)
			j(in, out)
		}(j, in)
		in = out
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	for data := range in {
		dataInt, ok := data.(int)
		if !ok {
			panic("Conversion failed")
		}
		dataStr := strconv.Itoa(dataInt)
		md := DataSignerMd5(dataStr)
		wg.Add(1)
		go func(md string, data string) {
			defer wg.Done()
			crcOut1 := make(chan interface{})
			crcOut2 := make(chan interface{})
			go performCrc32(data, crcOut1)
			go performCrc32(md, crcOut2)

			res1 := <-crcOut1
			str1, ok := res1.(string)
			if !ok {
				panic("Conversion failed")
			}
			res2 := <-crcOut2
			str2, ok := res2.(string)
			if !ok {
				panic("Conversion failed")
			}

			resStr := str1 + "~" + str2
			out <- resStr
		}(md, dataStr)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	for data := range in {

		dataStr, ok := data.(string)
		if !ok {
			panic("Conversion failed")
		}

		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			chans := make([]chan interface{}, 6)
			for i := 0; i < 6; i++ {
				chans[i] = make(chan interface{})
				go performCrc32(strconv.Itoa(i)+data, chans[i])
			}

			resStr := strings.Builder{}
			for i := 0; i < 6; i++ {
				res := <-chans[i]
				str, ok := res.(string)
				if !ok {
					panic("Conversion failed")
				}
				resStr.WriteString(str)
			}

			out <- resStr.String()
		}(dataStr)
	}
	wg.Wait()
}

func performCrc32(data string, out chan interface{}) {
	out <- DataSignerCrc32(data)
}

func CombineResults(in, out chan interface{}) {
	resArr := make([]string, 0)
	for data := range in {
		dataStr, ok := data.(string)
		if !ok {
			panic("Conversion failed")
		}
		resArr = append(resArr, dataStr)
	}

	sort.Strings(resArr)
	out <- strings.Join(resArr, "_")
}
