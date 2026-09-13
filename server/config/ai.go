package config

type AI struct {
	Enable       bool    `mapstructure:"enable" json:"enable" yaml:"enable"`
	Provider     string  `mapstructure:"provider" json:"provider" yaml:"provider"`
	BaseURL      string  `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	APIKey       string  `mapstructure:"api-key" json:"api-key" yaml:"api-key"`
	Model        string  `mapstructure:"model" json:"model" yaml:"model"`
	Temperature  float32 `mapstructure:"temperature" json:"temperature" yaml:"temperature"`
	MaxTokens    int     `mapstructure:"max-tokens" json:"max-tokens" yaml:"max-tokens"`
	ContextLimit int     `mapstructure:"context-limit" json:"context-limit" yaml:"context-limit"`
	DailyLimit   int     `mapstructure:"daily-limit" json:"daily-limit" yaml:"daily-limit"`
}
