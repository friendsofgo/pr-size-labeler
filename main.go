package main

import (
	"context"
	"strings"

	"github.com/google/go-github/v31/github"
	"github.com/kelseyhightower/envconfig"
	"github.com/posener/goaction"
	"github.com/posener/goaction/actionutil"
	"github.com/posener/goaction/log"
	"golang.org/x/oauth2"
)

func main() {
	if !goaction.CI {
		log.Warnf("Not in GitHub Action mode, quitting...")
		return
	}

	if goaction.Event != goaction.EventPullRequest {
		log.Debugf("Not a pull request action, nothing to do here...")
		return
	}

	var config Config
	err := envconfig.Process("input", &config)
	if err != nil {
		log.Fatalf("Required configuration is missing: %s", err)
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	event, err := goaction.GetPullRequest()
	if err != nil {
		log.Fatalf("Error happened while getting event info: %s", err)
	}

	filesToIgnore := strings.Fields(config.FilesToIgnore)

	files, _, err := client.PullRequests.ListFiles(ctx, *event.Repo.Owner.Name, *event.Repo.Name, *event.PullRequest.Number, nil)
	if err != nil {
		log.Fatalf("Error happened while getting files info: %s", err)
	}

	prSize := GetPrSize(config, files, filesToIgnore)

	gh := actionutil.NewClientWithToken(ctx, config.GitHubToken)

	_, _, err = gh.IssuesAddLabelsToIssue(ctx, goaction.PrNum(), []string{string(prSize)})
	if err != nil {
		log.Fatalf("Error happened while adding label: %s", err)
	}

	if prSize != XL {
		log.Debugf("Pull request successfully labeled")
		return
	}

	_, _, err = gh.PullRequestsCreateComment(ctx, goaction.PrNum(), &github.PullRequestComment{
		Body: github.String(config.MessageIfXL),
	})

	if err != nil {
		log.Fatalf("Error happened while adding comment: %s", err)
	}

	if config.FailIfXL {
		log.Fatalf("PR size is XL, make it shorter, please!")
	}

	log.Debugf("Pull request successfully labeled")
}

func GetPrSize(config Config, files []*github.CommitFile, filesToIgnore []string) PrSize {
	var totalModifications int

	checkIgnorableFiles := len(filesToIgnore) > 0

	for _, file := range files {
		if checkIgnorableFiles && isIgnorable(file.GetFilename(), filesToIgnore) {
			continue
		}

		totalModifications += file.GetChanges()
	}

	switch {
	case totalModifications < config.XsMaxSize:
		return XS
	case totalModifications < config.SMaxSize:
		return S
	case totalModifications < config.MMaxSize:
		return M
	case totalModifications < config.LMaxSize:
		return L
	default:
		return XL
	}
}

func isIgnorable(filename string, filesToIgnore []string) bool {
	for _, fti := range filesToIgnore {
		if strings.Contains(fti, "*") {
			ext := strings.Split(fti, "*")[1]
			if strings.Contains(filename, ext) {
				return true
			}
		}
		if fti == filename {
			return true
		}
	}

	return false
}

type PrSize string

const (
	XS PrSize = "size/xs"
	S  PrSize = "size/s"
	M  PrSize = "size/m"
	L  PrSize = "size/l"
	XL PrSize = "size/xl"
)

// Config is the data structure used to define the action settings.
type Config struct {
	GitHubToken   string `envconfig:"GITHUB_TOKEN" required:"true"`
	XsMaxSize     int    `envconfig:"XS_MAX_SIZE" default:"10"`
	SMaxSize      int    `envconfig:"S_MAX_SIZE" default:"100"`
	MMaxSize      int    `envconfig:"M_MAX_SIZE" default:"500"`
	LMaxSize      int    `envconfig:"L_MAX_SIZE" default:"1000"`
	FailIfXL      bool   `envconfig:"FAIL_IF_XL" default:"false"`
	MessageIfXL   string `envconfig:"MESSAGE_IF_XL" default:""`
	FilesToIgnore string `envconfig:"FILES_TO_IGNORE" default:""`
}
