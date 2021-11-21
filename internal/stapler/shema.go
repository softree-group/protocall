package stapler

type User struct {
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email"`
	NeedProtocol bool   `json:"need_protocol binding:"required"`
	Records       []string `json:"records" binding:"required"`
	Texts         []string `json:"texts" binding:"required"`
}

type ProtocolRequest struct {
	Users []User `json:"users" binding:"required"`
}
