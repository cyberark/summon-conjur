config:
  env
  conjur_rc

env:
  CONJUR_MAJOR_VERSION - must be set to 4 in order for summon-conjur to work with Conjur v4.9.
  config
    CONJUR_APPLIANCE_URL
    CONJUR_AUTHN_URL - not supported in summon-conjur but should become supported!
    ssl_cert
      CONJUR_CERT_FILE
      CONJUR_SSL_CERTIFICATE
  authn
    CONJUR_AUTHN_TOKEN - authenticator exists but not being used!
    CONJUR_AUTHN_TOKEN_FILE
    CONJUR_AUTHN_LOGIN, CONJUR_AUTHN_API_KEY

conjur_rc:
  $CONJURRC
  or
  $PWD/.conjurrc
  $HOME/.conjurrc
```
---
appliance_url: http://path/to/appliance
account: some account
cert_file: "/path/to/cert/file/pem"
```

ssl_cert:
  $CONJUR_SSL_CERTIFICATE
  $CONJUR_CERT_FILE

authn:
  creds
  net_rc

creds:
  $CONJUR_AUTHN_TOKEN_FILE
  $CONJUR_AUTHN_LOGIN, $CONJUR_AUTH_API_KEY

net_rc:
  $CONJUR_NETRC_PATH
  $HOME.netrc
  /etc/conjur.conf
```
machine [http|https]://path/to/conjur/authn/endpoint
  login admin
  password 242gyw4307d2sn16zr36d2z7s03y30mcmkr212g02d1gyn5tpjh4fve
```
+ net_rc only applies if the machine name is identical to `{the config's appliance_url}/auth`


TODO:
+ create an ASCII graphic of where config was taken from
+ verbose mode in api-go-go
+ several logging levels
