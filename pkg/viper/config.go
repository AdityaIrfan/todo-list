package viper

import "github.com/spf13/viper"

type EnvConfig struct {
	FileName string
	FileType string
	Path     string
}

func (e *EnvConfig) ReadConfig() error {
	viper.SetConfigName(e.FileName) // name of config file (without extension)
	viper.SetConfigType(e.FileType) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(e.Path)     // path to look for the config file in
	viper.AutomaticEnv()
	viper.WatchConfig()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return err
	}
	return nil
}
