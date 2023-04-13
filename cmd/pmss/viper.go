package pmss

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initViper() {
	viper.SetConfigName("pmss")               // name of config file (without extension)
	viper.AddConfigPath("/etc/pmss")          // path to look for the config file in
	viper.AddConfigPath("$HOME/.config/pmss") // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.pmss/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logrus.WithError(err).Fatal("Error reading config file")
		} else {
			logrus.Debug("No config file found")
		}
	}
	logrus.Debug("Using config file: ", viper.ConfigFileUsed())
}
