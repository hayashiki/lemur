package github

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/src-d/go-billy.v4/memfs"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

//go:generate mockgen -source github.go -destination mock/mock_github.go
type Client interface {
	PushNewFileToBranch(commit *GitCommit) error
	CreateNewPullRequest(commit *GitCommit) error
}

type client struct {
	gCli *github.Client
	Github
}

type Github struct {
	Organization string
	Repo         string
	Token        string
	User         string
}

type File struct {
	Name    string
	Path    string
	Content []byte
}

func (f File) FullPath() string {
	return f.Path + f.Name
}

type GitCommit struct {
	CommitAuthorName   string
	CommitAuthorEmail  string
	Branch             string
	FileName           string
	FileContent        string
	Files              []*File
	CommitMessage      string
	PullRequestMessage string
	PullRequestTitle   string
}

func NewGitCommit(files []*File, message string) *GitCommit {
	return &GitCommit{
		// TODO: set config
		CommitAuthorName:   "hayashiki",
		CommitAuthorEmail:  os.Getenv("EMAIL"),
		Files:              files,
		Branch:             strconv.FormatInt(time.Now().UnixNano(), 10),
		CommitMessage:      message,
		PullRequestMessage: message,
		PullRequestTitle:   message,
	}
}

func NewClient(org, repo, token, user string) Client {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	cli := github.NewClient(tc)

	return &client{
		cli,
		Github{
			org,
			repo,
			token,
			user,
		},
	}
}

func (c *client) PushNewFileToBranch(commit *GitCommit) error {
	f := memfs.New()

	url := fmt.Sprintf("https://%s:%s@github.com/%s/%s.git", c.User, c.Token, c.Organization, c.Repo)

	repo, err := git.Clone(memory.NewStorage(), f, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName("refs/heads/master"),
	})
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	err = w.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: plumbing.NewBranchReferenceName(commit.Branch),
	})
	if err != nil {
		return err
	}

	for _, file := range commit.Files {
		log.Printf("File name is %s", file.FullPath())

		// rangeでdefer i.Closeってどうすればいいんだっけ？
		i, err := f.OpenFile(file.FullPath(), os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Printf("file.OpenFile err: %v", file.Content)
			return err
		}

		if _, err = i.Write(file.Content); err != nil {
			log.Printf("git.Writes err: %v", file.Content)
			return err
		}

		if _, err := w.Add(file.FullPath()); err != nil {
			log.Printf("git.Add err: %v", file.Name)
			return err
		}
	}

	ref := plumbing.ReferenceName(commit.Branch)

	if err == nil {
		hash, err := w.Commit(commit.CommitMessage, &git.CommitOptions{
			Author: &object.Signature{
				Name:  commit.CommitAuthorName,
				Email: commit.CommitAuthorEmail,
				When:  time.Now(),
			},
		})

		if err != nil {
			return err
		}

		repo.Storer.SetReference(plumbing.NewReferenceFromStrings(commit.Branch, hash.String()))
	}

	originRefSpec := fmt.Sprintf("refs/heads/%s", commit.Branch)
	remote, err := repo.Remote("origin")
	if err == nil {
		err = remote.Push(&git.PushOptions{
			Progress: os.Stdout,
			RefSpecs: []config.RefSpec{
				config.RefSpec(ref + ":" + plumbing.ReferenceName(originRefSpec)),
			},
		})
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *client) CreateNewPullRequest(commit *GitCommit) error {

	newPR := &github.NewPullRequest{
		Title:               github.String(commit.PullRequestTitle),
		Head:                github.String(commit.Branch),
		Base:                github.String("master"),
		Body:                github.String(commit.PullRequestMessage),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := c.gCli.PullRequests.Create(context.Background(), c.Organization, c.Repo, newPR)
	if err != nil {
		return err
	}

	// TODO 同時タイミングでPRをおくるとエラーになるのでブランチプッシュだけにしておく
	//o := &github.PullRequestOptions{
	//	SHA:         pr.GetHead().GetSHA(),
	//	MergeMethod: "rebase",
	//}
	//
	//_, _, err = c.gCli.PullRequests.Merge(context.Background(), c.Organization, c.Repo, pr.GetNumber(), "Mereged!", o)
	//
	//if err != nil {
	//	return err
	//}

	log.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}
