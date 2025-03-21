# Contributing

For general contribution and community guidelines, please see the [community repo](https://github.com/cyberark/community).

## Contributing

1. [Fork the project](https://help.github.com/en/github/getting-started-with-github/fork-a-repo)
2. [Clone your fork](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/cloning-a-repository)
3. Make local changes to your fork by editing files
4. [Commit your changes](https://help.github.com/en/github/managing-files-in-a-repository/adding-a-file-to-a-repository-using-the-command-line)
5. [Push your local changes to the remote server](https://help.github.com/en/github/using-git/pushing-commits-to-a-remote-repository)
6. [Create new Pull Request](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request-from-a-fork)

From here your pull request will be reviewed and once you've responded to all
feedback it will be merged into the project. Congratulations, you're a
contributor!

## Development

You can start a Docker Compose development environment by running

```sh
./dev.sh
```

### Running tests

Automated CI pipelines:

- [Jenkinsfile](Jenkinsfile)

Run `./bin/test.sh`

## Releasing

Releases should be created by maintainers only. To create and promote a
release, follow the instructions in this section.

### Verify and update dependencies

1.  Review the changes to `go.mod` since the last release and make any needed
    updates to [NOTICES.txt](./NOTICES.txt):
    *   Verify that dependencies fit into supported licenses types:
        ```shell
         go-licenses check ./... --allowed_licenses="MIT,ISC,Apache-2.0,BSD-3-Clause" \
            --ignore github.com/cyberark/conjur-summon \
            --ignore $(go list std | awk 'NR > 1 { printf(",") } { printf("%s",$0) } END { print "" }')
        ```
        If there is new dependency having unsupported license, such license should be included to [notices.tpl](./notices.tpl)
        file in order to get generated in NOTICES.txt.  

        NOTE: The second ignore flag tells the command to ignore standard library packages, which
        may or may not be necessary depending on your local Go installation and toolchain.

    *   If no errors occur, proceed to generate updated NOTICES.txt:
        ```shell
         go-licenses report ./... --template notices.tpl > NOTICES.txt \
            --ignore github.com/cyberark/summon-conjur \
            --ignore $(go list std | awk 'NR > 1 { printf(",") } { printf("%s",$0) } END { print "" }')
         ```

### Update the changelog

**NOTE:** If the Changelog and NOTICES.txt are already up-to-date, skip this
step and promote the desired build from the main branch.

1. Create a new branch for the version bump.

2. Based on the changelog content, determine the new version number and update.

3. Review the git log and ensure the [changelog](CHANGELOG.md) contains all
   relevant recent changes with references to GitHub issues or PRs, if possible.

5. Commit these changes - `Bump version to x.y.z` is an acceptable commit
   message - and open a PR for review.

### Release and Promote

1. Merging into the main branch will automatically trigger a release.
   If successful, this release can be promoted at a later time.

2. Jenkins build parameters can be utilized to promote a successful release
   or manually trigger aditional releases as needed.

3. Reference the
   [internal automated release doc](https://github.com/conjurinc/docs/blob/master/reference/infrastructure/automated_releases.md#release-and-promotion-process)
   for releasing and promoting.
