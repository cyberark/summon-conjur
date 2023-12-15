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

You can start a Docker Compose development environment by running

```sh
$ ./dev.sh
```

### Running tests

Automated CI pipelines:
- [Jenkinsfile](Jenkinsfile)

Run `./bin/test.sh`

## Releasing

### Verify and update dependencies
1. Review the changes to `go.mod` since the last release and make any needed
   updates to [NOTICES.txt](./NOTICES.txt):
   - Add any dependencies that have been added since the last tag, including
     an entry for them alphabetically under the license type (make sure you
     check the license type for the version of the project we use) and a copy
     of the copyright later in the same file.
   - Update any dependencies whose versions have changed - there are usually at
     least two version entries that need to be modified, but if the license type
     of the dependency has also changed, then you will need to remove the old
     entries and add it as if it were a new dependency.
   - Remove any dependencies we no longer include.

   If no dependencies have changed, you can move on to the next step.

### Update the version and changelog
1. Create a new branch for the version bump.
1. Based on the unreleased content, determine the new version number and update
   the [version.go](pkg/summon_conjur/version.go) file.
1. Review the [changelog](CHANGELOG.md) to make sure all relevant changes since
   the last release have been captured. You may find it helpful to look at the
   list of commits since the last release - you can find this by visiting the
   [releases page](https://github.com/cyberark/summon-conjur/releases) and
   clicking the "`N commits` to main since this release" link for the latest
   release.

   This is also a good time to make sure all entries conform to our
   [changelog guidelines](https://github.com/cyberark/community/blob/main/Conjur/CONTRIBUTING.md#changelog-guidelines).
1. Commit these changes - `Bump version to x.y.z` is an acceptable commit message - and open a PR
   for review. Your PR should include updates to `pkg/summon_conjur/version.go`,
   `CHANGELOG.md`, and if there are any license updates, to `NOTICES.txt`.

### Add a git tag
1. Once your changes have been reviewed and merged into main, tag the version
   using `git tag -s v0.1.1`. Note this requires you to be  able to sign releases.
   Consult the [github documentation on signing commits](https://help.github.com/articles/signing-commits-with-gpg/)
   on how to set this up. `vx.y.z` is an acceptable tag message.
1. Push the tag: `git push vx.y.z` (or `git push origin vx.y.z` if you are working
   from your local machine).

### Make the release public
**Note:** Until the stable quality exercises have completed, the GitHub release
should be officially marked as a `pre-release` (eg "non-production ready")

1. The tagged commit should have caused a Draft release to be created in GitHub.
   Replace the commits in the Draft release's description with the relevant entries
   from the CHANGELOG.
1. If everything else looks good, release the draft.
1. Copy the `summon-conjur.rb` homebrew formula output by goreleaser
   to the [homebrew formula for Summon-Conjur](https://github.com/cyberark/homebrew-tools/blob/main/summon-conjur.rb)
   and submit a PR to update the version of Summon-Conjur available in brew.
