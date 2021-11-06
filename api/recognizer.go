package api

type TextRespone struct {
	Result []struct {
		Conf  float64 `json:"conf" binding:"required"`
		Start float64 `json:"start" binding:"required"`
		End   float64 `json:"end" binding:"required"`
		Word  string  `json:"word" binding:"required"`
	} `json:"result" binding:"required"`
	Text string `json:"text" binding:"required"`
}

var EOF = `{"eof" : 1}`
