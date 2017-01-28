#!/usr/bin/env bash
PASSWORD_GENERATED=${1:-123456}
echo "root:$PASSWORD_GENERATED" | chpasswd
unset PASSWORD_GENERATED

