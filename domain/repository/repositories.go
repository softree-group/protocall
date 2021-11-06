package repository

type Repositories interface {
	Bridge
	AsteriskAccountRepository
	User
	Conference
	ConferenceStorage
	ConferenceTranslator
}
