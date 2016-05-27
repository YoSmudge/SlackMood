package config

type config struct{
  SlackToken    string          `yaml:"slack_token"`
  Db            string          `yaml:"db_path"`
}
