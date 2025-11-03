# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.9.2] - 2025-11-03

### Changed
- Upgrade Go to 1.25 (CONJSE-2067)

## [0.9.1] - 2025-10-28

### Added
- Added `close-stale.yml` GitHub workflow

## [0.9.0] - 2025-10-27

### Added
- Added support for authn-iam, authn-gcp, and authn-azure authenticators. (CNJR-11056)

## [0.8.2] - 2025-09-15

### Changed
- Updated documentation to align with Conjur Enterprise name change to Secrets Manager. (CNJR-10977)

## [0.8.1] - 2025-07-23

### Changed
- Upgrade Go to 1.24.x
- Update Go dependencies to reflect conjur-api-go v0.13.2
- Improve error handling for missing .netrc file. ([cyberark/summon-conjur#83](https://github.com/cyberark/summon-conjur/issues/83), CNJR-10190)

### Fixed
- Fix inconsistent behavior when fetching large numbers of variables
  ([cyberark/summon#258](https://github.com/cyberark/summon/issues/258), CNJR-10325)

## [0.8.0] - 2024-06-06

### Changed
- Updated provider to stream secrets instead and leverage new Summon API (CNJR-4814)
- Upgraded Go to 1.22.4

## [0.7.2] - 2024-03-21

### Security
- Upgrade Go to 1.22 (CONJSE-1842)

## [0.7.1] - 2023-06-14

### Security
- Update golang.org/x/sys to v0.8.0, gopkg.in/yaml.v3 to v3.0.1, and Go to 1.20
  in Dockerfile.text
  [cyberark/summon-conjur#112](https://github.com/cyberark/summon-conjur/pull/112)

## [0.7.0] - 2023-03-10
### Added
- Added support for Conjur's OIDC and LDAP authenticators
  [cyberark/summon-conjur#108](https://github.com/cyberark/summon-conjur/pull/108)

### Changed
- Updated Golang to 1.19
  [cyberark/summon-conjur#108](https://github.com/cyberark/summon-conjur/pull/108)

### Security
- Update golang.org/x/sys to v0.1.0 for CVE-2022-29526 (not vulnerable)
  [cyberark/summon-conjur#110](https://github.com/cyberark/summon-conjur/pull/110)

### Removed
- Removed support for Conjur v4
  [cyberark/summon-conjur#108](https://github.com/cyberark/summon-conjur/pull/108)

## [0.6.4] - 2022-07-06
### Changed
- Updated direct dependencies (github.com/cyberark/conjur-api-go -> v0.10.1,
  github.com/stretchr/testify -> 1.7.2)
  [cyberark/summon-conjur#106](https://github.com/cyberark/summon-conjur/pull/106)

## [0.6.3] - 2022-05-19
### Changed
- Updated the Conjur API to 0.10.0 to support the new `CONJUR_AUTHN_JWT_HOST_ID` environment variable
  [cyberark/summon-conjur#103](https://github.com/cyberark/summon-conjur/pull/103/)

### Security
- Update test env Golang to 1.17 to fix CVE-2022-0778 and CVE-2022-1292.
  [cyberark/summon-conjur#102](https://github.com/cyberark/summon-conjur/pull/102/)

## [0.6.2] - 2022-02-25
### Changed
- Updated Conjur API to 0.9.0 to support authn-JWT
  [cyberark/summon-conjur#99](https://github.com/cyberark/summon-conjur/pull/99/)

## [0.6.1] - 2021-12-31
### Changed
- Updated Golang to 1.17 and the Conjur API to 0.8.1 
  [cyberark/summon-conjur#96](https://github.com/cyberark/summon-conjur/pull/96/)

## [0.6.0] - 2021-08-11
### Added
- Build for Apple M1 silicon.
  [cyberark/summon-conjur#88](https://github.com/cyberark/summon-conjur/issues/88)

## [0.5.5] - 2021-06-01
### Security
- Update golang.org/x/crypto to address CVE-2020-29652.
  [PR cyberark/summon-conjur#84](https://github.com/cyberark/summon-conjur/pull/84)

## [0.5.4] - 2021-03-16
### Added
- Update conjur-api-go dependency to v0.7.1.
- Preliminary support for building Solaris binaries.
  [cyberark/summon-conjur#67](https://github.com/cyberark/summon-conjur/issues/67)

### Fixed
- Verbose debug output with the -v flag, silently lost in v0.5.3 due to changes
  to the logging interface in
  [conjur-api-go](https://github.com/cyberark/conjur-api-go), is reintroduced.
  [cyberark/summon-conjur#77](https://github.com/cyberark/summon-conjur/issues/77)

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

[Unreleased]: https://github.com/cyberark/summon-conjur/compare/v0.9.1...HEAD
[0.9.1]: https://github.com/cyberark/summon-conjur/compare/v0.9.0...v0.9.1
[0.9.0]: https://github.com/cyberark/summon-conjur/compare/v0.8.1...v0.9.0
[0.8.1]: https://github.com/cyberark/summon-conjur/compare/v0.8.0...v0.8.1
[0.8.0]: https://github.com/cyberark/summon-conjur/compare/v0.7.2...v0.8.0
[0.7.2]: https://github.com/cyberark/summon-conjur/compare/v0.7.1...v0.7.2
[0.7.1]: https://github.com/cyberark/summon-conjur/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/cyberark/summon-conjur/compare/v0.6.4...v0.7.0
[0.6.4]: https://github.com/cyberark/summon-conjur/compare/v0.6.3...v0.6.4
[0.6.3]: https://github.com/cyberark/summon-conjur/compare/v0.6.2...v0.6.3
[0.6.2]: https://github.com/cyberark/summon-conjur/compare/v0.6.1...v0.6.2
[0.6.1]: https://github.com/cyberark/summon-conjur/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/cyberark/summon-conjur/compare/v0.5.5...v0.6.0
[0.5.5]: https://github.com/cyberark/summon-conjur/compare/v0.5.4...v0.5.5
[0.5.4]: https://github.com/cyberark/summon-conjur/compare/v0.5.3...v0.5.4
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
