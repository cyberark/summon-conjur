# summon-conjur

CyberArk Secrets Manager provider for [Summon](https://github.com/cyberark/summon).

[![GitHub release](https://img.shields.io/github/release/cyberark/summon-conjur.svg)](https://github.com/cyberark/summon-conjur/releases/latest)

[![Github commits (since latest release)](https://img.shields.io/github/commits-since/cyberark/summon-conjur/latest.svg)](https://github.com/cyberark/summon-conjur/commits/main)

---

## Install

Pre-built binaries and packages are available from GitHub releases
[here](https://github.com/cyberark/summon-conjur/releases).

### Using summon-conjur with Conjur Open Source

Are you using this project with [Conjur Open Source](https://github.com/cyberark/conjur)? Then we
**strongly** recommend choosing the version of this project to use from the latest [Conjur OSS
suite release](https://docs.conjur.org/Latest/en/Content/Overview/Conjur-OSS-Suite-Overview.html).
Conjur maintainers perform additional testing on the suite release versions to ensure
compatibility. When possible, upgrade your Conjur version to match the
[latest suite release](https://docs.conjur.org/Latest/en/Content/ReleaseNotes/ConjurOSS-suite-RN.htm);
when using integrations, choose the latest suite release that matches your Conjur version. For any
questions, please contact us on [Discourse](https://discuss.cyberarkcommons.org/c/conjur/5).

### Homebrew

```bash
brew tap cyberark/tools
brew install summon-conjur
```

### Linux (Debian and Red Hat flavors)

`deb` and `rpm` files are attached to new releases.
These can be installed with `dpkg -i summon-conjur_*.deb` and
`rpm -ivh summon-conjur_*.rpm`, respectively.

### Auto Install

**Note** Check the release notes and select an appropriate release to ensure support for your version of CyberArk Secrets Manager.

Use the auto-install script. This will install the latest version of summon-conjur.
The script requires sudo to place summon-conjur in dir `/usr/local/lib/summon`.

```bash
curl -sSL https://raw.githubusercontent.com/cyberark/summon-conjur/main/install.sh | bash
```

### Manual Install

Otherwise, download the [latest release](https://github.com/cyberark/summon-conjur/releases) and extract it to the directory `/usr/local/lib/summon`.

## Usage in isolation

Give summon-conjur a variable name and it will fetch it for you and print the value to stdout.

```shell
$ summon-conjur prod/aws/iam/user/robot/access_key_id
flgwkeatfghhdqkflaqiwoagsmfgxool
```

You can also use interactive mode by starting the command without any arguments 
and then passing paths to secrets one by one. This way you can fetch multiple values in a single command run.
Keep in mind that by using interactive mode outputted values will be in BASE64 format.

```shell
$ summon-conjur
prod/aws/iam/user/robot/access_key_id
Zmxnd2tlYXRmZ2hoZHFrZmxhcWl3b2Fnc21mZ3hvb2w=
prod/aws/s3/bucket_name/access_key_id
YWNudmdlb3dycmd4dW1ic2tncW51Zm50dmRvYWVic3A=
```

### Flags

```txt
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

### Example

As an example let's use the `env` command:

Following installation, define your keys in a `secrets.yml` file

```yml
AWS_ACCESS_KEY_ID: !var aws/iam/user/robot/access_key_id
AWS_SECRET_ACCESS_KEY: !var aws/iam/user/robot/secret_access_key
```

By default, summon will look for `secrets.yml` in the directory it is called from and export the secret values to the environment of the command it wraps.

Wrap the `env` in summon:

```sh
$ summon --provider summon-conjur env
...
AWS_ACCESS_KEY_ID=AKIAJS34242K1123J3K43
AWS_SECRET_ACCESS_KEY=A23MSKSKSJASHDIWM
...
```

`summon` resolves the entries in secrets.yml with the CyberArk Secrets Manager provider and makes the secret values available to the environment of the command `env`.

## Configuration

This provider uses the same configuration pattern as the [CyberArk Secrets Manager CLI](https://github.com/cyberark/conjur-cli-go)
to connect to Conjur. Specifically, it loads configuration from:

* `.conjurrc` files, located in the home and current directories, or at the
    path specified by the `CONJURRC` environment variable.
* Reads the `.conjurrc` file from `/etc/conjur.conf` on Linux/macOS and `C:\Windows\conjur.conf` on Windows.
* Environment variables:
  * Appliance URLs
    * `CONJUR_APPLIANCE_URL`
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
    * JWT Token
      * `CONJUR_AUTHN_JWT_SERVICE_ID`  (e.g. `kubernetes`)
      * `JWT_TOKEN_PATH` (optional)  (default: `/var/run/secrets/kubernetes.io/serviceaccount/token`)
    * AWS/Azure/GCP
      * `CONJUR_AUTHN_TYPE` (set to `iam`, `azure`, or `gcp`)
      * `CONJUR_SERVICE_ID` (except for GCP)
      * `CONJUR_AUTHN_JWT_HOST_ID`
      * `CONJUR_AUTHN_JWT_TOKEN` (optional - if not set, token will be read from the metadata service)


If `CONJUR_AUTHN_LOGIN` and `CONJUR_AUTHN_API_KEY` or `CONJUR_AUTHN_TOKEN` or `CONJUR_AUTHN_TOKEN_FILE` or `CONJUR_AUTHN_JWT_SERVICE_ID` are not provided, the username and API key are read from system keychain or `~/.netrc`, stored there by `conjur login`.

On systems that support keychain storage, that will be used by default, and if that fails the `~/.netrc` file will be used,
though this behavior can be modified in the `.conjurrc` file:

```yaml
...
credential_storage: "netrc"
netrc_path: "/etc/conjur.identity"
...
```

The provider will fail unless all of the following values are provided:

* An appliance url (`CONJUR_APPLIANCE_URL`)
* An organization account (`CONJUR_ACCOUNT`)
* A valid authentication method (e.g., username/api key, token, or JWT or cloud auth configuration)
* A path to (`CONJUR_CERT_FILE`) **or** content of (`CONJUR_SSL_CERTIFICATE`) the appliance's public SSL certificate

---

## Contributing

We welcome contributions of all kinds to this repository. For instructions on how to get started and descriptions of our development workflows, please see our [contributing
guide][contrib].

[contrib]: CONTRIBUTING.md
