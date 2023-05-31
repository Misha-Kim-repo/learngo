package main

import (
	"learngo/JopScrapper/scrapper"
	"os"
	"strings"

	"github.com/labstack/echo"
)

func handleHome(c echo.Context) error {
	//return c.String(http.StatusOK, "Hello, World!")
	//HTML을 전달하기 위해 기존에 사용하던 멤버 함수를 FIle() 함수로 변경한다.
	//매개 변수는 먼저 생성한 home.html을 지정한다.
	return c.File("home.html")
}

// 이름을 변경하지 않도록 *.csv 파일명의 const변수를 선언한다.
const fileName string = "jobs.csv"

// main() 함수에서 scrapper 기능을 전달하기 위해 해당 함수를 사용한다.
// 해당 함수는 term이라는 검색어를 Scrape() 함수에 전달하게 된다.
// scrapper에서 작업 후 만들어진 *.csv 파일을 반환한다.
func handleScrape(c echo.Context) error {

	//작업 완료 후 더 이상 필요 없는 *.csv 파일은 삭제한다.
	defer os.Remove(fileName)

	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	return c.Attachment(fileName, fileName)
}

// 기본적으로 go echo 서버에 나와있는 내용을 토대로 작성
// GET에 들어가는 익명 함수를 따로 추출하여 handleHome 함수로 작성
// 해당 함수를 매개 변수로 지정한다.
func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}
