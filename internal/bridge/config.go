package bridge

type Config struct {
	Application       string `"yaml:"application"`
	SnoopyApplication string `"yaml:"snoopy"`
	URL               string `"yaml:"url"`
	Websocket         string `"yaml:"ws"`
	User              string `"yaml:"user"`
	Password          string `"yaml:"password"`
}
