package model

type GUIConfig struct {
	Title  string `yaml:"title"`
	Width  int    `yaml:"width"`
	Height int    `yaml:"height"`
}

type MoverConfig struct {
	Interval int `yaml:"interval"`
}

type AppConfig struct {
	GUI   GUIConfig   `yaml:"gui"`
	Mover MoverConfig `yaml:"mover"`
}
