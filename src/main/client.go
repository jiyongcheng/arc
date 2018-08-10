package main

import (
	"log"
	"net/url"

	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

func newClient(service string, gitHostUrl string) interface{} {
	var gitHostUrlParsed *url.URL
	var err error

	// if a git host url has been passed in, we assume it's a gitlab installation
	if len(gitHostUrl) != 0 {
		gitHostUrlParsed, err = url.Parse(gitHostUrl)

		if err != nil {
			log.Fatalf("invalid gitlab url: %s", gitHostUrl)
		}
		api, _ := url.Parse("api/v4/")
		gitHostUrlParsed = gitHostUrlParsed.ResolveReference(api)
	}

	if service == "github" {
		// githubToken := os.Getenv("GITHUB_TOKEN")
		githubToken := "d51e1adf3a246899e3a05270c2dfd3dddf0f7259"
		if githubToken == "" {
			log.Fatal("GITHUB_TOKEN environment variable not set")
		}

		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
		tc := oauth2.NewClient(oauth2.NoContext, ts)

		client := github.NewClient(tc)
		if gitHostUrlParsed != nil {
			client.BaseURL = gitHostUrlParsed
		}

		return client
	}

	if service == "gitlab" {
		// gitlabToken := os.Getenv("GITLAB_TOKEN")
		gitlabToken := "yNsVSdXgsUzTUyxVDBav"
		if gitlabToken == "" {
			log.Fatal("GITLAB_TOKEN environment variable not set")
		}
		client := gitlab.NewClient(nil, gitlabToken)

		if gitHostUrlParsed != nil {
			client.SetBaseURL(gitHostUrlParsed.String())
		}

		return client
	}

	return nil
}
