# summon-conjur

Conjur provider for [Summon](https://github.com/cyberark/summon).

[![GitHub release](https://img.shields.io/github/release/cyberark/summon-conjur.svg)](https://github.com/cyberark/summon-conjur/releases/latest)
[![pipeline status](https://gitlab.com/cyberark/summon-conjur/badges/master/pipeline.svg)](https://gitlab.com/cyberark/summon-conjur/pipelines)

[![Github commits (since latest release)](https://img.shields.io/github/commits-since/cyberark/summon-conjur/latest.svg)](https://github.com/cyberark/summon-conjur/commits/master)

---

**Note** Use the [summon-conjurcli](https://github.com/conjurinc/summon-conjurcli) provider if you are on Conjur v4.4.0 or earlier.

**Note** You **must** set environment variable `CONJUR_MAJOR_VERSION=4` for this provider to work with Conjur v4.9.

## Install

Pre-built binaries and packages are available from GitHub releases
[here](https://github.com/cyberark/summon-conjur/releases).

### Using summon-conjur with Conjur OSS 

Are you using this project with [Conjur OSS](https://github.com/cyberark/conjur)? Then we 
**strongly** recommend choosing the version of this project to use from the latest [Conjur OSS 
suite release](https://docs.conjur.org/Latest/en/Content/Overview/Conjur-OSS-Suite-Overview.html). 
Conjur maintainers perform additional testing on the suite release versions to ensure 
compatibility. When possible, upgrade your Conjur version to match the 
[latest suite release](https://docs.conjur.org/Latest/en/Content/ReleaseNotes/ConjurOSS-suite-RN.htm); 
when using integrations, choose the latest suite release that matches your Conjur version. For any 
questions, please contact us on [Discourse](https://discuss.cyberarkcommons.org/c/conjur/5).

### Homebrew

```
brew tap cyberark/tools
brew install summon-conjur
```

### Linux (Debian and Red Hat flavors)

`deb` and `rpm` files are attached to new releases.
These can be installed with `dpkg -i summon-conjur_*.deb` and
`rpm -ivh summon-conjur_*.rpm`, respectively.

### Auto Install

**Note** Check the release notes and select an appropriate release to ensure support for your version of Conjur.

Use the auto-install script. This will install the latest version of summon-conjur.
The script requires sudo to place summon-conjur in dir `/usr/local/lib/summon`.

```
curl -sSL https://raw.githubusercontent.com/cyberark/summon-conjur/master/install.sh | bash
```

### Manual Install
Otherwise, download the [latest release](https://github.com/cyberark/summon-conjur/releases) and extract it to the directory `/usr/local/lib/summon`.

## Usage in isolation

Give summon-conjur a variable name and it will fetch it for you and print the value to stdout.

```sh-session
$ # export CONJUR_MAJOR_VERSION=4 for Conjur v4.9
$ summon-conjur prod/aws/iam/user/robot/access_key_id
8h9psadf89sdahfp98
```

### Flags

```
Usage of summon-conjur:
  -h, --help
	show help (default: false)
  -V, --version
	show version (default: false)
  -v, --verbose
	be verbose (default: false)
```

## Usage as a provider for Summon

[Summon](https://github.com/cyberark/summon/) is a command-line tool that reads a file in secrets.yml format and injects secrets as environment variables into any process. Once the process exits, the secrets are gone.

*Example*

As an example let's use the `env` command:

Following installation, define your keys in a `secrets.yml` file

```yml
AWS_ACCESS_KEY_ID: !var aws/iam/user/robot/access_key_id
AWS_SECRET_ACCESS_KEY: !var aws/iam/user/robot/secret_access_key
```

By default, summon will look for `secrets.yml` in the directory it is called from and export the secret values to the environment of the command it wraps.

Wrap the `env` in summon:

```sh
$ # export CONJUR_MAJOR_VERSION=4 for Conjur v4.9
$ summon --provider summon-conjur env
...
AWS_ACCESS_KEY_ID=AKIAJS34242K1123J3K43
AWS_SECRET_ACCESS_KEY=A23MSKSKSJASHDIWM
...
```

`summon` resolves the entries in secrets.yml with the conjur provider and makes the secret values available to the environment of the command `env`.

## Configuration

This provider uses the same configuration pattern as the [Conjur CLI
Client](https://github.com/conjurinc/api-ruby#configuration) to connect to Conjur.
Specifically, it loads configuration from:

 * `.conjurrc` files, located in the home and current directories, or at the
    path specified by the `CONJURRC` environment variable.
 * Read `/etc/conjur.conf` as a `.conjurrc` file.
 * Read `/etc/conjur.identity` as a `netrc` file. Note that the user running must either be in the group `conjur` or root to read the identity file.
 * Environment variables:
   * Version
     * `CONJUR_MAJOR_VERSION` - must be set to `4` in order for summon-conjur to work with Conjur v4.9.
   * Appliance URLs
     * `CONJUR_APPLIANCE_URL`
     * `CONJUR_CORE_URL`
     * `CONJUR_AUTHN_URL`
   * SSL certificate
     * `CONJUR_CERT_FILE`
     * `CONJUR_SSL_CERTIFICATE`
   * Authentication
     * Account
       * `CONJUR_ACCOUNT`
     * Login
       * `CONJUR_AUTHN_LOGIN`
       * `CONJUR_AUTHN_API_KEY`
     * Token
       * `CONJUR_AUTHN_TOKEN`
       * `CONJUR_AUTHN_TOKEN_FILE`

If `CONJUR_AUTHN_LOGIN` and `CONJUR_AUTHN_API_KEY` or `CONJUR_AUTHN_TOKEN` or `CONJUR_AUTHN_TOKEN_FILE` are not provided, the username and API key are read from `~/.netrc`, stored there by `conjur authn login`.

In general, you can ignore the `CONJUR_CORE_URL` and `CONJUR_AUTHN_URL` unless
you need to specify, for example, an authn proxy.

The provider will fail unless all of the following values are provided:

- `CONJUR_MAJOR_VERSION=4` for Conjur v4.9
- An appliance url (`CONJUR_APPLIANCE_URL`)
- An organization account (`CONJUR_ACCOUNT`)
- A username and api key, or Conjur authn token, or a path to `CONJUR_AUTHN_TOKEN_FILE` a dynamic Conjur authn token
- A path to (`CONJUR_CERT_FILE`) **or** content of (`CONJUR_SSL_CERTIFICATE`) the appliance's public SSL certificate

---

## Contributing

We welcome contributions of all kinds to this repository. For instructions on how to get started and descriptions of our development workflows, please see our [contributing
guide][contrib].

[contrib]: CONTRIBUTING.md
