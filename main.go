package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
)

type Config struct {
	GhToken string
	GhOwner string
	GhRepo  string
}

// Model
type Package struct {
	FullName      string
	Description   string
	StarsCount    int
	ForksCount    int
	LastUpdatedBy string
}

func main() {

	config := getConfig()

	ctx := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GhToken},
	)
	tokenClient := oauth2.NewClient(ctx, tokenService)

	client := github.NewClient(tokenClient)

	repo, _, err := client.Repositories.Get(ctx, config.GhOwner, config.GhRepo)

	if err != nil {
		fmt.Printf("Problem in getting repository information %v\n", err)
		os.Exit(1)
	}

	pack := &Package{
		FullName:    *repo.FullName,
		Description: *repo.Description,
		ForksCount:  *repo.ForksCount,
		StarsCount:  *repo.StargazersCount,
	}

	fmt.Printf("%+v\n", pack)

	commitInfo, _, err := client.Repositories.ListCommits(ctx, config.GhOwner, config.GhRepo, nil)

	if err != nil {
		fmt.Printf("Problem in commit information %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", commitInfo[0]) // Last commit information

	// repository readme information
	readme, _, err := client.Repositories.GetReadme(ctx, config.GhOwner, config.GhRepo, nil)
	if err != nil {
		fmt.Printf("Problem in getting readme information %v\n", err)
		return
	}

	// get content
	content, err := readme.GetContent()
	if err != nil {
		fmt.Printf("Problem in getting readme content %v\n", err)
		return
	}

	fmt.Println(content)

	// Get Rate limit information

	rateLimit, _, err := client.RateLimits(ctx)
	if err != nil {
		fmt.Printf("Problem in getting rate limit information %v\n", err)
		return
	}

	fmt.Printf("Limit: %d \nRemaining %d \n", rateLimit.Core.Limit, rateLimit.Core.Remaining) // Last commit information

}

func getConfig() Config {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
		panic("No ./config.json file found")
	}
	return config
}
