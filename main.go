package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type query struct {
	RepositoryOwner struct {
		Repositories struct {
			TotalCount graphql.Int
			Nodes      []struct {
				Name            graphql.String
				PrimaryLanguage struct {
					Name graphql.String
				}
			}
			PageInfo struct {
				HasNextPage graphql.Boolean
				EndCursor   graphql.String
			}
		} `graphql:"repositories(first:$pageSize,after:$repositoriesCursor)"`
	} `graphql:"repositoryOwner(login:$user)"`
}

func main() {

	ctx := context.Background()
	_ = ctx
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN"),
	})

	tc := oauth2.NewClient(ctx, ts)

	ghv4Client := graphql.NewClient("https://api.github.com/graphql", tc)

	orgLogin := os.Args[1]

	variables := map[string]interface{}{
		"pageSize":           graphql.Int(100),
		"repositoriesCursor": (*graphql.String)(nil),
		"user":               graphql.String(orgLogin),
	}

	err := fetchRepoDetails(ctx, ghv4Client, variables)

	if err != nil {
		log.Fatal(err)
	}

}

func fetchRepoDetails(ctx context.Context, cli *graphql.Client,
	variables map[string]interface{}) error {
	var q query
	err := cli.Query(ctx, &q, variables)

	if err != nil {
		return err
	}

	for _, repo := range q.RepositoryOwner.Repositories.Nodes {
		fmt.Printf("%s|%s|%d\n", repo.Name,
			repo.PrimaryLanguage.Name,
			q.RepositoryOwner.Repositories.TotalCount)
	}

	if !q.RepositoryOwner.Repositories.PageInfo.HasNextPage {
		return nil
	}

	variables["repositoriesCursor"] = q.RepositoryOwner.Repositories.PageInfo.EndCursor

	return fetchRepoDetails(ctx, cli, variables)

}
