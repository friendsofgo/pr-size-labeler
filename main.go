package main

import (
	"context"

	"github.com/google/go-github/v31/github"
	"github.com/kelseyhightower/envconfig"
	"github.com/posener/goaction"
	"github.com/posener/goaction/actionutil"
	"github.com/posener/goaction/log"
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

	event, err := goaction.GetPullRequest()
	if err != nil {
		log.Fatalf("Error happened while getting event info: %s", err)
	}

	prSize := GetPrSize(config, event.PullRequest)

	ctx := context.Background()
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

func GetPrSize(config Config, pr *github.PullRequest) PrSize {
	totalModifications := pr.GetAdditions() + pr.GetDeletions()

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
	GitHubToken string `envconfig:"GITHUB_TOKEN" required:"true"`
	XsMaxSize   int    `envconfig:"XS_MAX_SIZE" default:"10"`
	SMaxSize    int    `envconfig:"S_MAX_SIZE" default:"100"`
	MMaxSize    int    `envconfig:"M_MAX_SIZE" default:"500"`
	LMaxSize    int    `envconfig:"L_MAX_SIZE" default:"1000"`
	FailIfXL    bool   `envconfig:"FAIL_IF_XL" default:"false"`
	MessageIfXL string `envconfig:"MESSAGE_IF_XL" default:""`
}
