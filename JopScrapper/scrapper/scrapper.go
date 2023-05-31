package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	ccsv "github.com/tsak/concurrent-csv-writer"
)

//패키지화를 위한 전역변수 삭제
//var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=python"

type extractedJob struct {
	title    string
	location string
	deadline string
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status Code: ", res.StatusCode)
	}
}

// cleanString함수 작성
// strings 패키지의 strings 객체를 활용하여 접근한다.
// TrimSpace(s string): 스페이스(" ")를 없애버리는 함수
// Fields(s string): 텍스트로만 이루어진 배열을 생성하는 함수
// Join(a []string, sep string) string : 배열을 패러미터로 받아 sep을 통해 문자열을 합치는 함수
// >> 따라서 Fiels 함수에서 얻어낸 string 형의 배열을 Join함수에서 하나로 합쳐 string 형으로 반환한다

// 패키지화의 일환으로 해당 함수를 노출시키기 위해 대문자화
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// 실행 시간을 측정하기 위한 함수 생성
func ElapsedTime(tag string, msg string) func() {
	if msg != "" {
		log.Printf("[%s] %s", tag, msg)
	}

	start := time.Now()
	return func() { log.Printf("[%s] elapsed Time: %s ", tag, time.Since(start)) }
}

// 주어진 page 매개변수를 통해 각 page에 접근할 수 있는 URL을 생성한다.
// main() 함수와의 비동기 처리를 위해 channel 매개변수를 추가한다.
// func getPage(page int) []extractedJob {
// func getPage(page int, mainC chan<- []extractedJob) {
// scrapper 패키지화를 위해 함수의 매개변수가 변경된다 (string 타입의 url 변수 추가)
func getPage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob

	//extractJob()의 비동기화를 위한 channel 생성
	//전달받은 값을 리턴하는 대신 channel에서 전달받는 데이터를 수신하기 위핸 변수를 선언한다.
	c := make(chan extractedJob)
	pageURL := url + "&recruitPage=" + strconv.Itoa(page+1) + "&recruitPageCount=" + strconv.Itoa(40)
	fmt.Println("Requesting: " + pageURL)

	//page를 요청해보기
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	//fmt.Println(doc)

	//jobCard 찾기
	//현재 사*인 사이트를 통해 해당 부분을 search하고 있어 JobCard가 아닌 jobTitle을 찾는다
	//Find-Each는 아래와 같이 분리시켜 사용이 가능하다
	//selection 포인터는 각각의 채용정보 단 건을 의미한다. 따라서 각 채용정보의 정보를 취득하고 싶을 경우,
	//selection 포인터 객체를 통해 접근한다.
	//해당 부분에서 추출하는 것은 titleJob/location/deadline 총 세가지다
	searchJobTitle := doc.Find(".item_recruit")
	searchJobTitle.Each(func(i int, card *goquery.Selection) {
		//기존에 수행했던 부분을 신규 구성한 함수 extractJob(card *goquery.Selection)에서 실행한다.
		//해당 결과를 job 객체에 입력
		//append 함수를 통해 job 객체를 추가한다
		//extractJob() 함수를 처리하기 위해 goroutine 처리를 진행한다.
		//로직이 변경되었으므로, 해당 작업이 끝난 상태에서 append가 필요하므로, append() 함수 사용처를 변경
		go extractJob(card, c)
		//jobs = append(jobs, job)
	})

	//비동기화 된 extracteJob() 함수의 값을 저장하기 위해 for문을 통한 순회로 값을 추가(append)시킨다.
	for i := 0; i < searchJobTitle.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	//전달받은 값을 리턴해준다
	//return jobs
	//비동기화를 위해 값을 저장한 jobs 배열을 mainC 채널로 전송한다.
	mainC <- jobs

}

// getPages(page int) 함수에서 Job 단 건(card)의 정보를 추출하기 위한 함수를 신규 생성
// 기존의 getPages(page int) 함수에서는 div-class item_recruit의 문서 정보만 읽는다
// 해당 값들을 반환하기 위해 extractedJob struct를 사용한다
// getPage()에서 비동기화 처리를 위해 리턴값을 가지는 대신 채널을 수신한다
// func extractJob(card *goquery.Selection) extractedJob {
func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	title := CleanString(card.Find(".area_job>.job_tit>a").Text())
	location := CleanString(card.Find(".area_job>.job_condition>span>a").Text())
	deadline := CleanString(card.Find(".area_job>.job_date>span").Text())

	// 리턴 값 대신 비동기화를 위해 채널로 해당 값들을 전송한다.
	// return extractedJob{
	// 	title:    title,
	// 	location: location,
	// 	deadline: deadline,
	// }
	//fmt.Println(id, location, deadline)
	c <- extractedJob{
		title:    title,
		location: location,
		deadline: deadline,
	}
}

func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	//I/O를 통해 열어준 Body를 닫는다.
	//메모리 누수를 방지
	defer res.Body.Close()

	//res.Body로 받은 내용을 열어 해당 내용을 doc 및 err에 대입한다
	//res.Body는 HTML 문서를 로드한다
	//res.Body: 기본적으로 byte이며, I/O을 의미한다.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	//doc.Find를 통해 무언가를 찾아올 것이며, 찾아온 아이템으로 추가적인 조치를 취할 수 있다.
	//Each()에서는 selection을 통해 각 pagination에서 받아온 값으로 추가적인 처리를 진행한다.
	//여기에서는 selection으로 각 페이지의 총 갯수(a href로 시작하는)를 받아온다.
	//Find-Each()함수를 통해 받아온 pages를 반환하는 형식으로 진행
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

// 받아온 채용정보를 csv에 쓰는 역할을 하는 함수
// func writeJobs(jobs []extractedJob) {

// 	//go언어에서 제공하는 os 패키지를 사용한다
// 	//해당 패키지를 이용하여 jobs.csv 파일을 생성할 수 있다.
// 	file, err := os.Create("jobs.csv")
// 	checkErr(err)

// 	//writer를 사용하여 csv 신규 파일을 구성한다.
// 	//구성한 파일을 통해 함수가 끝나는 시점에서 Flush() 함수를 사용하여 데이터를 파일에 저장한다
// 	w := csv.NewWriter(file)
// 	defer ElapsedTime("writeJobs()", "start")()
// 	defer w.Flush()

// 	headers := []string{"title", "location", "deadLine"}

// 	wErr := w.Write(headers)
// 	checkErr(wErr)

// 	//매개변수로 받은 jobs 파일을 순회하며 신규 배열(jobSlice)에 쓴다
// 	for _, job := range jobs {
// 		jobSlice := []string{job.title, job.location, job.deadline}
// 		jwErr := w.Write(jobSlice)
// 		checkErr(jwErr)
// 	}
// }

// 받아온 채용정보를 csv에 쓰는 역할을 하는 함수 - mk.2
// 기존 함수와 다르게 goroutine-channel을 통한 비동기화 처리
// 기존 함수 대비 0.2초 정도의 속도가 향상된 것을 확인할 수 있다
func writeJob(jobs []extractedJob) {
	csv, err := ccsv.NewCsvWriter("jobs.csv")
	checkErr(err)

	defer csv.Close()

	//blocking operation을 위한 channel 생성
	done := make(chan bool)

	//csv 파일에 headers를 쓴다
	headers := []string{"title", "location", "deadLine"}
	wErr := csv.Write(headers)
	checkErr(wErr)

	//for문 순회 및 goroutine을 활용한 비동기화 진행
	for _, job := range jobs {
		go func(job extractedJob) {
			jobSlice := []string{job.title, job.location, job.deadline}
			jwErr := csv.Write(jobSlice)
			checkErr(jwErr)
			<-done
		}(job)
	}
}

// 패키지화를 위해 함수명 변경
// main() > scraper()로 변경한다.
// go echo 서버를 사용하기 위해 매개변수 term을 받는다.
// baseURL 전역변수 삭제를 위해 Scrape() 함수에서 대신 URL을 선언하고 있으며, 해당 변수는 getpage()에서 사용한다.
// func main() {
func Scrape(term string) {
	//Job에 대한 정보를 받기 위한 extractedJob 객체 slice 생성
	//getPage() 함수를 수신하기 위한 채널 신규 생성
	//패키지화를 위해 URL을 받는다.
	var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=" + term
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages(baseURL)

	for i := 0; i < totalPages; i++ {
		//getPages() 함수를 통해 신규 slice인 extractedJobs를 채운다
		//해당 값을 통해 신규 job을 찾을 때 마다 그 내용을 배열에 추가(append)
		//extractedJobs := getPage(i)
		//jobs = append(jobs, extractedJobs...)
		// getPage() 함수를 비동기 처리하기 위해 goroutine을 사용한다.
		// scrapper 패키지화에 따른, getPage()의 매개변수 추가
		go getPage(i, baseURL, c)
	}

	//처리된 데이터를 채널을 통해 수신
	//수신된 데이터를 기존 배열에 totalPages 만큼 추가한다.(append)
	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}
	//main 함수에서의 writeJobs 함수 사용
	writeJob(jobs)

	fmt.Println("Done, extracted", len(jobs))
}
