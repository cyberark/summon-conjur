# cauldron-conjur

Conjur provider for [Cauldron](https://conjurinc.github.io/cauldron/).

**Note** Use the [cauldron-conjurcli](https://github.com/conjurinc/cauldron-conjurcli) provider if you are on Conjur v4.4.0 or earlier.

## Usage

Give cauldron-conjur a variable name and it will fetch it for you and print
the value to stdout.

```sh-session
$ cauldron-conjur prod/aws/iam/user/robot/access_key_id
8h9psadf89sdahfp98
```

## Configuration

This provider uses the same configuration pattern as the [Conjur CLI
Client](https://github.com/conjurinc/api-ruby#configuration) to connect to Conjur.  
Specifically, it loads configuration from:

 * `.conjurrc` files, located in the home and current directories, or at the 
    path specified by the `CONJURRC` environment variable.
 * Read `/etc/conjur.conf` as a `.conjurrc` file.
 * Environment variables:
    * `CONJUR_AUTHN_LOGIN`
    * `CONJUR_API_KEY`
    * `CONJUR_CERT_FILE`
    * `CONJUR_APPLIANCE_URL`
    * `CONJUR_CORE_URL`
    * `CONJUR_AUTHN_URL`
 * A username and api key can be read from `~/.netrc` if stored there by
    `conjur authn login`

In general, you can ignore the `CONJUR_CORE_URL` and `CONJUR_AUTHN_URL` unless
you need to specify, for example, and authn proxy.

The provider will fail unless all of the following values are provided:

 * The appliance url
 * A username and api key
 * A path to the appliance's SSL certificate


