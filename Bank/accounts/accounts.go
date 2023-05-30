package accounts

import (
	"errors"
	"fmt"
)

// Account Struct
type Account struct {
	owner   string
	balance int
}

// 신규 Account를 생성하기 위한 함수(constructor)
// NewAccount creates Account
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

// method를 활용한 Deposit 함수. 예금된 액수에 추가된 만큼 금액을 더하는 작업을 수행
// Deposit x amount on your account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

// Balance of your Account
func (a Account) Balance() int {
	return a.balance
}

// 에러 처리를 위한 에러 변수 선언
var errNoMoney = errors.New("Can't withdraw")

// Withdraw x amount from your Account
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil
}

// ChangeOwner of the account
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

// Owner of the account
func (a Account) Owner() string {
	return a.owner
}

// 객체 정보를 출력하기 위한 함수(String) 재정의
func (a Account) String() string {
	return fmt.Sprint(a.Owner(), "'s account.\nHas: ", a.Balance())
}
