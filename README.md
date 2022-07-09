# azure-auth - tools for interacting with Azure and Microsoft Office365 system.
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/CyCoreSystems/azure-auth)

The primary impetus here was to support email-oriented operations.  Azure uses a
version of XOAUTH2-based authentication, but its currently-distributed Go
libraries do not support things like refresh tokens.

Two CLI tools exist here:
 - [azure-token](cmd/azure-token):  retrieves and caches an access token for use with subsequent
   operations, such as IMAP and SMTP calls.
 - [send-azure](cmd/send-azure):  a sendmail drop-in which uses XOAUTH2 for authentication.

Additionally, some useful library packages are available to support these and
other operations:
 - [token](pkg/token): deals with token fetching, required human interation, and caching
 - [smtp](pkg/smtp): XOAUTH2 generator and SMTP client
 - [config](pkg/config): configuration file reader for key and endpoint data

## Configuration file

Most of these tools expect a configuration file to exist at
`${XDG_CONFIG_PATH}/azure/config.yaml`.  It generally looks like:

```yaml
username: me@mydomain.com
tenantID: d6e2c434-8a03-4eb7-87b5-d96dc44a651c
clientID: 44ad7be1-4a36-42bf-8fdd-6ba2e2b53890
clientSecret:  blah-BLAH-bl@h.blah~blah
scopes:
  - https://outlook.office.com/SMTP.Send
  - https://outlook.office.com/IMAP.AccessAsUser.All
  - offline_access
redirect:
  host: localhost
  port: 5050
  path: "getToken"
```
