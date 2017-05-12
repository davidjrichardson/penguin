package main


import (
	"os"
	"log"

	"./src/youtubeDL"
	"./src/admin"
	"./src/playlist"
	"./src/server"
	"./src/config"
	"./src/help"
	"./src/videoplayer"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Pass config file as first argument to program")
		os.Exit(1)
	}
	config.Init(os.Args[1])
	server.Init()
	playlist.Init()
	admin.Init()
	youtubeDL.Init()
	youtubeDL.Update()
	videoplayer.Init()
	config.Destroy()

	help.PrintMasthead()
	//go videoplayer.Start()
	server.Run()
}
