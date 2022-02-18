package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/StevenACoffman/commentary/pkg/github"
	"github.com/StevenACoffman/commentary/pkg/middleware"
	"github.com/StevenACoffman/commentary/pkg/types"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "commentary",
	Short: "Get Open GitHub Pull Requests for you and your Team",
	Long:  `This lets you `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		defer func() {
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}()

		key := os.Getenv("GITHUB_TOKEN")
		if key == "" {
			err = fmt.Errorf("must set GITHUB_TOKEN=<github token>")
			return
		}
		sha := os.Getenv("GITHUB_SHA")
		if sha == "" {
			err = fmt.Errorf("must set GITHUB_SHA=<git commit SHA1>")
			return
		}
		owner := os.Getenv("GITHUB_REPOSITORY_OWNER")
		if owner == "" {
			err = fmt.Errorf("must set GITHUB_REPOSITORY_OWNER=<git repository owner>")
			return
		}
		ownerAndRepo := os.Getenv("GITHUB_REPOSITORY")
		if ownerAndRepo == "" {
			err = fmt.Errorf("must set GITHUB_REPOSITORY=<git repository and owner>")
			return
		}

		baseRef := os.Getenv("GITHUB_BASE_REF")
		if baseRef == "" {
			err = fmt.Errorf("must set GITHUB_BASE_REF=<main or whatever>")
			return
		}
		headRef := os.Getenv("GITHUB_HEAD_REF")
		if headRef == "" {
			err = fmt.Errorf("must set GITHUB_BASE_REF=<develop or whatever>")
			return
		}
		githubRefName := os.Getenv("GITHUB_REF_NAME")
		if githubRefName == "" {
			err = fmt.Errorf("must set GITHUB_REF_NAME=<1/merge>")
			return
		}
		prNum, err := strconv.Atoi(BeforeLastSlash(githubRefName))
		fmt.Println("unable to parse integer out of", githubRefName)

		repo := AfterLastSlash(ownerAndRepo)
		fmt.Printf(`GITHUB_SHA:%s
GITHUB_REPOSITORY_OWNER:%s
GITHUB_REPOSITORY:%s
repo:%s
`, sha, owner, ownerAndRepo, repo)
		httpClient := middleware.NewBearerAuthHTTPClient(key)

		graphqlClient := graphql.NewClient("https://api.github.com/graphql", httpClient)

		// App Starting
		logger := log.New(os.Stdout,
			"INFO: ",
			log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("main : Started")

		ctx := context.Background()
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			if strings.Contains(pair[0], "REF") {
				fmt.Println(e)
			}
		}

		//pr, comments, err := github.GetPullRequestAndCommentsForCommit(ctx, graphqlClient, sha, repo, owner)
		pr, comments, err := github.GetPullRequestByBranch(ctx, graphqlClient, owner, repo, headRef, baseRef)
		if err != nil {
			fmt.Println("ERROR", err)
		}

		fmt.Println("Found ", len(comments), " comments on this PR")
		commentID := filterComments(comments)
		if pr.Number != 0 {
			//github.GetPullRequestByURI(ctx, graphqlClient, "https://github.com/StevenACoffman/commentary/pull/1")
			//pr, comments, err = github.GetPullRequestByBranch(ctx, graphqlClient, owner, repo, headRef, baseRef)
			prNum = pr.Number
		}
		fmt.Println("Got PR#", prNum)
		prURL := fmt.Sprintf("https://github.com/%s/pull/%d", ownerAndRepo, prNum)

		fmt.Println(prURL)

		now := time.Now().Format(time.RFC1123)
		newMessage := fmt.Sprintf("%s\nThe current date is %s. This comment will be updated by %s.",
			Marker, now, MarkActionTypeID)
		if commentID != "" {
			id, err := github.UpdateComment(ctx, graphqlClient, commentID, newMessage)
			if err != nil {
				fmt.Println("ERROR", err)
				os.Exit(1)
			}
			fmt.Println(id)
		} else if pr.ID != "" {
			id, err := github.CreateNewPullRequestComment(ctx, graphqlClient, pr.ID, newMessage)
			if err != nil {
				fmt.Println("ERROR", err)
				os.Exit(1)
			}
			fmt.Println(id)
		} else {
			fmt.Println("Could not find the Pull Request")
			os.Exit(1)
		}

		if err == nil {
			logger.Println("finished clean")
			os.Exit(0)
		} else {
			logger.Printf("Got error: %v", err)
			os.Exit(1)
		}
	},
}

func AfterLastSlash(path string) (file string) {
	slash := "/"
	li := strings.LastIndex(path, slash)
	if li == -1 {
		return path
	}
	return path[li+1:]
}

func BeforeLastSlash(path string) (file string) {
	slash := "/"
	li := strings.LastIndex(path, slash)
	if li == -1 {
		return path
	}
	return path[:li]
}

var (
	MarkFmt          = `<!-- {"write":"github-pr-comment-api", "v":1, "id":"%s"} -->`
	MarkActionTypeID = getEnv("COMMENTARY_ACTION_TYPE", "default")
	Marker           = fmt.Sprintf(MarkFmt, MarkActionTypeID)
)

func filterComments(comments []types.CommentNodes) string {
	for _, comment := range comments {
		if strings.HasPrefix(comment.Body, Marker) {
			return comment.ID
		}
	}
	return ""
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default is $HOME/.commentary.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".commentary" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".commentary")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
