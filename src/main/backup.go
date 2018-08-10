package main

import (
	"log"
	"net/url"
	"os/exec"
	"path"
	"sync"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

// we have them here so that we can override these in the tests

var execCommand = exec.Command
var appFS = afero.NewOsFs()
var gitCommand = "git"

//check if we have a copy of the repo already, if
//we do , we update the repo, else we do a fresh clone
func backUp(backupDir string, repo *Repository, wg *sync.WaitGroup) ([]byte, error) {
	defer wg.Done()

	repoDir := path.Join(backupDir, repo.namespace, repo.name)
	_, err := appFS.Stat(repoDir)

	var stdoutStderr []byte
	if err == nil {
		log.Printf("%s exists, updating.\n", repo.name)
		cmd := execCommand(gitCommand, "-C", repoDir, "pull")
		stdoutStderr, err = cmd.CombinedOutput()
	} else {
		log.Printf("Cloning %s \n", repo.name)
		cmd := execCommand(gitCommand, "clone", repo.gitUrl, repoDir)
		stdoutStderr, err = cmd.CombinedOutput()
	}

	return stdoutStderr, err
}

func setupBackupDir(backupDir string, service string, githostUrl string) string {
	if len(backupDir) == 0 {
		homeDir, err := homedir.Dir()
		if err == nil {
			service = service + ".com"
			backupDir = path.Join(homeDir, ".gitbackup", service)
		} else {
			log.Fatal("could not determine home directory and backup directory not specified")
		}
	} else {
		if len(githostUrl) == 0 {
			if service == "gitlab" {
				service = "https://git-scm.co"
			} else {
				service = service + ".com"
			}

			backupDir = path.Join(backupDir, service)

		} else {
			u, err := url.Parse(githostUrl)
			if err != nil {
				panic(err)
			}
			backupDir = path.Join(backupDir, u.Host)

		}
	}

	_, err := appFS.Stat(backupDir)
	if err != nil {
		log.Printf("%s doesn't exist , creating it \n", backupDir)
		err := appFS.MkdirAll(backupDir, 0771)
		if err != nil {
			log.Fatal(err)
		}
	}
	return backupDir
}
