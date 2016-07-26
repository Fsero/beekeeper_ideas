#!/bin/bash

echo "* Ensuring git is installed as ansible-galaxy requires it"
apt-get -qq update
apt-get -qq install \
  --no-install-recommends \
  git-core

echo "* Installing role dependencies to temporary directory"
ansible-galaxy install \
  --force \
  --roles-path=tests/roles/ \
  --role-file=requirements.yml

echo "* Running ansible tests"
ansible-playbook \
  --inventory-file /etc/ansible/hosts \
  --connection local \
  tests/test.yml

