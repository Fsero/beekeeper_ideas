# apt

Ansible role for managing apt repositories and keeping the local apt cache up to date. By default, no new repositories are added unless `apt_repos` is defined.

## Dependencies

None

## Platforms

### Debian

* wheezy
* jessie

### Ubuntu

* precise
* trusty

## Variables

See [`defaults/main.yml`](defaults/main.yml) for default values.

These variables should be set either as a dependency in `meta/main.yml` for your role, or as a part of your playbook. See below for examples.

Variable            | Type      | Description
--------            | ----      | -----------
apt_pkg             | List      | Packages to be installed to support apt module
apt_cache_time      | Integer   | Time in seconds before apt cache is considered stale

### apt_repo

Variable        | Type        | Description
--------        | ----        | -----------
repo_url        | String      | URI of repository to be added
key_id          | String      | ID of GPG key for repo
key_server      | String      | Keyserver to use for downloading key

## Tasks

### main

- Adds repositories specified in `apt_repos` variable list
- Updates apt-cache if older than value for `apt_cache_time`
- Installs Python libraries used by Ansible for managing apt

## Examples

Add to your role as a dependency in `meta/main.yml` with the repository to be added as a YAML list:

    ---
    dependencies:
      - role: apt
        apt_repo:
          - repo_url: "deb https://your-apt-repo.tld stable main"
            key_id: "012345678"
            key_url: "https://your-apt-repo.tld/repo.gpg.key"

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
[repo]: https://github.com/cspicer/ansible-apt
[issues]: https://github.com/cspicer/ansible-apt/issues
[license]: https://github.com/cspicer/ansible-apt/blob/master/LICENSE
