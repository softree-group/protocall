package account

type AccountStorage interface {
	GetAccount(account string) *Account
	GetFree() *Account
	TakeAccount(account string, userID string)
	FreeAccount(account string)
	SaveAccount(id int, account *Account)
	Who(account string) string
}

type Application struct {
	storage AccountStorage
}

func NewApplication(storage AccountStorage) *Application {
	return &Application{
		storage: storage,
	}
}

func (a *Application) GetFree() *Account {
	return a.storage.GetFree()
}

func (a *Application) Get(account string) *Account {
	return a.storage.GetAccount(account)
}

func (a *Application) Take(account, userID string) {
	a.storage.TakeAccount(account, userID)
}

func (a *Application) Free(account string) {
	a.storage.FreeAccount(account)
}

func (a *Application) Who(account string) string {
	return a.storage.Who(account)
}
