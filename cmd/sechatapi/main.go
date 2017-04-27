package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-osman/sechatapi"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "sechatapi"
	app.Usage = "HTTP api for the Stack Exchange chat network"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "addr",
			Usage:  "run server at `ADDR`",
			EnvVar: "ADDR",
			Value:  "127.0.0.1:0",
		},
		cli.StringFlag{
			Name:   "email",
			Usage:  "use `EMAIL` for authentication",
			EnvVar: "EMAIL",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "use `PASSWORD` for authentication",
			EnvVar: "PASSWORD",
		},
	}
	app.Action = func(c *cli.Context) {

		log := logrus.WithField("context", "main")

		if len(c.String("email")) == 0 || len(c.String("password")) == 0 {
			log.Fatal("neither email nor password may be blank")
		}

		// Start the server
		s, err := sechatapi.New(&sechatapi.Config{
			Addr:     c.String("addr"),
			Email:    c.String("email"),
			Password: c.String("password"),
		})
		if err != nil {
			log.Fatal(err)
		}
		defer s.Close()

		// Display the address
		log.Infof("listening at %s...", s.Addr())

		// Wait for a signal before shutting down
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	}
	app.Run(os.Args)
}
