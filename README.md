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
