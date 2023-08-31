package configs

import "github.com/namsral/flag"

var (
	GinMode, BaseUrl string
	IsGinInDebug     bool
)

func ConfigGin(commandSet *flag.FlagSet) {
	if command := commandSet.Lookup("gin-mode"); command != nil {
		GinMode = command.Value.String()
	}
	if command := commandSet.Lookup("base-url"); command != nil {
		BaseUrl = command.Value.String()
	}
	IsGinInDebug = GinMode == "debug"
}
