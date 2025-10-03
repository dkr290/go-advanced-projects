package main

import "model-image-deployer/cmd/api"

func main() {
	app := api.NewApp()
	app.Run()
}

