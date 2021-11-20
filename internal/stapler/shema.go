package stapler

type ProtocolRequest struct {
	Records []string `json:"records" binding:"required"`
	To      []string `json:"to" binding:"required"`
}
