package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/satori/go.uuid"
	"github.com/urfave/cli"
)

var version = "v0.1.0"

func main() {
	var token string
	var repo string
	var file string
	var path string

	app := cli.NewApp()
	app.Name = "media-to-github"
	app.Version = version
	app.HelpName = "media-to-github"
	app.Usage = "上传图片到github"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "t, token",
			Usage:       "github token",
			Destination: &token,
		},
		cli.StringFlag{
			Name:        "r, repo",
			Usage:       "github repo",
			Destination: &repo,
		},
		cli.StringFlag{
			Name:        "f, file",
			Usage:       "file path or url",
			Destination: &file,
		},
		cli.StringFlag{
			Name:        "p, path",
			Usage:       "where file store",
			Destination: &path,
		},
	}

	app.Action = func(c *cli.Context) error {
		if token == "" || repo == "" || file == "" {
			return cli.ShowAppHelp(c)
		}
		if path == "" {
			path = fmt.Sprintf("img/%s.png", uuid.NewV4().String())
		}

		var content []byte
		var err error

		if strings.HasPrefix(file, "http") {
			resp, err := http.Get(file)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			content, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
		} else {
			content, err = ioutil.ReadFile(file)
			if err != nil {
				return err
			}
		}

		return upload(repo, token, path, content)
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
