Bach
====

Bach is a set of low level tools for composing services. Bach aims to
make it easy to wrap other tools, adding service discovery, env
injections, configuration management and anything else that is
necessary to support legacy applications on a modern cloud / container
base platform.


Withenv
-------

Withenv takes YAML files and applies them to the environment prior to
running a command. The idea is that rather than relying on shell
variables being set, you can explicitly define what environment
variables you want to use.

For example, here is some YAML that might be helpful when connecting
to an OpenStack Cloud.

```yaml
---
OS_USERNAME: user
OS_PASSWORD: secret
OS_URL: http://api.example.org/rest/api
```

You can then run commands (ie in Ansible for example) that read this
information.

```bash
$ we -e test_user_creds.yml ansible-play install-app
```

ToConfig
--------

Toconfig allows you to write a configuration file before starting a
process based on a template. The configuration format allows pulling
values from environment variables and inject them in your config
file.

```bash
$ DBURL=db.example.com toconfig \
    --template /etc/templates/mydb.cnf.tmpl \
	--config /etc/db/mydb.cnf \
	startdb
```

Present
-------

Present allows running a script before a process starts and another
when it finishes. For example, if you wanted to start up some process
like `etcd` before starting our app, you could write a small bash
script that starts it and stops when the program exits.

```bash
$ present --start dbstart.sh --after dbkill.sh myapp
```

Present can be used as a simple cluster management system. Present can
utilize a
[gossip protocol](http://www.cs.cornell.edu/%7Easdas/research/dsn02-swim.pdf)
to communicate between the cluster. This doesn't require a data store,
so it may not be viable for some applications, but for simple service
discovery, it should work well.

```bash
$ $SOME_CLUSTER_ADDRESS=172.18.0.4:7001 we -s 'present nodes $SOME_CLUSTER_ADDRESS' \
    toconfig -t my.conf.tmpl -s my.conf myapp
```

The `present nodes` command can find the nodes in the cluster and
print a JSON file. Withenv can then read the result, add it to the
enviornment, at which point, toconfig can write values in the config
file.

To join the cluster, you can use the default command.

```bash
$ present --name appserver --cluster-addr 172.18.0.4 myapp
```

This joins the cluster with the service name of `appserver` that can
be used by other cluster members to discover the service. At this
point there is nothing in place to do things like register specific
ports, dns lookups or anything else terribly special.

Why?
---

There are a lot of great tools out there that can run or build
applications, but with each tool, there become new requirements on how
you run it. You can use a tool, if you meet these sets of
requirements. This is generally not a huge deal, but it often means
changing the way your application works in order to accept the new
state your application needs to run. It is very difficult to adopt new
deployment, service discovery and configuration systems without
changing the way your data accepts state.

For example, lets say you have an app that needs a URL of some service
in order to run. You want to run the app in a docker container on some
cluster management system. The cluster management system allows you to
set an environment variable defining your service URL, but in order
for you to do that, you need to rewrite some code in you application
in order read the environment variable.

The other option is to use something like consul-template to rewrite
the config based some script or something similar. This adds a level
of complexity as you need to figure out a way to run the service
discovery database.

All of this is needed to try out your new deployment strategy.

A better tact is to assume no changes your app. So, step one is to
write your config with your services at run time before starting the
command.

```bash
$ SERVICE_URL=http://example.com/api/ toconfig -t myapp.conf.tmpl -s myapp.conf myapp
```

Now we can test the service makes it to our config without changing
the app.

Next up we want to look up where the service is. We can write a script
to do that.

```python
#!/usr/bin/python

import requests

resp = requests.get('http://appprovider.services/foo')
print('SERVICE_URL: %s' % resp.json()['url'])
```

Now we use withenv to load this in our environment.

```bash
$ we -s findfoo.py toconfig -t myapp.conf.tmpl -s myapp.conf myapp
```

Again we can test each portion and be confident that each step should
work.

Finally, we can use this chain of commands in our Dockerfile and test
that.

```
# Dockerfile
CMD ["we", "-s", "findfoo.py", "toconfig", "-t", "myapp.conf.tmpl", "-s", "myapp.conf", "myapp"]
```

Now we can be confident that our app is the same and it should work
with the contract that the deployment system requires.
