package main

import (
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v32/github"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/context"
)

var repoNameRegex = regexp.MustCompile(`^terraform-provider-(.+)$`)

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "version",
				Value: "*",
				Usage: "Semver matcher for the versions to consider. Always the highest version matching will be used.",
			},
		},
		Action: func(c *cli.Context) error {

			if c.Args().Len() == 0 {
				return errors.New("repository name must be specified")
			}

			repositoryName := c.Args().First()

			version, err := semver.NewConstraint(c.String("version"))
			if err != nil {
				return errors.Wrapf(err, "while parsing version constraint %s", c.String("version"))
			}

			parts := strings.Split(repositoryName, "/")
			if len(parts) != 2 {
				return errors.New("repository must be in the format: owner/name")
			}

			owner := parts[0]
			repo := parts[1]

			match := repoNameRegex.FindStringSubmatch(repo)

			if match == nil {
				return errors.New("repository name must start with 'terraform-provider-' ")
			}

			pluginName := match[1]

			gcl := github.NewClient(nil)

			ctx := context.Background()

			release, err := findLatestReleaseWithPluginAssets(ctx, gcl, owner, repo, version)
			if err != nil {
				return err
			}

			return release.generateShims("terraform.d", pluginName)

		},
	}

	app.RunAndExitOnError()

}

// terraform-provider-linuxbox_v0.0.13_darwin_amd64.tar.gz

var assetNamePattern = regexp.MustCompile(`^terraform-provider-[^_]+_(v[0-9]+\.[0-9]+\.[0-9]+)_([^_]+_[^_]+)(.tar.gz)$`)

func matchTerraformProviderAssetName(s string) (bool, string) {
	r := assetNamePattern.FindStringSubmatch(s)
	if r == nil {
		return false, ""
	}

	return true, r[2]
}
