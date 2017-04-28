package main

import (
	"encoding/json"
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
			Name:   "file",
			Usage:  "write address to `FILE`",
			EnvVar: "FILE",
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

		// Create the logger for main and switch to JSON output
		log := logrus.WithField("context", "main")
		logrus.SetFormatter(&logrus.JSONFormatter{})

		if len(c.String("email")) == 0 || len(c.String("password")) == 0 {
			log.Fatal("neither email nor password may be blank")
		}

		// Generate token
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

		// Display the address and write it to file if requested
		filename := c.String("file")
		if len(filename) != 0 {
			err := func() error {
				w, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
				if err != nil {
					return err
				}
				defer w.Close()
				v := map[string]interface{}{
					"address": s.Addr(),
					"token":   token,
				}
				if err := json.NewEncoder(w).Encode(v); err != nil {
					return err
				}
				return nil
			}()
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(filename)
		} else {
			log.Infof("listening at %s...", s.Addr())
		}

		// Wait for a signal before shutting down
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	}
	app.Run(os.Args)
}
