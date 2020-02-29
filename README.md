## Presentation

Golem is a tool that creates development servers from a config file.
Basically, you can create a static files server, a JSON database on the fly or routes based web server


## Getting started 

> Golem is a draft, there no pre-compiled executables yet.
> To use it you have to build it from sources.

To start a new Golem server you have to create a config file named golem.yaml
This file should contain at least the following entries:

* __port:__ the main port used by the server
* __services:__ a list of one or more services mounted by Golem

Each service should contain: 

* A name
* A type 
* A configuration

There is an example of a server that's return a JSON object 

```
port: "7171"
services: 
  - name: "Ping Server"
    http_config: 
      routes:
        "/ping": 
          body: '{"message": "pong!"}'
```

From this configuration file, Golem will start a server listening on the port `7171`

You also use parameters on the route declaration and retrieve them as template variables: 

```
port: "7171"
services: 
  - name: "Move"
    http_config: 
      routes:
        "/ping": 
          body: '{"message": "pong!"}'
        "/echo/:message":
          handler:
            type: "template"
            template: "${message}"
```

As you can see, to serve a static response we can simply use the `body` attribute.
To use a template as response we have to use the `handler` attribute with two childs: 

* __type:__ to specicy the handler type. `template` in this case
* __template:__ contains the template value
