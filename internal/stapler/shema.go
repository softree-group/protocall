package stapler

type User struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Record   string `json:"record" binding:"required"`
	Text     string `json:"text" binding:"required"`
}

type ProtocolRequest struct {
	Users []User `json:"users" binding:"required"`
}
