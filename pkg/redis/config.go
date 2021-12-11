package redis

type Config struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	MaxIdle     int    `yaml:"maxIdle"`
	IdleTimeout int    `yaml:"idleTimeout"`
}
