package main

import (
	"fmt"
	"time"
)

func sexyCount(person string) {
	for i := 0; i < 10; i++ {
		fmt.Println(person, "is Sexy", i)
		time.Sleep(time.Second)
	}
}

func main() {
	go sexyCount("nico")
	go sexyCount("flynn")
	time.Sleep(time.Second * 10)
}
