package config

type RedisConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	DB   int    `yaml:"db"`
	Pass string `yaml:"pass"`
}

type Config struct {
	Styles []Style     `yaml:"styles"`
	Redis  RedisConfig `yaml:"redis"`
}

type Style struct {
	Name          string  `yaml:"name"`
	Image         string  `yaml:"image"`
	Title         string  `yaml:"title"`
	X             float64 `yaml:"x"`
	Y             float64 `yaml:"y"`
	ImageWidth    uint    `yaml:"imageWidth"`
	ImageHeight   uint    `yaml:"imageHeight"`
	LabelWidth    uint    `yaml:"labelWidth"`
	LabelHeight   uint    `yaml:"labelHeight"`
	TitleFontSize float64 `yaml:"titleFontSize"`
	OemFontSize   float64 `yaml:"oemFontSize"`
}

func (config *Config) GetStyle(style string) Style {
	for _, s := range config.Styles {
		if s.Name == style {
			return s
		}
	}
	return Style{}
}
