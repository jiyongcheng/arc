package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	gitlab "github.com/xanzy/go-gitlab"
)

type Response struct {
	*http.Response

	// These fields provide the page values for paginating through a set of
	// results.  Any or all of these may be set to the zero value for
	// responses that are not part of a paginated set, or for which there
	// are no additional pages.
	nextPage  int
	prevPage  int
	firstPage int
	lastPage  int
}

type Repository struct {
	gitUrl    string
	name      string
	namespace string
}

func getRepositories(client interface{}, service string, githubRepoType string, gitlabRepoVisibility string) ([]*Repository, error) {
	if client == nil {
		log.Fatalf("Couldn't acquire a client to talk to %s", service)
	}

	var repositories []*Repository

	if service == "github" {
		ctx := context.Background()
		options := github.RepositoryListOptions{Type: githubRepoType}

		for {
			repos, resp, err := client.(*github.Client).Repositories.List(ctx, "", &options)

			if err == nil {
				for _, repo := range repos {
					namespace := strings.Split(*repo.FullName, "/")[0]
					repositories = append(repositories, &Repository{gitUrl: *repo.GitURL, name: *repo.Name, namespace: namespace})
				}
			} else {
				return nil, err
			}
			if resp.NextPage == 0 {
				break
			}
			options.ListOptions.Page = resp.NextPage
		}
	}

	if service == "gitlab" {
		var visibility gitlab.VisibilityValue

		switch gitlabRepoVisibility {
		case "public":
			visibility = gitlab.PublicVisibility
		case "private":
			visibility = gitlab.PrivateVisibility
		case "internal":
			fallthrough
		case "default":
			visibility = gitlab.InternalVisibility
		}

		options := gitlab.ListProjectsOptions{Visibility: &visibility}

		for {
			repos, resp, err := client.(*gitlab.Client).Projects.ListProjects(&options)

			if err == nil {
				for _, repo := range repos {
					namespace := strings.Split(repo.PathWithNamespace, "/")[0]
					repositories = append(repositories, &Repository{gitUrl: repo.SSHURLToRepo, name: repo.Name, namespace: namespace})
				}
			} else {
				return nil, err
			}
			if resp.NextPage == 0 {
				break
			}
			options.ListOptions.Page = resp.NextPage
		}
	}

	return repositories, nil
}
