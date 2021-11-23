package stapler

type User struct {
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	NeedProtocol bool     `json:"need_protocol`
	Records      []string `json:"records"`
	Texts        []string `json:"texts"`
}

type ProtocolRequest struct {
	Users []User `json:"users"`
}
