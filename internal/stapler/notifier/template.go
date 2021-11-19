package notifier

import (
	"bytes"
	"html/template"

	"protocall/internal/stapler"
	"protocall/pkg/logger"
)

var (
	notifyMessage *template.Template
)

const (
	subject = "Протокол вашей встречи"

	body = `
<h1>Cтенограмма вашей конференции.</h1>
{{range $phrase := .Phrases}}
	<div class="message">
		<p>{{$phrase.Text}}</p>
		<p>Пользователь {{$phrase.User}}</p>
		<p>Время {{$phrase.Time}}</p>
	</div>
{{end}}
`
)

func init() {
	notifyMessage = template.Must(template.New("base").Parse(body))
}

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
