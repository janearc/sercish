package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

// we're just gonna hamfist this and not make a separate package
// because really we're only doing one thing and the heavy lifting
// should be performed by sux

type Config struct {
	Slack struct {
		appId         string `yaml:"app_id"`
		clientId      string `yaml:"client_id"`
		clientSecret  string `yaml:"client_secret"`
		signingSecret string `yaml:"signing_secret"`
		oauthToken    string `yaml:"oauth_token"`
	}
	NOAA struct {
		apiURL    string `yaml:"api_url"`
		stationID string `yaml:"station_id"`
	}
	Web struct {
		Port int `yaml:"port"`
	}
	Version struct {
		BuildDate string `yaml:"build_date"`
		Build     string `yaml:"build"`
		Branch    string `yaml:"branch"`
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

func (c *Config) OAuthToken() string {
	return c.Slack.oauthToken
}

func (c *Config) ApiURL() string {
	return c.NOAA.apiURL
}

func (c *Config) StationID() string {
	return c.NOAA.stationID
}

// there are no mutators because you're not allowed to mutate these values

// LoadConfig smashes a bunch of yaml into a config object we'll need everywhere
func LoadConfig(configFileName string, versionFileName string, secretsFileName string) (*Config, error) {
	if configFileName == "" {
		configFileName = "/app/config/config.yml"
	}

	file, err := os.Open(configFileName)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to open config file %s", configFileName)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.WithError(err).Fatalf("Failed to close config file %s", configFileName)
		}
		logrus.Infof("Parsed config file %s", configFileName)
	}(file)

	var config Config

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		logrus.WithError(err).Fatalf("Failed to decode config file %s", configFileName)
		return nil, err
	}

	// same as above, but for the version file
	if versionFileName == "" {
		versionFileName = "/app/config/version.yml"
	}
	vf, err := os.Open(versionFileName)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to open version file %s", versionFileName)
		return nil, err
	}
	defer func(vf *os.File) {
		err := vf.Close()
		if err != nil {
			logrus.WithError(err).Fatalf("Failed to close version file %s", versionFileName)
		}
		logrus.Infof("Parsed version file %s", versionFileName)
	}(vf)

	// Decode the version file
	vfDecoder := yaml.NewDecoder(vf)
	if err := vfDecoder.Decode(&config); err != nil {
		logrus.WithError(err).Fatalf("Failed to decode version file %s", versionFileName)
		return nil, err
	}

	if secretsFileName == "" {
		secretsFileName = "/app/config/secrets.yml"
	}

	sf, err := os.Open(secretsFileName)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to open secrets file %s", secretsFileName)
		return nil, err
	}
	defer func(sf *os.File) {
		err := sf.Close()
		if err != nil {
			logrus.WithError(err).Fatalf("Failed to close secrets file %s", secretsFileName)
		}
		logrus.Infof("Parsed secrets file %s", secretsFileName)
	}(sf)

	secDecoder := yaml.NewDecoder(sf)
	if err := secDecoder.Decode(&config); err != nil {
		logrus.WithError(err).Fatalf("Failed to decode secrets file %s", secretsFileName)
		return nil, err
	}

	logrus.Info("Chomp chomp")
	return &config, nil
}
