# Todo
[![CI](https://github.com/nargesbyt/todo.go/actions/workflows/ci.yaml/badge.svg)](https://github.com/nargesbyt/todo.go/actions/workflows/ci.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/nargesbyt/todo.go)](https://goreportcard.com/report/github.com/nargesbyt/todo.go)

## Description
Todo is a REST API server that provides a task management service for you so that each user can define tasks.each task has status field
that shows the state of task for example pending,in progress,finished,....
There is an authentication system that prevents users who doesn't sign in in service create or retrieve tasks.


## Features
**`Todo`** supports:
- Support varius databases: PostgreSQL, SQLite3 and MySQL
- RESTful API
- JSON:API Specification
- Authentication: PAT, Basic and OIDC

## Install

### Docker images

Docker images are available on [Docker Hub](https://hub.docker.com/repository/docker/nargesbyt/todo/general).
You can launch a Todo container for trying it out with

```bash
docker run --name todo -d -p 127.0.0.1:8080:8080 nargesbyt/todo
```

Todo will now be reachable at <http://localhost:8080/>.

### Building from source

To build Todo from source code, You need:
* Go [version 1.18 or greater](https://golang.org/doc/install).

Start by cloning the repository:

```bash
go install https://github.com/nargesbyt/todo.go
```

## Contributing

## Roadmap
- [] we can add deadline to tasks that send notification to assigned user.

## License

Copyright 2023 Narges Bayat

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.