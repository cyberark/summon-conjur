# Contributing

For general contribution and community guidelines, please see the [community repo](https://github.com/cyberark/community).

## Contributing

1. [Fork the project](https://help.github.com/en/github/getting-started-with-github/fork-a-repo)
2. [Clone your fork](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/cloning-a-repository)
3. Make local changes to your fork by editing files
3. [Commit your changes](https://help.github.com/en/github/managing-files-in-a-repository/adding-a-file-to-a-repository-using-the-command-line)
4. [Push your local changes to the remote server](https://help.github.com/en/github/using-git/pushing-commits-to-a-remote-repository)
5. [Create new Pull Request](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request-from-a-fork)

From here your pull request will be reviewed and once you've responded to all
feedback it will be merged into the project. Congratulations, you're a
contributor!

## Development

You can start a docker-compose development environment by running

```sh
$ ./dev.sh
```

### Dependency management With dep
[dep](https://golang.github.io/dep/docs/introduction.html) is being used to manage dependencies.

When you add a new package, or change the version of an existing package, run (in the `dev` container)

```sh
/go/src/github.com/cyberark/summon-conjur# dep ensure
```

to update `Gopkg.toml` and `Gopkg.lock`.

### Running tests

Automated CI pipelines:
- [.gitlab.ci.yml](.gitlab.ci.yml)
- [Jenkinsfile](Jenkinsfile)

Run `./test.sh oss` for OSS tests, `./test.sh enterprise` for Enterprise tests.
This defaults to both.
