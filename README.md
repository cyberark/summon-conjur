# summon-conjur

Conjur provider for [Summon](https://conjurinc.github.io/summon/).

**Note** Use the [summon-conjurcli](https://github.com/conjurinc/summon-conjurcli) provider if you are on Conjur v4.4.0 or earlier.

## Install

Download the [latest release](https://github.com/conjurinc/summon-conjur/releases) and extract it to the directory `/usr/local/lib/summon`.

## Usage

Give summon-conjur a variable name and it will fetch it for you and print
the value to stdout.

```sh-session
$ summon-conjur prod/aws/iam/user/robot/access_key_id
8h9psadf89sdahfp98
```

## Configuration

This provider uses the same configuration pattern as the [Conjur CLI
Client](https://github.com/conjurinc/api-ruby#configuration) to connect to Conjur.
Specifically, it loads configuration from:

 * `.conjurrc` files, located in the home and current directories, or at the
    path specified by the `CONJURRC` environment variable.
 * Read `/etc/conjur.conf` as a `.conjurrc` file.
 * Read `/etc/conjur.identity` as a `netrc` file. Note that the user running must either be in the group `conjur` or root to read the identity file.
 * Environment variables:
   * Appliance URLs
     * `CONJUR_APPLIANCE_URL`
     * `CONJUR_CORE_URL`
     * `CONJUR_AUTHN_URL`
   * SSL certificate
     * `CONJUR_CERT_FILE`
     * `CONJUR_SSL_CERTIFICATE`
   * Authentication
     * `CONJUR_AUTHN_LOGIN`
     * `CONJUR_AUTHN_API_KEY`
     * `CONJUR_AUTHN_TOKEN`

If `CONJUR_AUTHN_LOGIN` and `CONJUR_AUTHN_API_KEY` or `CONJUR_AUTHN_TOKEN` are not provided, the username and API key are read from `~/.netrc`, stored there by `conjur authn login`.

In general, you can ignore the `CONJUR_CORE_URL` and `CONJUR_AUTHN_URL` unless
you need to specify, for example, an authn proxy.

The provider will fail unless all of the following values are provided:

 * The appliance url
 * A username and api key, or Conjur authn token
 * A path to (`CONJUR_CERT_FILE`) **or** content of (`CONJUR_SSL_CERTIFICATE`) the appliance's
 public SSL certificate
