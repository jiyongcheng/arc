package main

import (
	"bytes"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"os"
	"path/filepath"
	"log"
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	SENDGRID_API_KEY string
	EMAIL_END_POINT string
	EMAIL_HOST string
	EMAIL_REQUEST_METHOD string
	EMAIL_FROM string
	EMAIL_TO []string
	EMAIL_TEMPLATE_PATH string
	EMAIL_SUBJECT string
	JIRA_URL string
}

func sendReleaseNotesEmail(tickets []jira.Issue) []byte {
	cwd, _ := os.Getwd()
	configuration := Configuration{}
	err := gonfig.GetConf(filepath.Join(cwd, "./config/development.json"), &configuration)
	ReleaseVersion := "v2018.51"
	Subject := fmt.Sprintf(configuration.EMAIL_SUBJECT, ReleaseVersion)
	issues := make(map[string][]jira.Issue)
	for _, ticket := range tickets {
		issueType := ticket.Fields.Type.Name
		issues[issueType] = append(issues[issueType], ticket)
	}

	m := mail.NewV3Mail()

	from := mail.NewEmail("", configuration.EMAIL_FROM)
	m.SetFrom(from)

	templateData := struct {
		Notes   map[string][]jira.Issue
		JiraUrl string
		Version string
	}{
		Notes:   issues,
		JiraUrl: configuration.JIRA_URL,
		Version: ReleaseVersion,
	}

	path := filepath.Join(cwd, configuration.EMAIL_TEMPLATE_PATH)

	t, err := template.ParseFiles(path)
	if err != nil {
		fmt.Println(err)
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		fmt.Println(err)
	}

	content := mail.NewContent("text/html", buf.String())
	m.AddContent(content)

	personalization := mail.NewPersonalization()

	for _,to := range configuration.EMAIL_TO {
		personalization.AddTos(mail.NewEmail("", to))
	}

	personalization.Subject = Subject

	m.AddPersonalizations(personalization)

	request := sendgrid.GetRequest(configuration.SENDGRID_API_KEY, configuration.EMAIL_END_POINT, configuration.EMAIL_HOST)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Sprintf("Email sent status: %d", response.StatusCode)
	}

	return nil
}
