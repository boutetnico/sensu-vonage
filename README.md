[![Bonsai Asset Badge](https://img.shields.io/badge/Sensu%20Slack%20Handler-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/boutetnico/sensu-vonage-handler)

# Sensu Slack Handler

- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Handler definition](#handler-definition)
  - [Check definition](#check-definition)
- [Installation from source and
  contributing](#installation-from-source-and-contributing)

## Overview


The [Sensu Vonage Handler][0] is a [Sensu Event Handler][3] that sends event data via SMS.

## Usage examples

Help:

```
Usage:
  sensu-vonage-handler [flags]
  sensu-vonage-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -k, --api-key string                The Vonage API key
  -s, --api-secret string             The Vonage API secret
  -h, --help                          help for sensu-vonage-handler
  -f, --from string                   The name/number of the sender
  -r, --recipients string             Comma-separated list of recipients
```

## Configuration

### Asset registration

Assets are the best way to make use of this handler. If you're not using an asset, please consider doing so! If you're using sensuctl 5.13 or later, you can use the following command to add the asset:

`sensuctl asset add sensu/sensu-vonage-handler`

If you're using an earlier version of sensuctl, you can download the asset
definition from [this project's Bonsai Asset Index
page][6].

### Handler definition

Create the handler using the following handler definition:

```yml
---
api_version: core/v2
type: Handler
metadata:
  namespace: default
  name: vonage
spec:
  type: pipe
  command: sensu-vonage-handler --from Sensu --recipients 4499999999,4488888888
  filters:
  - is_incident
  runtime_assets:
  - sensu/sensu-vonage-handler
  secrets:
  - name: VONAGE_API_KEY
    secret: vonage-api-key
  - name: VONAGE_API_SECRET
    secret: vonage-api-secret
  timeout: 10
```

### Check definition

```
api_version: core/v2
type: CheckConfig
metadata:
  namespace: default
  name: dummy-app-healthz
spec:
  command: check-http -u http://localhost:8080/healthz
  subscriptions:
  - dummy
  publish: true
  interval: 10
  handlers:
  - vonage
```

### Customizing configuration options via checks and entities

All configuration options of this handler can be overridden via the annotations
of checks and entities. For example, to customize the recipients for a given
entity, you could use the following sensu-agent configuration snippet:

```yml
# /etc/sensu/agent.yml example
annotations:
  sensu.io/plugins/vonage/config/recipients: '4477777777'
```

### Proxy Support

This handler supports the use of the environment variables HTTP_PROXY,
HTTPS_PROXY, and NO_PROXY (or the lowercase versions thereof). HTTPS_PROXY takes
precedence over HTTP_PROXY for https requests.  The environment values may be
either a complete URL or a "host[:port]", in which case the "http" scheme is assumed.

## Installing from source and contributing

Download the latest version of the sensu-vonage-handler from [releases][4],
or create an executable script from this source.

### Compiling

From the local path of the sensu-vonage-handler repository:
```
go build
```

To contribute to this plugin, see [CONTRIBUTING](https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md)

[0]: https://github.com/boutetnico/sensu-vonage-handler
[1]: https://github.com/sensu/sensu-go
[3]: https://docs.sensu.io/sensu-go/latest/reference/handlers/#how-do-sensu-handlers-work
[4]: https://github.com/boutetnico/sensu-vonage-handler/releases
[5]: https://docs.sensu.io/sensu-go/latest/reference/secrets/
[6]: https://bonsai.sensu.io/assets/boutetnico/sensu-vonage-handler
