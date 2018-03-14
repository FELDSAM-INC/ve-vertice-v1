
vertice
=======

[![Travis Build Status](https://travis-ci.org/VirtEngine/vertice.svg?branch=1.5.2)](https://travis-ci.org/VirtEngine/vertice)

Vertice is the core engine for VirtEngine Vertice 1.5.x and is open source.


# Roadmap for 2.0

Read the [Deployment design](https://github.com/virtengine/verticedev/blob/master/proposals/01.deployments.md).

## Where is the code for 2.0

We have moved the development to private gitlab as it will have enterprise features.

## When can i get it in my anxious hands

`2.0` will be released on `Sep 30 2017` or less.

`2.0` developed with private enterprise features and is moved to gitlab.


### Requirements

>
[Golang 1.8 > +](http://www.golang.org/dl)

>
[NSQ 0.3.x](http://nsq.io/deployment/installing.html)

>
[Cassandra 3 +](https://wiki.apache.org/cassandra/GettingStarted)

## Usage

``vertice -v start``

### Compile from source

```bash
$ mkdir -p $GOPATH/src/github.com/virtengine
$ cd $GOPATH/src/github.com/virtengine
$ git clone https://github.com/virtengine/vertice.git
$ cd vertice
$ make build
```

You can use the following command to run it:

```bash
$ ./vertice -v start --config ./conf/vertice.conf
```

### Documentation

[development documentation] (https://github.com/virtengine/verticedev/tree/master/development)

[documentation] (http://docs.virtengine.com) for usage.




We are glad to help if you have questions, or request for new features..

[twitter @virtengine](http://twitter.com/virtengine) [email support@virtengine.com](<support@virtengine.com>)

[devkit] (https://github.com/virtengine/verticedev)

# License


|                      |                                          |
|:---------------------|:-----------------------------------------|
| **Author:**          | Rajthilak (<rajthilak@megam.io>)
| 	                   | KishorekumarNeelamegam (<nkishore@megam.io>)
|                      | Ranjitha  (<ranjithar@megam.io>)
|                      | MVijaykanth  (<mvijaykanth@megam.io>)
| **Copyright:**       | Copyright (c) 2013-2017 Megam Systems.
| **License:**         | Apache License, Version 2.0

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
