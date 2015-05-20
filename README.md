# cauldron-conjur

Cauldron implementation with Conjur backend

## Usage

Give cauldron-conjur a variable name and it will fetch it for you and print
the value to stdout.

```sh-session
$ cauldron-conjur prod/aws/iam/user/robot/access_key_id
8h9psadf89sdahfp98
```

Reading Conjur configuration isn't implemented yet, so set these environment
variables for development.

```
GO_CONJUR_APPLIANCE_URL
GO_CONJUR_AUTHN_LOGIN
GO_CONJUR_AUTHN_API_KEY
GO_CONJUR_SSL_CERTIFICATE_PATH
```

This implementation has no 3rd-party dependencies.
