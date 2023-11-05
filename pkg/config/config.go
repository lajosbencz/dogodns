package config

const DefaultPath = "dogodns.yaml"
const DefaultPIP = "https://ipecho.net/plain"
const DefaultInterval = 60
const DefaultTTL = 300

type Config struct {
	Domain   *string  `yaml:"domain,omitempty"`
	Domains  []string `yaml:"domains"`
	Token    string   `yaml:"token"`
	PIP      string   `yaml:"pip"`
	Interval int      `yaml:"interval"`
	TTL      int      `yaml:"ttl"`
}

func (c *Config) GetDomains() []string {
	l := []string{}
	if c.Domain != nil {
		l = append(l, *c.Domain)
	}
	l = append(l, c.Domains...)
	return l
}

func DefaultConfig(domain, token string) *Config {
	return &Config{
		Domains:  []string{domain},
		Token:    token,
		PIP:      DefaultPIP,
		Interval: DefaultInterval,
		TTL:      DefaultTTL,
	}
}
