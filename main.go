package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/github/hub/git"
	"github.com/github/hub/github"
)

const (
	UmengLabel = "umeng"
)

func getToken() (token string, err error) {
	return getConfig("github.token")
}

func getConfig(key string) (value string, err error) {
	value, err = git.Config(key)
	if value != "" && err == nil {
		return
	}
	value, err = git.GlobalConfig(key)
	return
}

func getOnwerRepo() (owner string, name string, err error) {
	repo, err := github.LocalRepo()
	if err != nil {
		return
	}
	project, err := repo.CurrentProject()
	if err != nil {
		return
	}
	return project.Owner, project.Name, err
}

func handleError(err error) {
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "umeng2github"
	app.Usage = "Import umeng error data which exported from umeng website as csv file to Github issues"

	app.Commands = []cli.Command{
		{
			Name: "import",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "token, k",
					Usage:  "github personal token to use, can also set in git config with github.token",
					EnvVar: "GITHUB_TOKEN",
				},
				cli.StringFlag{
					Name:   "owner, o",
					Usage:  "owner of the repo to which the issues are created",
					EnvVar: "OWNER",
				},
				cli.StringFlag{
					Name:   "repo, r",
					Usage:  "repo name to which the issues are created",
					EnvVar: "REPO",
				},
				cli.BoolFlag{
					Name:   "short",
					Usage:  "make output short, i.e. report success only",
					EnvVar: "SHORT",
				},
			},
			Usage:       "import umeng error data csv to github",
			Description: "The only argument is the filepath",
			Action: func(c *cli.Context) {
				token := c.String("token")
				owner := c.String("owner")
				repo := c.String("repo")
				if token == "" {
					tokenFromConfig, err := getToken()
					handleError(err)
					token = tokenFromConfig
				}
				if token == "" {
					fmt.Println("Cannot find github.token in either global git config or local git config")
					os.Exit(1)
				}
				if owner == "" || repo == "" {
					ownerFromConfig, repoFromConfig, err := getOnwerRepo()
					handleError(err)
					if owner == "" {
						owner = ownerFromConfig
					}
					if repo == "" {
						repo = repoFromConfig
					}
				}

				if len(c.Args()) < 1 {
					fmt.Println("Missing file path argument")
					os.Exit(1)
				}

				client := NewClient(token, owner, repo)

				shortResultFlag := c.Bool("short")

				for _, arg := range c.Args() {
					report, err := ReadReport(arg)
					handleError(err)
					for _, record := range report {
						body, err := FormatRecord(record)
						handleError(err)

						issue, err := client.CreateIssue(record.Title, body)
						handleError(err)
						if !shortResultFlag {
							fmt.Printf("issue #%d created at %s\n", *issue.Number, *issue.HTMLURL)
						}
					}
				}

				if shortResultFlag {
					fmt.Println("success")
				}
			},
		},
	}

	app.Run(os.Args)
}
