package account

type Account struct {
	Username string `redis:"username" json:"username"`
	Password string `redis:"password" json:"password"`
	UserID   string `redis:"id"`
}

func (a *Account) Less(then Account) bool {
	return a.Username < then.Username
}
