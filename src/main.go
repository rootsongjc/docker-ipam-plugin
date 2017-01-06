package main

import (
	"github.com/codegangsta/cli"
	"command"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "docker-ipam-plugin"
	app.Version = "0.1.0"
	app.Author = "rootsongjc@gmail.com"
	app.Usage = "Docker network plugin with remote IPAM"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "cluster-store", Value: "http://127.0.0.1:2379", Usage: "the key/value store endpoint url. [$CLUSTER_STORE]"},
		cli.BoolFlag{Name: "debug", Usage: "debug mode [$DEBUG]"},
	}
	app.Commands = []cli.Command{
		command.NewServerCommand(),
		command.NewIPRangeCommand(),
		command.NewReleaseIPCommand(),
		command.NewHostRangeCommand(),
		command.NewReleaseHostCommand(),
		command.NewCreateNetworkCommand(),
	}
	app.Run(os.Args)
}
