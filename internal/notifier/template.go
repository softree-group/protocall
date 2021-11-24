package notifier

import (
	"bytes"
	"html/template"

	"protocall/internal/stapler"
	"protocall/pkg/logger"
)

const (
	subject = "Протокол вашей встречи"

	body = `
<!DOCTYPE html>
<html>
	<head>
		<title>Protocall</title>
		<meta charset="utf-8">
	</head>
	<body>
		<h1>Протокол вашей конференции.</h1>
		{{range $phrase := .Phrases}}
			<div class="message">
				<p class="message_text">{{$phrase.Text}}</p>
				<span class="message_user">{{$phrase.User}}</span>
				<span class="message_time">{{$phrase.Time.Format "15:04"}}</span>
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
		return ""
	}
	return res.String()
}
