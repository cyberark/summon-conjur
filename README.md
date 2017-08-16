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

This provider uses a similar configuration pattern as the [Conjur CLI
Client](https://github.com/conjurinc/api-ruby#configuration) to connect to Conjur.
Specifically, it loads configuration from:

 * Environment variables:
   * Config
     * `CONJUR_APPLIANCE_URL`
     * `CONJUR_ACCOUNT`
   * Authentication
     * `CONJUR_AUTHN_LOGIN`
     * `CONJUR_AUTHN_API_KEY`
     * `CONJUR_AUTHN_TOKEN_FILE`

The provider will fail unless all of the following values are provided:

 * An appliance url
 * An organisation account
 * A username and api key, or Conjur authn token file
