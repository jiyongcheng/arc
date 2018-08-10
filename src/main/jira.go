package main

import (
	"github.com/andygrunwald/go-jira"
	"fmt"
	"strings"
)

func getJiraClient()(client *jira.Client) {

	jiraURL := "https://formcorp.atlassian.net"

	tp := jira.BasicAuthTransport{
		Username:"aaronji@formcorp.com.au",
		Password:"mDnmsedT7bBQnX25a8Uc9E64",
	}

	client, err := jira.NewClient(tp.Client(), strings.TrimSpace(jiraURL))

	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	return client
}

func getTicketInformation(client *jira.Client, issueKey string){
	issue, _, _ := client.Issue.Get(issueKey, nil)


	fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
	fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
	fmt.Println(issue.Fields.Epic)
}

func test() {
	client := getJiraClient()
	getTicketInformation(client, "FPD-2450")
}
