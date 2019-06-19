<a id="markdown-component-connstate---settings-component-for-collecting-http-server-metrics" name="component-connstate---settings-component-for-collecting-http-server-metrics"></a>
# component-connstate - Settings component for collecting HTTP server metrics

<!-- TOC -->

- [component-connstate - Settings component for collecting HTTP server metrics](#component-connstate---settings-component-for-collecting-http-server-metrics)
    - [Overview](#overview)
    - [Quick Start](#quick-start)
    - [Status](#status)
    - [Contributing](#contributing)
        - [Building And Testing](#building-and-testing)
        - [License](#license)
        - [Contributing Agreement](#contributing-agreement)

<!-- /TOC -->

<a id="markdown-overview" name="overview"></a>
## Overview

This is a [`settings`](https://github.com/asecurityteam/settings) that enables
constructing a metrics gathering component that integrates with the
[`http.Server.ConnState`](https://golang.org/pkg/net/http/#Server) hook. It
works by tracking the state of all active connections and any state changes of
a connection while the server runs.

<a id="markdown-quick-start" name="quick-start"></a>
## Quick Start

```golang
package main

import (
    "context"
    "net/http"

    stat "github.com/asecurityteam/component-stat"
    stat "github.com/asecurityteam/component-connstate"
    "github.com/asecurityteam/settings"
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

<a id="markdown-status" name="status"></a>
## Status

This project is in incubation which means we are not yet operating this tool in
production and the interfaces are subject to change.

<a id="markdown-contributing" name="contributing"></a>
## Contributing

<a id="markdown-building-and-testing" name="building-and-testing"></a>
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

<a id="markdown-license" name="license"></a>
### License

This project is licensed under Apache 2.0. See LICENSE.txt for details.

<a id="markdown-contributing-agreement" name="contributing-agreement"></a>
### Contributing Agreement

Atlassian requires signing a contributor's agreement before we can accept a patch. If
you are an individual you can fill out the [individual
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=3f94fbdc-2fbe-46ac-b14c-5d152700ae5d).
If you are contributing on behalf of your company then please fill out the [corporate
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=e1c17c66-ca4d-4aab-a953-2c231af4a20b).
