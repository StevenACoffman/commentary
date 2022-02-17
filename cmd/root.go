package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"github.com/StevenACoffman/commentary/pkg/github"
	"github.com/StevenACoffman/commentary/pkg/middleware"
	"github.com/StevenACoffman/commentary/pkg/types"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
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
		a := strings.Split(ownerAndRepo, "/")
		repo := a[len(a)-1]
		fmt.Println(repo)
		httpClient := middleware.NewBearerAuthHTTPClient(key)

		graphqlClient := graphql.NewClient("https://api.github.com/graphql", httpClient)

		// App Starting
		logger := log.New(os.Stdout,
			"INFO: ",
			log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("main : Started")

		ctx := context.Background()
		pr, comments, err := github.GetPullRequestAndCommentsForCommit(ctx, graphqlClient, sha, repo, owner)
		if err != nil {
			fmt.Println("ERROR", err)
		}
		commentID := filterComments(comments)
		fmt.Println("Got PR#", pr)
		prURL := fmt.Sprintf("https://github.com/%s/pull/%d", ownerAndRepo, pr)

		fmt.Println(prURL)

		now := time.Now().Format(time.RFC1123)
		newMessage := fmt.Sprintf("%s\nThe current date is %s. This comment will be updated.", Marker, now)
		if commentID != "" {
			id, err := github.UpdateComment(ctx, graphqlClient, commentID, newMessage)
			if err != nil {
				fmt.Println("ERROR", err)
			}
			fmt.Println(id)
		} else {
			id, err := github.CreateNewPullRequestComment(ctx, graphqlClient, pr.ID, newMessage)
			if err != nil {
				fmt.Println("ERROR", err)
			}
			fmt.Println(id)
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

// 	Populated during build
var (
	MarkFmt          = `<!-- {"write":"github-pr-comment-api", "v":1, "id":"%s"} -->`
	MarkActionTypeID = "default"
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
