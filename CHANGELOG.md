# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.3] - 2019-02-06
### Changed
- Go modules are now used for dependency management
- Newer goreleaser syntax is used to build artifacts
- Fixed issues with spaces in variable IDs (via conjur-api-go version increase)
- Fixed issues with homedir pathing (via conjur-api-go version increase)

## [0.5.2] - 2018-08-06
### Added
- deb and rpm packages
- homebrew formula

### Changed
- Update build/package process to use [goreleaser](https://github.com/goreleaser/goreleaser).

## [0.5.1] - 2018-07-19
### Added
- Add some logging to help debug configuration [PR #31](https://github.com/cyberark/summon-conjur/pull/31).

### Changed
- Update to the latest version of conjur-api-go.

## [0.5.0] - 2017-11-20
### Added
- Support new v5 token format and summon-conjur version flag [PR #23](https://github.com/cyberark/summon-conjur/pull/23).

## [0.4.0] - 2017-09-19
### Changed
- Support v4, https and configuration from machine identity files, [PR #20](https://github.com/cyberark/summon-conjur/pull/20).

## [0.3.0] - 2017-08-16
### Changed
- Provider updated to use [cyberark/conjur-api-go](https://github.com/cyberark/conjur-api-go). This provides compatibility with [cyberark/conjur](https://github.com/cyberark/conjur), Conjur 5 CE. PR [#13](https://github.com/cyberark/summon-conjur/pull/13).

## [0.2.0] - 2016-07-20
### Added
- `CONJUR_SSL_CERTIFICATE` can now be passed (content of cert file) [#3](https://github.com/conjurinc/summon-conjur/issues/3)
- netrc file is now only read if required [#4](https://github.com/conjurinc/summon-conjur/issues/4)
- `CONJUR_AUTHN_TOKEN` can now be used for identity [#5](https://github.com/conjurinc/summon-conjur/issues/5)

## [0.1.4] - 2016-02-29
### Fixed
- A friendly error is now returned when no argument is given [GH-2](https://github.com/conjurinc/summon-conjur/issues/2)

## [0.1.3] - 2016-02-24
### Changed
- Config now looks at `netrc_path` in conjurrc to find identity.file

## [0.1.2] - 2015-12-09
### Changed
- Config now uses env var `CONJUR_AUTHN_API_KEY` instead of `CONJUR_API_KEY`.

## [0.1.1] - 2015-10-08
### Fixed
- Fixed an issue authenticating hosts - `/` is now properly escaped.

## 0.1.0 - 2015-06-04
### Added
- Initial release

[Unreleased]: https://github.com/cyberark/summon-conjur/compare/v0.5.3...HEAD
[0.5.3]: https://github.com/cyberark/summon-conjur/compare/v0.5.2...v0.5.3
[0.5.2]: https://github.com/cyberark/summon-conjur/compare/v0.5.1...v0.5.2
[0.5.1]: https://github.com/cyberark/summon-conjur/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/cyberark/summon-conjur/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/cyberark/summon-conjur/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/cyberark/summon-conjur/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/cyberark/summon-conjur/compare/v0.1.4...v0.2.0
[0.1.4]: https://github.com/cyberark/summon-conjur/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/cyberark/summon-conjur/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/cyberark/summon-conjur/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/cyberark/summon-conjur/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/cyberark/summon-conjur/releases/tag/v0.1.0
