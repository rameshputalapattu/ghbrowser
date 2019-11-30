package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type query struct {
	Organization struct {
		Repositories struct {
			TotalCount githubv4.Int
			Nodes      []struct {
				Name            githubv4.String
				PrimaryLanguage struct {
					Name githubv4.String
				}
			}
			PageInfo struct {
				HasNextPage githubv4.Boolean
				EndCursor   githubv4.String
			}
		} `graphql:"repositories(first:$pageSize,after:$repositoriesCursor)"`
	} `graphql:"organization(login:$user)"`
}

func main() {

	ctx := context.Background()
	_ = ctx
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
	})

	tc := oauth2.NewClient(ctx, ts)

	ghv4Client := githubv4.NewClient(tc)

	orgLogin := os.Args[1]

	variables := map[string]interface{}{
		"pageSize":           githubv4.Int(100),
		"repositoriesCursor": (*githubv4.String)(nil),
		"user":               githubv4.String(orgLogin),
	}

	err := fetchRepoDetails(ctx, ghv4Client, variables)

	if err != nil {
		log.Fatal(err)
	}

}

func fetchRepoDetails(ctx context.Context, cli *githubv4.Client,
	variables map[string]interface{}) error {
	var q query
	err := cli.Query(ctx, &q, variables)

	if err != nil {
		return err
	}

	for _, repo := range q.Organization.Repositories.Nodes {
		fmt.Printf("%s|%s|%d\n", repo.Name,
			repo.PrimaryLanguage.Name,
			q.Organization.Repositories.TotalCount)
	}

	if !q.Organization.Repositories.PageInfo.HasNextPage {
		return nil
	}

	variables["repositoriesCursor"] = q.Organization.Repositories.PageInfo.EndCursor

	return fetchRepoDetails(ctx, cli, variables)

}
