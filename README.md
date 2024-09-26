# component-connstate - Settings component for collecting HTTP server metrics
[![GoDoc](https://godoc.org/github.com/asecurityteam/component-connstate?status.svg)](https://godoc.org/github.com/asecurityteam/component-connstate)
<!-- TOC -->autoauto- [component-connstate - Settings component for collecting HTTP server metrics](#component-connstate---settings-component-for-collecting-http-server-metrics)auto    - [Overview](#overview)auto    - [Quick Start](#quick-start)auto    - [Status](#status)auto    - [Contributing](#contributing)auto        - [Building And Testing](#building-and-testing)auto        - [License](#license)auto        - [Contributing Agreement](#contributing-agreement)autoauto<!-- /TOC -->

## Overview

This is a [`settings`](https://github.com/asecurityteam/settings) that enables
constructing a metrics gathering component that integrates with the
[`http.Server.ConnState`](https://golang.org/pkg/net/http/#Server) hook. It
works by tracking the state of all active connections and any state changes of
a connection while the server runs.

## Quick Start

```golang
package main

import (
    "context"
    "net/http"

    stat "github.com/asecurityteam/component-stat"
    stat "github.com/asecurityteam/component-connstate"
    "github.com/asecurityteam/settings/v2"
)

func main() {
    ctx := context.Background()
    envSource := settings.NewEnvSource(os.Environ())

    s, _ := stat.New(ctx, envSource)
    cs, _ := connstate.Load(ctx, envSource, connstate.NewComponent().WithStat(s))
    defer cs.Close()

    srv := &http.Server{
        ConnState: cs.HandleEvent,
    }
    go cs.Report()

    _ = srv.ListenAndServe()
}
```

## Status

This project is in incubation which means we are not yet operating this tool in
production and the interfaces are subject to change.

## Contributing

### Building And Testing

We publish a docker image called [SDCLI](https://github.com/asecurityteam/sdcli) that
bundles all of our build dependencies. It is used by the included Makefile to help
make building and testing a bit easier. The following actions are available through
the Makefile:

-   make dep

    Install the project dependencies into a vendor directory

-   make lint

    Run our static analysis suite

-   make test

    Run unit tests and generate a coverage artifact

-   make integration

    Run integration tests and generate a coverage artifact

-   make coverage

    Report the combined coverage for unit and integration tests

### License

This project is licensed under Apache 2.0. See LICENSE.txt for details.

### Contributing Agreement

Atlassian requires signing a contributor's agreement before we can accept a patch. If
you are an individual you can fill out the [individual
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=3f94fbdc-2fbe-46ac-b14c-5d152700ae5d).
If you are contributing on behalf of your company then please fill out the [corporate
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=e1c17c66-ca4d-4aab-a953-2c231af4a20b).
