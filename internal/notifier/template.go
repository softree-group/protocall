package notifier

import (
	"bytes"
	"html/template"

	"protocall/internal/stapler"
	"protocall/pkg/logger"
)

const (
	subject    = "Протокол вашей встречи"
	errMessage = "Данных нет"
	body       = `
<!DOCTYPE html>
<html>
	<head>
		<title>Protocall</title>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		<meta name="theme-color" content="#282c34" />
	</head>
	<body style="margin: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothing: grayscale; background: #282c34; font-size: 16px;">
		<h1 style="font-size: 21px; color: gainsboro; margin: 10px;">Протокол вашей конференции от {{(index .Phrases 0).Time.Format "02.01"}}.</h1>
		{{range $phrase := .Phrases}}
			<div class="message" style="background: #3f4148; padding: 10px; margin: 5px; border-radius: 10px; color: rgb(223, 223, 223);">
				<div class="message_user-container" style="position: relative; font-weight: 700; width: 100%; margin-bottom: 5px;">
					<span class="message_user">{{$phrase.User}}</span>
					<span class="message_time" style="position: absolute; right: 0; font-weight: 200;">{{$phrase.Time.Format "15:04"}}</span>
				</div>
				<p class="message_text">{{$phrase.Text}}</p>
			</div>
		{{end}}
	</body>
</html>
`
)

var notifyMessage = template.Must(template.New("base").Parse(body))

func render(phrases []stapler.Phrase) string {
	res := &bytes.Buffer{}
	if err := notifyMessage.ExecuteTemplate(res, "base", struct {
		Phrases []stapler.Phrase
	}{
		Phrases: phrases,
	}); err != nil {
		logger.L.Error(err)
		return errMessage
	}
	return res.String()
}
