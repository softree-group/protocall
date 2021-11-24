package yastt

type Specification struct {
	LanguageCode      string `yaml:"languageCode"      json:"languageCode,omitempty"`
	Model             string `yaml:"model"             json:"model omitempty"`
	ProfanityFilter   string `yaml:"profanityFilter"   json:"profanityFilter,omitempty"`
	AudioEncoding     string `yaml:"audioEncoding"     json:"audioEncoding,omitempty"`
	SampleRateHertz   int    `yaml:"sampleRateHertz"   json:"sampleRateHertz,omitempty"`
	AudioChannelCount int    `yaml:"audioChannelCount" json:"audioChannelCount,omitempty"`
	RawResults        bool   `yaml:"rawResults"        json:"rawResults,omitempty"`
}

type Config struct {
	Specification Specification `json:"specification"`
}

type Audio struct {
	URI string `json:"uri"`
}

type RecognizerRequest struct {
	Config Config `json:"config"`
	Audio  Audio  `json:"audio"`
}

type Word struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Word      string `json:"word"`
}

type Alternative struct {
	Words []Word `json:"words"`
	Text  string `json:"text"`
}

type Chunk struct {
	Alternatives []Alternative `json:"alternatives"`
	ChannelTag   string        `json:"channelTag"`
}

type Response struct {
	Chunks []Chunk `json:"chunks"`
}

type RecognizerResponse struct {
	Done       bool     `json:"done"`
	ID         string   `json:"id"`
	CreatedAt  string   `json:"createdAt"`
	CreatedBy  string   `json:"createdBy"`
	ModifiedAt string   `json:"modifiedAt"`
	Response   Response `json:"response"`
}
