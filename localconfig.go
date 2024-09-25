package main

// we're just gonna hamfist this and not make a separate package
// because really we're only doing one thing and the heavy lifting
// should be performed by sux

type Config struct {
	Slack struct {
		appId         string `yaml:"app_id"`
		clientId      string `yaml:"client_id"`
		clientSecret  string `yaml:"client_secret"`
		signingSecret string `yaml:"signing_secret"`
	}
}

// let's just do some accessors because accessors are cool

func (c *Config) AppId() string {
	return c.Slack.appId
}

func (c *Config) ClientId() string {
	return c.Slack.clientId
}

func (c *Config) ClientSecret() string {
	return c.Slack.clientSecret
}

func (c *Config) SigningSecret() string {
	return c.Slack.signingSecret
}

// there are no mutators because you're not allowed to mutate these values
