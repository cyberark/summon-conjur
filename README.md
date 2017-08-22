# summon-conjur

Conjur provider for [Summon](https://cyberark.github.io/summon/).

**Note** Use the [summon-conjurcli](https://github.com/conjurinc/summon-conjurcli) provider if you are on Conjur v4.4.0 or earlier.

## Install

**Note** Check the release notes and select an appropriate release to ensure support for your version of Conjur.

Download the [latest release](https://github.com/cyberark/summon-conjur/releases) and extract it to the directory `/usr/local/lib/summon`.

## Usage in isolation

Give summon-conjur a variable name and it will fetch it for you and print the value to stdout.

```sh-session
$ summon-conjur prod/aws/iam/user/robot/access_key_id
8h9psadf89sdahfp98
```

## Usage as a provider for Summon

[Summon](https://cyberark.github.io/summon/) is a command-line tool that reads a file in secrets.yml format and injects secrets as environment variables into any process. Once the process exits, the secrets are gone.

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
$ summon --provider summon-conjur env
...
AWS_ACCESS_KEY_ID=xxxxxxxxxxxxxxxxxxx
AWS_SECRET_ACCESS_KEY=yyyyyyyyyyyyyyyy
...
```

`summon` resolves the entries in secrets.yml with the conjur provider and makes the secret values available to the environment of the command `env`.

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
 
### Extract environment variables from machine identity files
 
Previous versions of summon-conjur support loading configuration from [machine identity files](https://developer.conjur.net/key_concepts/machine_identity.html#storage-files). While that support is not available in this release, it is relatively straight-forward to extract the relevant environment variables from the machine identity files.
 
```bash
CONJUR_IDENTITY_PATH=/etc/conjur.conf
export CONJUR_ACCOUNT=$(awk '/account:/ {print $2}' $CONJUR_IDENTITY_PATH)
export CONJUR_APPLIANCE_URL=$(awk '/appliance_url:/ {print $2}' $CONJUR_IDENTITY_PATH)
CONJUR_NETRC_PATH=$(awk '/netrc_path:/ {print $2}' $CONJUR_IDENTITY_PATH)
: ${CONJUR_NETRC_PATH:="/etc/conjur.identity"}
```
#### Extract Login Strategy Credentials
```bash
export CONJUR_AUTHN_LOGIN=$(awk -v CONJUR_APPLIANCE_URL=$CONJUR_APPLIANCE_URL '$0 ~ CONJUR_APPLIANCE_URL {f=1} f && /login/ {print $2;f=0}' $CONJUR_NETRC_PATH)
export CONJUR_AUTHN_API_KEY=$(awk -v CONJUR_APPLIANCE_URL=$CONJUR_APPLIANCE_URL '$0 ~ CONJUR_APPLIANCE_URL {f=1} f && /password/ {print $2;f=0}' $CONJUR_NETRC_PATH)
```

#### Extract Token File Strategy Credentials
```bash
function finish {
    echo "Cleaning up..."
    kill $CONJUR_AUTHN_TOKEN_FILE_PID
    echo "All done."
}
trap finish EXIT
while true; do echo $(conjur authn authenticate) > /tmp/CONJUR_AUTHN_TOKEN_FILE ; sleep 2; done &
export CONJUR_AUTHN_TOKEN_FILE_PID=$!
export CONJUR_AUTHN_TOKEN_FILE=/tmp/CONJUR_AUTHN_TOKEN_FILE
```
