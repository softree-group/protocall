package stapler

type ProtocolRequest struct {
	ConferenceID string   `json:"conference_id" binding:"required"`
	To           []string `json:"to" binding:"required"`
}
