# ansible-docker

Ansible role to install Docker Engine via the official Docker apt repositories
for Debian-family distributions. This role assumes a 64-bit installation and a
kernel supported by Docker. See the Docker install guides for
[Debian](http://docs.docker.com/engine/installation/debian/) or
[Ubuntu](http://docs.docker.com/engine/installation/ubuntulinux/) for more.

## Requirements

* Ansible 2.0+
* Debian-family distribution (current stable/LTS or newer)

## Dependencies

###[apt](https://github.com/cspicer/ansible-apt)
The `apt` role is used to add the official Docker repo along with signing keys.

## Variables

See [`defaults/main.yml`](defaults/main.yml) for default values.

### docker

These variables should be set either as a dependency in `meta/main.yml` for
your role, or as a part of your include statement. See below for examples.

Variable        | Type        | Description
--------        | ----        | -----------
`docker_pkg`    | List        | List of apt packages to be installed

## Testing

This role includes tests which are run via [Docker](https://www.docker.com)
along with a Makefile to simplify the testing process.

External role dependencies are defined in a [requirements](requirements.yml)
file which is used by the [`ansible-galaxy`](http://docs.ansible.com/ansible/galaxy.html#the-ansible-galaxy-command-line-tool)
command line tool. Roles are copied into [`tests/roles`](tests/roles) which is
configured in [`ansible.cfg`](ansible.cfg) as a search path.

A [small shell script](tests/test.sh) is used to resolve role dependencies
and run Ansible within the container.

Note that because this role installs Docker within a Docker container,
privledged mode is required by the top level container. This may have security
implications depending on your setup.

Starting a test container:
* `make test`

Cleaning up:
* `make clean`

## Examples

Add to your playbook as a role include:

```yaml
---
- hosts: docker-hosts
  roles:
    - ansible-docker
```

## Development

* Source hosted at [GitHub][repo]
* Report issues/questions/feature requests on [GitHub Issues][issues]

Pull requests are very welcome! Make sure your patches are well tested.
Ideally create a topic branch for every separate change you make. For
example:

1. Fork the repo
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## Author Information

Created and maintained by [Chris Spicer][cspicer] (<github@cspicer.ca>).

## License

MIT License (see [LICENSE][license])

[cspicer]: https://github.com/cspicer
[repo]: https://github.com/cspicer/ansible-docker
[issues]: https://github.com/cspicer/ansible-docker/issues
[license]: https://github.com/cspicer/ansible-docker/blob/master/LICENSE
