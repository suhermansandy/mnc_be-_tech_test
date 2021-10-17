package main

import (
	"mnc-be-tech-test/app"
	"mnc-be-tech-test/config"
	"runtime"
)

func main() {
	app := &app.App{}
	app.Initialize()
	app.Run(":" + config.Env.HTTPPort)
	runtime.Goexit()
}
