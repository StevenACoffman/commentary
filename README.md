### Commentary - Update a GitHub comment

This is a small demo that will either create a new comment or update an existing comment
on a pull request in GitHub.

I got the idea from [Ben Limmer](https://benlimmer.com/2021/12/20/create-or-update-pr-comment/), but I did it in Go.

This seemed like a good way to test how fast the various methods of running GitHub actions authored in Go would be.
My theory was that npm had an unfair advantage as it was already baked into the standard runner OS, but (SPOILER) I was wrong!

### GitHub Actions - Many Ways to Go!

Although a GitHub Action Ubuntu Runner comes with scripting languages like Bash, NodeJS, Python, etc.
it doesn't come with Go 😞. As a result, achieving quick feedback from GitHub Actions _written_ in Go, is all about
minimizing data transfer size in order to execute your written-in-Go step.

So to execute GitHub Action steps _written_ in Go, I'm aware of these possibilities:
1. Setup Go in the action, and then just shell out and `go run main.go`
2. [package the Go using npm](https://github.com/sanathkr/go-npm) and further tweaked [like this](https://blog.xendit.engineer/how-we-repurposed-npm-to-publish-and-distribute-our-go-binaries-for-internal-cli-23981b80911b) (private or public npm registry as you please)
3. [package your Go as a docker container](https://www.sethvargo.com/writing-github-actions-in-go/) (private or public registry)
4. attach pre-built Go artifacts to a GitHub release and YOLO it (download and execute artifact via shell)
5. [attach pre-built Go artifacts to a GitHub release and run those using JS wrappers](https://full-stack.blend.com/how-we-write-github-actions-in-go.html)

I have started to add all these methods to this demo repository to see how they perform. So far I'm done with the first 4, but I imagine the Node wrapper would be the
same as the YOLO Bash one if the node script has no need to install npm dependencies.

| Which Runner    | Elapsed | Download Size |
|-----------------|---------|---------------|
| YOLOCommenter   | 3s      | 2.92 MB       |
| DockerCommenter | 4s      | 10MB          |
| NodeCommenter   | 6s      | 41 MB         |
| GoRunCommenter  | 25s     | 134 MB        |

The overhead of having to download a Go environment makes `go run main.go` the slowest, but if your action needs the Go tools anyway,
maybe that doesn't add anything. You can also use the `actions/cache@v2` to further reduce startup time. It is certainly the simplest and easiest to maintain.

### Running as a GitHub Action
There are several environment variables that this needs.
+ `COMMENTARY_ACTION_TYPE` -  you can have multiple actions all racing without stepping on each other
+ `GITHUB_TOKEN` - This should be a secret, but is the personal access token of the service account (or your real GitHub account)
+ `GITHUB_REPOSITORY` - Set by GitHub as an `owner/repo`
+ `GITHUB_REPOSITORY_OWNER` - Set by GitHub as `owner` 
+ `GITHUB_BASE_REF` - Set by GitHub as `main`
+ `GITHUB_HEAD_REF` - Set by Github to be the branch name, e.g. `mybranch`
+ `GITHUB_REF_NAME` - Set by Github, for example, `1/merge`
+ `GITHUB_SHA` - Set by GitHub as the commit sha1, and used to look up the PR.

### Publishing your Go to NPM
After you (er... me?) cut a new release, you (I?)  need to increment the package.json `version` to match,
and then run `npm publish`. If you get a 404, you might need to `npm login` again.

This Go-published-to-NPM idea was pioneered by [Sanath Kumar Ramesh](https://github.com/sanathkr/go-npm) and further tweaked [like this](https://blog.xendit.engineer/how-we-repurposed-npm-to-publish-and-distribute-our-go-binaries-for-internal-cli-23981b80911b).

My silly npm package is [commentary-cli](https://www.npmjs.com/package/commentary-cli) has a 39.44 MiB tar ball. 

### Publishing to Docker
After you cut a new release, you need to push to Docker Hub. There's a handy ko.sh script
that will do it for you (er... me?) until I Go-ify it into the magefile.

This uses the excellent [`google/ko`](https://github.com/google/ko) to make a fairly small 10 MB image [over here](https://hub.docker.com/repository/docker/stevenacoffman/commentary)

### YOLO style Run from release binary
The installation instructions for a lot of things (e.g. homebrew) have you download and execute a shell script, which
seems wildly insecure, so I refer to it as "YOLO" style installation ("You Only Live Once"). 

There's a script in here `./download-extract-execute-artifact.sh` that does this.

I don't know how to execute it in a one liner, but I imagine there's a trick to do it

### Mage

Instead of `make` and `Makefile`, I used [mage](https://magefile.org/) and made a [magefile](https://github.com/StevenACoffman/teamboard/blob/main/magefile.go).

If you do `brew install mage` then you can run here:
+ `mage -v run` - will run the webserver by doing `go run main.go`
+ `mage generate` - will re-generate the genqlient code by doing `go generate ./...`
+ `mage install` - will build and install the commentary application binary
+ `mage release <v0.?.0>` - will generate a new release using Go releaser

Or just run the go commands by hand.

### Stuff inside
I use [genqlient](https://github.com/Khan/genqlient) to communicate with the GitHub GraphQL API.
I used [Cobra](https://github.com/spf13/cobra) to scaffold out the cli.
I use Mage for task automation.

### GitHub Action details

https://github.com/embano1/ci-demo-app/blob/main/DETAILS.md
