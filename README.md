# ds-load
CLI pipeline for populating the directory

## Arguments

```
Usage: ds-load <command>

Directory loader

Commands:
  exec                  import data in directory
  get-plugin            download plugin
  set-default-plugin    sets a plugin as default
  list-plugins          list available plugins
  version               version information

Flags:
  -h, --help                  Show context-sensitive help.
  -c, --config=CONFIG-FLAG    Path to the config file. Any argument provided to the CLI will take precedence.
  -v, --verbosity=INT         Use to increase output verbosity.
```
The default command for ds-load is `exec`, so when running `ds-load` without any command, the exec parameters need to be passed.
The CLI parameters need to be passed first, and then, an arbitrary list of positional parameters can be passed. The first positional parameter is the plugin name we want to invoke followed by the plugin parameters. `ds-load --host=<directory host> auth0 --domain=<auth0 domain>`.

For viewing the plugin help: `ds-load auth0 --help`. 

### Parameter position
When running `ds-load auth0 --key`, `auth0` is a positional parameter, so `--key` will be run in the context of the plugin. If we run `ds-load --key auth0`, `--key` would be a parameter to `ds-load exec`.

### ds-load exec
`exec` is the default command. it will invoke a plugin with the specified parameters reading it's output and importing the resulting data into the directory.

```
Usage: ds-load exec <command> ...

import data in directory

Arguments:
  <command> ...    available commands are: auth0|okta

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

### Env variables
Parameters can also be passed by environment, as seen in the help message of each command, but the ones from config files and command line take precedence.

## Config files

config files are in yaml format:
```yaml
---
command:
  arg: value
  another-arg: value
another-command:
  arg: value
```

When passing custom config files to both the cli and the plugin, use `ds-load -c <config-path> plugin -c <plugin-config-path>` 

### CLI config
default location: `~/.config/ds-load/cfg/config.yaml` can be overridden using `-c/--config`

#### example
```yaml
---
exec:
  host: directory.eng.aserto.com:8443
  api-key: secretapikey
  tenant-id: your-tenant-id
```

### Plugin config
default location: `~/.config/ds-load/cfg/<plugin-name>.yaml` can be overridden using `-c/--config`

#### example
```yaml
---
exec:
  domain: "domain.auth0.com"
  client-id: "clientid"
  client-secret: "clientsupersecret"
  template-file: "/path/to/transform.file"
  max-chunk-size: 20
fetch:
  domain: "domain.auth0.com"
  client-id: "clientid"
  client-secret: "clientsupersecret"
```

## Transform
The data received from the fetcher is being transformed using a transformation template, which is written as a go template and it outputs objects and relations.

The default transformation template can be exported using `ds-load <plugin-name> export-transform`.

A custom transformation file can be provided when running the plugin in `exec` or `transform` mode via the `--template-file` parameter.

## Logs

Logs are printed to `stdout`. You can increase detail using the verbosity flag (e.g. `-vvv`).

## Usage examples

### Import from auth0 into the directory
```
ds-load --host=<directory-host> --api-key=<directory-api-key> --tenant-id=<tenant-id> auth0 --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret> --max-chunk-size=20
```

### Import data with custom transformation file
```
ds-load --host=<directory-host> --api-key=<directory-api-key> --tenant-id=<tenant-id> auth0 --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret> --template-file=<template-path> --max-chunk-size=20
```

### View data from auth0
```
ds-load auth0 fetch --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret> --max-chunk-size=20
```

### Transform data from a previously saved auth0 fetch
we use `-p` in order to just print the transform data
```
ds-load auth0 fetch --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret> > auth0.data
cat auth0.data | ds-load -p auth0 transform
```

### Transform and import data from a previously saved auth0 fetch
```
ds-load auth0 fetch --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret> > auth0.data
cat auth0.data | ds-load --host=<directory-host> --api-key=<directory-api-key> --tenant-id=<tenant-id> auth0 transform
```

### Pipe data from fetch to transform
```
ds-load auth0 fetch --domain=<auth0-domain> --client-id=<auth0-client-id> --client-secret=<auth0-client-secret> | ds-load -p auth0 transform
```
