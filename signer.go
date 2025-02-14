package main

import (
	"fmt"
	"strconv"
)

func ExecutePipeline(jobs ...job) {
	fmt.Println("ExecutePipeline")
	in, out := make(chan interface{}, 100), make(chan interface{}, 100)

	for i, j := range jobs { //TODO: сделать правильную обработку большого количества входящих чисел (больше одного)
		fmt.Println(i)
		j(in, out) // здесь надо go
		in, out = out, make(chan interface{}, 100)
	}

}

func SingleHash(in, out chan interface{}) {
	data, _ := <-in

	dataInt, ok := data.(int)
	if !ok {
		panic("Conversion failed")
	}

	dataStr := strconv.Itoa(dataInt)

	md := DataSignerMd5(dataStr)
	crcOut1 := make(chan interface{})
	crcOut2 := make(chan interface{})
	go performCrc32(dataStr, crcOut1)
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
	fmt.Println(resStr)
}

func MultiHash(in, out chan interface{}) {
	data := <-in
	fmt.Println(data)

	dataStr, ok := data.(string)
	if !ok {
		panic("Conversion failed")
	}

	chans := make([]chan interface{}, 6)
	for i := 0; i < 6; i++ {
		chans[i] = make(chan interface{})
		go performCrc32(strconv.Itoa(i)+dataStr, chans[i])
	}

	resStr := ""
	for i := 0; i < 6; i++ {
		res := <-chans[i]
		str, ok := res.(string)
		if !ok {
			panic("Conversion failed")
		}
		resStr += str
	}

	out <- resStr
	fmt.Println(resStr)
}

func performCrc32(data string, out chan interface{}) {
	out <- DataSignerCrc32(data)
}

func CombineResults(in, out chan interface{}) {

}
