package user

type StorageRepository interface {
	FindUser(sessionID string) *User
	SaveUser(user *User)
	DeleteUser(sessionID string)
}

type Application struct {
	storage StorageRepository
}

func NewApplication(storage StorageRepository) *Application {
	return &Application{storage}
}

func (a *Application) Find(sessionID string) *User {
	return a.storage.FindUser(sessionID)
}

func (a *Application) Save(user *User) {
	a.storage.SaveUser(user)
}

func (a *Application) Delete(sessionID string) {
	a.storage.DeleteUser(sessionID)
}
