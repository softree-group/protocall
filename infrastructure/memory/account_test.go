package memory

import (
	"fmt"
	"testing"
)

func TestAsteriskAccountMemory_Take(t *testing.T) {
	repo := NewAsteriskAccountMemory()

	account := repo.GetFree()
	fmt.Printf("%v+", account)
	repo.Take(account.Username, account.Username)
	account = repo.GetFree()
	fmt.Printf("%v+", account)
	repo.Take(account.Username, account.Username)
	account = repo.GetFree()
	fmt.Printf("%v+", account)
	repo.Take(account.Username, account.Username)
	account = repo.GetFree()
	fmt.Printf("%v+", account)
	repo.Free("1234")
	account = repo.GetFree()
	fmt.Printf("%v+", account)
}
