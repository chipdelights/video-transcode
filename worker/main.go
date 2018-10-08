package main

import (
	"../common"
)

func main() {

	server, err := common.MachineryInstance()
	if err != nil {
		panic("Could not create server")
	}

	server.RegisterTask("Transcode", Transcode)

	worker := server.NewWorker("worker-1", 10)
	err = worker.Launch()
	if err != nil {
		panic("Could not launch worker")
	}

}
