package main

import (
	"fmt"
	"time"
)

// channel 매개변수는 type이 중요하다
// 여기에서는 string 타입의 값을 파이프로 전달한다(c chan string)
func isSexy(person string, c chan string) {
	time.Sleep(time.Second * 10)
	fmt.Println(person)
	c <- person + " is sexy"
}

func main() {
	c := make(chan string)
	people := [2]string{"nico", "flynn"}
	for _, person := range people {
		go isSexy(person, c)
	}
	fmt.Println("Waiting for messages")
	// resultOne := <-c
	// resultTwo := <-c
	// resultThree := <-c

	//blocking operation: 프로그램을 병렬로 처리할 수 있게끔 해당 작업이 끝날때까지 대기한다
	//channel로부터 메시지를 가져온다
	// fmt.Println("Received this message: ", resultOne)
	// fmt.Println("Received this message: ", resultTwo)
	// fmt.Println("Received this message: ", resultThree)

	//변수를 각각 선언하여 데이터를 수신하는 것은 Deadlock을 발생시킬 수 있다
	//따라서 Deadlock을 방지하기 위해 for문과 iterator를 사용하여 people을 순회하는 방식으로 진행
	for i := 0; i < len(people); i++ {
		fmt.Println("waiting for ", i)
		fmt.Println(<-c)
	}
}
