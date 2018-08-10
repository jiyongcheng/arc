package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
	"fmt"
	"github.com/araddon/dateparse"
	"time"
	"regexp"
	"github.com/andygrunwald/go-jira"
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

		type JiraIssue struct {
			summary string
			key string
			epic string
			issueType string
			priority string
		}

		//repo, resp, err := client.(*gitlab.Client).Projects.GetProject("formcorp/web")
		//repo, resp, err := client.(*gitlab.Client).Projects.GetProjectEvents("formcorp/web")
		listCommitsOptions := gitlab.ListCommitsOptions{}
		//listCommitsOptions.Since = time.Parse(time.RFC822, "01 Jan 15 10:00 UTC")

		jiraClient := getJiraClient()
		from,to := "2018-08-01","2018-08-10"
		loc, err := time.LoadLocation("Local")
		timeFrom, err := dateparse.ParseIn(from,loc)
		timeTo, err := dateparse.ParseIn(to,loc)

		refName := "staging"
		listCommitsOptions.Since = &timeFrom
		listCommitsOptions.Until = &timeTo
		listCommitsOptions.RefName = &refName
		statuses, _, err := client.(*gitlab.Client).Commits.ListCommits("formcorp/web",  &listCommitsOptions)
		if err == nil {

			var tickets []string
			var ticketsInfo []jira.Issue
			for _, commit := range statuses {
				r := regexp.MustCompile(`(?i)(FPD)-\d+`)
				ticket := r.FindString( commit.Title)
				if ticket != "" && !stringInSlice(ticket, tickets) {
					tickets = append(tickets, ticket)
					issue, _, _ := jiraClient.Issue.Get(ticket, nil)
					ticketsInfo= append(ticketsInfo, *issue)
				}
			}
			fmt.Println(len(tickets))
			fmt.Println(tickets)
			fmt.Println(ticketsInfo)
			//for key,item := range ticketsInfo {
			//	fmt.Println(key)
			//	fmt.Println(item.Fields.Summary)
			//}

			sendReleaseNotesEmail(ticketsInfo)

		} else {
			fmt.Println(err)
		}


		return repositories, nil;

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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
