package applications

<<<<<<< HEAD
import (
	"protocall/domain/entity"
)
=======
import "protocall/domain/entity"
>>>>>>> 977da2b (rebase inbloud)

type Conference interface {
	StartConference(user *entity.User) (*entity.Conference, error)
	JoinToConference(user *entity.User, meetID string) (*entity.Conference, error)
	IsExist(meetID string) bool
<<<<<<< HEAD
	StartRecord(user *entity.User, meetID string) error
	Get(meetID string) *entity.Conference
	StartRecordUser(user *entity.User, meetID string) error
=======
>>>>>>> 977da2b (rebase inbloud)
}
