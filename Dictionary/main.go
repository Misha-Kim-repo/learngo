package main

import (
	"fmt"

	"learngo/Dictionary/mydict"
)

func main() {
	//Map 형 자료구조를 활용하여 Search 함수로 value 값 검색해보기
	// dictionary := mydict.Dictionary{"First": "First word"}
	// definition, err := dictionary.Search("Second")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(definition)
	// }

	//Map 형 자료구조를 활용하여 Add 함수로 값(Key-Value) 입력해보기
	// dictionary := mydict.Dictionary{}
	// word := "Hello"
	// definition := "Greeting"
	// errHello := dictionary.Add(word, definition)
	// if errHello != nil {
	// 	fmt.Println(errHello)
	// }
	// hello, errHello := dictionary.Search(word)
	// fmt.Println(hello)

	//이미 추가한 단어를 Map 형 자료구조에서 검색하여 에러 반환해보기
	// errAdd := dictionary.Add(word, definition)
	// if errAdd != nil {
	// 	fmt.Println(errAdd)
	// }

	//이미 추가한 단어를 갱신해보자
	// dictionary := mydict.Dictionary{}
	// baseWord := "Hello"
	// dictionary.Add(baseWord, "First")
	// err := dictionary.Update(baseWord, "Second")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// word, _ := dictionary.Search(baseWord)
	// fmt.Println(word)

	//이미 추가한 단어를 삭제해보자
	dictionary := mydict.Dictionary{}
	baseWord := "Hello"
	dictionary.Add(baseWord, "First")
	dictionary.Search(baseWord)
	err := dictionary.Delete(baseWord)
	word, _ := dictionary.Search(baseWord)
	if err != nil {
		fmt.Println(word)
	}
	fmt.Println("Can't find the word")
}
