package notifier

import (
	"bytes"
	"html/template"
	"time"

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
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
            }

        body {
            margin: 0;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;
            background: #282c34; 
            font-size: 16px;
        }

        body h1 {
            font-size: 21px;
            color: gainsboro;
            margin: 10px;
        }

        .message {
            background: #3f4148;
            padding: 10px;
            margin: 5px;
            border-radius: 10px;
            color: rgb(223, 223, 223);
        }

        .message_user-container {
            font-weight: 700;
            width: 100%;
            margin-bottom: 5px;
            display: flex;
            justify-content: space-between;
        }

        .message_time {
            font-weight: 200;
        }
    </style>
	<body>
		<h1>Протокол вашей конференции от {{.Time.Format "02.01"}}.</h1>
		{{range $phrase := .Phrases}}
			<div class="message">
                <div class="message_user-container">
				    <span class="message_user">{{$phrase.User}}</span>
				    <span class="message_time">{{$phrase.Time.Format "15:04"}}</span>
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
		Time time.Time
	}{
		Phrases: phrases,
		Time: phrases[0].Time,
	}); err != nil {
		logger.L.Error(err)
		return ""
	}
	return res.String()
}
