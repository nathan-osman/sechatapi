package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-osman/sechatapi"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "sechatapi"
	app.Usage = "HTTP api for the Stack Exchange chat network"
	app.Flags = []cli.Flag{
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

		// Create the logger for main and switch to JSON output
		log := logrus.WithField("context", "main")
		logrus.SetFormatter(&logrus.JSONFormatter{})

		// Basic sanity check to ensure credentials were supplied
		if len(c.String("email")) == 0 || len(c.String("password")) == 0 {
			log.Fatal("neither email nor password may be blank")
		}

		// Generate a unique token to secure requests
		token := uuid.NewV4().String()

		// Start the server
		s, err := sechatapi.New(&sechatapi.Config{
			Email:    c.String("email"),
			Password: c.String("password"),
			Token:    token,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer s.Close()

		// Use the logger to output the address that the server is listening on
		log.WithField("address", s.Addr()).Infof("listening on %s", s.Addr())

		// Wait for a signal before shutting down
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	}
	app.Run(os.Args)
}
