# ds-load

`ds-load` is a CLI pipeline for populating the Topaz directory (or directories that are contract-compatible, such as the Aserto directory).

## Arguments

```
Usage: ds-load <command>

Directory loader

Commands:
  exec                  import data in directory by running fetch, transform and publish
  publish               load data from stdin into directory
  get-plugin            download plugin
  set-default-plugin    sets a plugin as default
  list-plugins          list available plugins
  version               version information

Flags:
  -h, --help                  Show context-sensitive help.
  -c, --config=CONFIG-FLAG    Path to the config file. Any argument provided to the CLI will take precedence.
  -v, --verbosity=INT         Use to increase output verbosity.
```

The `ds-load` pipeline has three stages: `fetch`, `transform`, and `publish`:
* `fetch` retrieves plugin-specific data from the source in a native JSON format
* `transform` converts this data into directory objects and relations
* `publish` loads the objects and relations data into the directory

The default command for ds-load is `exec`, which executes all three stages. When running `ds-load` without any command, the exec parameters need to be passed.

Note: the plugin examples below all use `auth0`, but every one of the plugins (`azuread`, `cognito`, `google`, `okta`, etc) follow the same patterns.

### ds-load vs plugin parameters
The ds-load CLI parameters need to be passed first, and can be followed by an arbitrary list of positional parameters. The first positional parameter is the plugin name we want to invoke followed by the plugin's parameters.

Example: `ds-load --host=<directory host> auth0 --domain=<auth0 domain>`
* `--host` is a CLI parameter for `ds-load`
* `auth0` is the plugin name
* `--domain` is a parameter for the `auth0` plugin

For viewing the plugin help, use the following format: `ds-load auth0 --help`.

Tip: when running `ds-load auth0 --key`, `auth0` is a positional parameter, so `--key` will be run in the context of the plugin. If we run `ds-load --key auth0`, `--key` would be a parameter to `ds-load exec`.

### ds-load exec
`exec` is the default command. it will invoke a plugin with the specified parameters reading its output and importing the resulting data into the directory.

```
Usage: ds-load exec <command> ...

import data in directory by running fetch, transform and publish

Arguments:
  <command> ...    available commands are: auth0|azuread|cognito|google|okta

Flags:
  -h, --help                  Show context-sensitive help.
  -c, --config=CONFIG-FLAG    Path to the config file. Any argument provided to the CLI will take precedence.
  -v, --verbosity=INT         Use to increase output verbosity.

  -s, --host=STRING           Directory host address ($DIRECTORY_HOST)
  -k, --api-key=STRING        Directory API Key ($DIRECTORY_API_KEY)
  -i, --insecure              Disable TLS verification
  -t, --tenant-id=STRING      Directory Tenant ID ($DIRECTORY_TENANT_ID)
  -p, --print                 print output to stdout
```

`-p/--print` is enabled by default when invoking a plugin with `fetch/version/export-transform/--help`

### Environment variables
Parameters can also be passed using environment variables, as seen in the help message of each command, but the ones from config files and command line take precedence.

## Config files

Config files are in yaml format:
```yaml
---
arg: value
another-arg: value
<plugin-name>:
  arg: value
```

When passing custom config files to both the cli and the plugin, use `ds-load -c <config-path> <plugin-name> <command>` 

### CLI config
The default location for the configuration file is `~/.config/ds-load/cfg/config.yaml`. It can be overridden using the `-c/--config` flag.

#### example with auth0 plugin
```yaml
---
host: directory.prod.aserto.com:8443
api-key: "secretapikey"
tenant-id: your-tenant-id
auth0:
  domain: "domain.auth0.com"
  client-id: "clientid"
  client-secret: "clientsupersecret"
  template: "/path/to/transform.file"
```

### Plugin config
The default location for plugin configuration files is `~/.config/ds-load/cfg/<plugin-name>.yaml`. It can be overridden using the `-c/--config` flag.

#### example for auth0
```yaml
---
auth0:
  domain: "domain.auth0.com"
  client-id: "clientid"
  client-secret: "clientsupersecret"
  template: "/path/to/transform.file"
```

## Transform
The data received from the fetcher is transformed into objects and relations using a transformation template, which is uses the go template syntax.

The default transformation template can be exported using `ds-load <plugin-name> export-transform`.

A custom transformation file can be provided when running the plugin in `exec` or `transform` mode via the `--template` parameter.

More information on the transformation template language can be found in the [tranform template docs](./docs/templates.md).

## Logs

Logs are printed to `stdout`. You can increase detail using the verbosity flag (e.g. `-vvv`).

## Fetching source data

Most plugins fetch entire objects that are available in the transform template. `azuread` and `azureadb2c` use the msgraph api to query the data and by default only query only properties that are needed in the default transform template.

### AzureAD and AzureADB2C

To use a custom porperty list in the query, you can use the CLI parameters or configure them in your config:

```
azuread:
  tenant: "tenant-id"
  client-id: "client-id"
  client-secret: "secret"
  groups: true
  user-properties: ["id", "displayName"]
  group-properties: ["id", "displayName"]
```

AzureAD Groups:

- displayName
- id
- mail
- createdDateTime
- mailNickname

AzureADB2C Users:

- displayName
- id
- mail
- createdDateTime
- mobilePhone
- userPrincipalName
- accountEnabled
- identities
- creationType

AzureADB2C Groups:

- displayName
- id
- mail
- createdDateTime
- mailNickname
- members
- transitiveMembers

A list of all available properties is available on the Microsoft website for [user object type](https://learn.microsoft.com/en-us/graph/api/resources/user?view=graph-rest-1.0#properties) and [group object type](https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-1.0#properties)

## Usage examples

### Import from auth0 into the directory
```
ds-load --host=<directory-host> --api-key=<directory-api-key> --tenant-id=<tenant-id> auth0 --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret>
```

### Import data with a custom transformation file
```
ds-load --host=<directory-host> --api-key=<directory-api-key> --tenant-id=<tenant-id> auth0 --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret> --template=<template-path>
```

### Fetch data from auth0 without importing it
```
ds-load auth0 fetch --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret>
```

### Transform data from a previously saved auth0 fetch
Note: we use `-p` in order to just print the transform data.
```
ds-load auth0 fetch --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret> > auth0.json
cat auth0.json | ds-load -p auth0 transform
```

### Transform and import data from a previously saved auth0 fetch
```
ds-load auth0 fetch --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret> > auth0.json

cat auth0.json | ds-load --host=<directory-host> --api-key=<directory-api-key> --tenant-id=<tenant-id> auth0 transform
```

### Pipe data from fetch to transform
```
ds-load auth0 fetch --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret> | ds-load -p auth0 transform
```

### Use config file to import data from auth0 into the directory

config.yaml
```yaml
---
host: "directory.prod.aserto.com:8443"
api-key: "secretapikey"
tenant-id: "your-tenant-id"
auth0:
  domain: "domain.auth0.com"
  client-id: "clientid"
  client-secret: "clientsupersecret"
```

```
ds-load -c ./config.yaml auth0
```

### Load directory data from a file
```
ds-load -p auth0 --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret> > auth0.json

cat auth0.json | ds-load publish --host=<directory-host> --api-key=<directory-api-key> --tenant-id=<tenant-id>
```

