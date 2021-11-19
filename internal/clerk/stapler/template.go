package stapler

import "html/template"

var (
	smtpTmpl *template.Template
)

const protocallTmpl = `<p>По этой <a href="{{.Link}}">ссылке</a> вы можете скачать стенограмму конференции.</p>`

func init() {
	smtpTmpl = template.Must(template.New("base").Parse(protocallTmpl))
}
