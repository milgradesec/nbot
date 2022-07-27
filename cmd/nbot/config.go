package main

import "github.com/spf13/viper"

func init() {
	viper.SetEnvPrefix("nbot")
	viper.AutomaticEnv()
	viper.BindEnv("AWS_ACCESS_KEY_ID", "AWS_ACCESS_KEY_ID")
	viper.BindEnv("AWS_SECRET_ACCESS_KEY", "AWS_SECRET_ACCESS_KEY")
	viper.BindEnv("AWS_DEFAULT_REGION", "AWS_DEFAULT_REGION")
}
