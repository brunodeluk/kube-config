package source

import (
	"context"
	"fmt"
	gitclient "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"path/filepath"
)

type Git struct {
	Token  string
	URL    string
	Branch string
	Dir    string
}

func (g *Git) Fetch(ctx context.Context) (string, error) {
	fmt.Printf("Cloning %s repo...\n", g.URL)
	wd, _ := os.Getwd()
	path, err := os.MkdirTemp(wd, "repo")
	if err != nil {
		os.RemoveAll(path)
		return "", err
	}

	_, err = gitclient.PlainCloneContext(ctx, path, false, &gitclient.CloneOptions{
		URL: g.URL,
		Auth: &http.BasicAuth{
			Username: "git",
			Password: g.Token,
		},
		RemoteName:    gitclient.DefaultRemoteName,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", g.Branch)),
		SingleBranch:  true,

		NoCheckout: false,
		Progress:   nil,
	})

	fmt.Printf("finished cloning %s repo\n", g.URL)
	return filepath.Join(path, g.Dir), err
}
