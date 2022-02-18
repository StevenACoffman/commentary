### Commentary - Update a github comment

This is a small demo that will either create a new or update an existing comment
on a pull request in GitHub.

I got the idea from [Ben Limmer](https://benlimmer.com/2021/12/20/create-or-update-pr-comment/), but I did it in Go.

This seemed like a good way to test how fast the various methods of running github actions would be.

So writing GitHub actions in Go, I'm aware of 4 possibilities:
1. just shell out and `go run main.go`
2. [package the Go using npm](https://github.com/sanathkr/go-npm) or [like this](https://blog.xendit.engineer/how-we-repurposed-npm-to-publish-and-distribute-our-go-binaries-for-internal-cli-23981b80911b) (private or public npm registry as you please)
3. [package your Go as a docker container](https://www.sethvargo.com/writing-github-actions-in-go/) (private or public registry)
4. [attach pre-built Go artifacts to a github release and run those using js wrappers](https://full-stack.blend.com/how-we-write-github-actions-in-go.html)

I have started to add all these to this repository to see how they perform.

### Running as a GitHub Action
There are several environment variables that this needs.
+ `COMMENTARY_ACTION_TYPE` -  you can have multiple actions all racing without stepping on each other
+ `GITHUB_TOKEN` - This should be a secret, but is the personal access token of the service account (or your real github account)
+ `GITHUB_REPOSITORY` - Set by GitHub as an `owner/repo`
+ `GITHUB_REPOSITORY_OWNER` - Set by GitHub as `owner` 
+ `GITHUB_SHA` - Set by GitHub as the commit sha1, and used to look up the PR.


### Mage

Instead of `make` and `Makefile`, I used [mage](https://magefile.org/) and made a [magefile](https://github.com/StevenACoffman/teamboard/blob/main/magefile.go).

If you do `brew install mage` then you can run here:
+ `mage -v run` - will run the webserver by doing `go run main.go`
+ `mage generate` - will re-generate the genqlient code by doing `go generate ./...`
+ `mage install` - will build and install the commentary application binary
+ `mage release` - will generate a new release

Or just run the go commands by hand.