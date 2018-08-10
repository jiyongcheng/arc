package main

import (
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sendgrid/sendgrid-go"
			"fmt"
		"html/template"
	"bytes"
	"github.com/andygrunwald/go-jira"
)

func helloEmail(tickets []jira.Issue) []byte {
	issues := make(map[string][]jira.Issue)
	for _,ticket :=range tickets {
		issueType := ticket.Fields.Type.Name
		issues[issueType] = append(issues[issueType], ticket)
	}

	notes := ""
	for key,ticket := range issues {
		notes = notes + key + "\n"
		for _,issue := range ticket {
			notes = notes + `<a href=https://formcorp.atlassian.net/browse/` + issue.Key + `>`+ issue.Key + `</a>`
			notes = notes + "-" + issue.Fields.Summary + "\n"
		}
	}

	address := "aaronji@transformd.com"
	name := "Example User"
	from := mail.NewEmail(name, address)
	subject := "Hello World"
	address = "aaronji@transformd.com"
	name = "Example User"
	to := mail.NewEmail(name, address)
	templateData := struct {
		Name string
		URL  string
		Notes string
	}{
		Name: "Dhanush",
		URL:  "http://geektrust.in",
		Notes:notes,
	}
	t, err := template.ParseFiles("/home/aaronji/go/arc/src/main/emailTemplate.html")
	if err != nil {
		fmt.Println(err)
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		fmt.Println(err)
	}

	content := mail.NewContent("text/html", buf.String())
	m := mail.NewV3MailInit(from, subject, to, content)
	return mail.GetRequestBody(m)
}

func sendReleaseNotesEmail(tickets []jira.Issue) {
	request := sendgrid.GetRequest("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", "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = getReleaseNotesEmailBody(tickets)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func getReleaseNotesEmailBody(tickets []jira.Issue) []byte {
	issues := make(map[string][]jira.Issue)
	for _,ticket :=range tickets {
		issueType := ticket.Fields.Type.Name
		issues[issueType] = append(issues[issueType], ticket)
	}

	notes := ""
	for key,ticket := range issues {
		notes = notes + key + "\n"
		for _,issue := range ticket {
			notes = notes + `<a href=https://formcorp.atlassian.net/browse/` + issue.Key + `>`+ issue.Key + `</a>`
			notes = notes + "-" + issue.Fields.Summary + "\n"
		}
	}

	address := "aaronji@transformd.com"
	name := "Example User"
	from := mail.NewEmail(name, address)
	subject := "Hello World"
	address = "aaronji@transformd.com"
	name = "Example User"
	to := mail.NewEmail(name, address)
	templateData := struct {
		Notes string
	}{
		Notes:notes,
	}
	t, err := template.ParseFiles("/home/aaronji/go/arc/src/main/emailTemplate.html")
	if err != nil {
		fmt.Println(err)
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		fmt.Println(err)
	}

	content := mail.NewContent("text/html", buf.String())
	m := mail.NewV3MailInit(from, subject, to, content)
	return mail.GetRequestBody(m)
}
