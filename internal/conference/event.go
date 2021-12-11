package conference

import (
	"protocall/internal/translator"
	"protocall/internal/user"
)

type Event struct {
	Record       *translator.Record
	User         *user.User
	ConferenceID string
	Text         string
}
