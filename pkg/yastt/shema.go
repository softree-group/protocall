package yastt

import "time"

type Specification struct {
	LanguageCode      string `yaml:"languageCode"      json:"languageCode"`
	Model             string `yaml:"model"             json:"model"`
	ProfanityFilter   string `yaml:"profanityFilter"   json:"profanityFilter"`
	AudioEncoding     string `yaml:"audioEncoding"     json:"audioEncoding"`
	SampleRateHertz   string `yaml:"sampleRateHertz"   json:"sampleRateHertz"`
	AudioChannelCount int    `yaml:"audioChannelCount" json:"audioChannelCount"`
	RawResults        bool   `yaml:"rawResults"        json:"rawResults"`
}

type Config struct {
	Specification Specification `json:"specification" binding:"required"`
}

type Audio struct {
	URI string `json:"uri" binding:"required"`
}

type RecognizerRequest struct {
	Config Config `json:"config" binding:"required"`
	Audio  Audio  `json:"audio" binding:"required"`
}

type Word struct {
	StartTime time.Duration `json:"startTime" binding:"required"`
	EndTime   time.Duration `json:"endTime" binding:"required"`
	Word      string        `json:"word" binding:"required"`
}

type Alternative struct {
	Words []Word `json:"words" binding:"required"`
	Text  string `json:"text" binding:"required"`
}

type Chunk struct {
	Alternatives []Alternative `json:"alternatives" binding:"required"`
	ChannelTag   int           `json:"channelTag" binding:"required"`
}

type Response struct {
	Chunks []Chunk `json:"chunks"`
}

type RecognizerResponse struct {
	Done       bool     `json:"done" binding:"required"`
	ID         string   `json:"id" binding:"required"`
	CreatedAt  string   `json:"createdAt" binding:"required"`
	CreatedBy  string   `json:"createdBy" binding:"required"`
	ModifiedAt string   `json:"modifiedAt" binding:"required"`
	Response   Response `json:"response" binding:"required"`
}
