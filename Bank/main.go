package main

import (
	"fmt"

	"learngo/accounts"
)

func main() {
	account := accounts.NewAccount("nico")
	account.Deposit(10)
	err := account.Withdraw(20)
	if err != nil {
		fmt.Println(err)
	}
	//account를 Print하는 것으로 객체 정보가 출력된다(String() 메소드)
	fmt.Println(account)
}
