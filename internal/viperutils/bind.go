package viperutils

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func BindPFlags(flags *pflag.FlagSet, replacer *strings.Replacer) (err error) {
	flags.VisitAll(func(flag *pflag.Flag) {
		key := flag.Name
		if replacer != nil {
			key = replacer.Replace(key)
		}
		if err = viper.BindPFlag(key, flag); err != nil {
			return
		}
	})
	return nil
}
