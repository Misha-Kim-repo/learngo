package main

import (
	"errors"
	"fmt"
	"net/http"
)

var errRequestFailed = errors.New("Request Failed")

//goroutines 추가로 인하여 해당 함수는 더 이상 error 값을 변환하지 않는다
// func hitURL(url string) error {
// 	fmt.Println("Checking: ", url)
// 	resp, err := http.Get(url)
// 	if err != nil || resp.StatusCode >= 400 {
// 		fmt.Println(err, resp.StatusCode)
// 		return errRequestFailed
// 	}
// 	return nil
// }

// channel을 전달하기 위한 struct 선언
type requestResult struct {
	url    string
	status string
}

// 신규 추가 함수
// channel을 처리하기 위해 신규 channel 타입 패러미터 추가
// 해당 패러미터는 send only이다
func hitURL(url string, c chan<- requestResult) {
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
	}
	c <- requestResult{url: url, status: status}
}

func main() {
	//channel을 사용하기 위해 struct 형 channel 추가
	var results = make(map[string]string)
	c := make(chan requestResult)
	urls := []string{
		"https://www.airbnb.com/",
		"https://www.google.com",
		"https://www.amazon.com/",
		"https://www.reddit.com",
		"https://www.google.com",
		"https://soundcloud.com",
		"https://www.facebook.com",
		"https://instagram.com",
		"https://academy.nomadcoders.co",
	}
	for _, url := range urls {
		go hitURL(url, c)
		// goroutines 추가로 인한 코드 삭제 처리 - 1
		// if err != nil {
		// 	result = "FAILED"
		// }
		// results[url] = result
	}
	// goroutines 추가로 인한 코드 삭제 처리 - 2
	// for url, result := range results {
	// 	fmt.Println(url, result)
	// }

	//hitURL() 함수를 통해 전송된 결괏값 대입
	//for문을 통해 urls slice를 차례대로 순회
	for i := 0; i < len(urls); i++ {
		result := <-c
		results[result.url] = result.status
	}

	//전달받은 값을 for문을 통해 순서대로 출력
	for url, status := range results {
		fmt.Println(url, status)
	}
}
