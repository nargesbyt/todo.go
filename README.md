# Todo
[![Go Report Card](https://goreportcard.com/badge/github.com/nargesbyt/todo.go)](https://goreportcard.com/report/github.com/nargesbyt/todo.go)

## Description
Todo is a REST API server that provides a task management service for you so that each user can define tasks.each task has status field
that shows the state of task for example pending,in progress,finished,....
There is an authentication system that prevents users who doesn't sign in in service create or retrieve tasks.


##Features
**`Todo`** supports:
- http,https protocol
- Postgres,sqlite3,Mysql databases
- RESTful API
- JSON:API Specification
- Token based Authentication
- Basic Authentication

## Install

### Docker images

Docker images are available on [Docker Hub](https://hub.docker.com/repository/docker/nargesbyt/todo/general).
You can launch a Todo containerfor trying it out with

```bash
docker run --name todo -d -p 127.0.0.1:8080:8080 todo
```

Todo will now be reachable at <http://localhost:8080/>.

### Building from source

To build Todo from source code, You need:
* Go [version 1.18 or greater](https://golang.org/doc/install).

Start by cloning the repository:

```bash
git clone https://github.com/nargesbyt/todo.go.git
cd todo
```

## Contributing

## Roadmap
we can add deadline to tasks that send notification to assigned user.

## License