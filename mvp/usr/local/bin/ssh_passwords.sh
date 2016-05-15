#!/usr/bin/env bash

OPTIONS="-v -A -s 4096 -pc --unbuffered -F -j"
mkdir -p /var/log/ssh
#contrasenyas
nohup sysdig $OPTIONS fd.num=4 and evt.is_io_write=true and proc.name=sshd > /var/log/ssh/passwords.json  2> /dev/null &

#paquetes entrantes
sysdig $OPTIONS proc.name=sshd and fd.port!=22 and syscall.type=recvto > /var/log/ssh/inconnections.json 2> /dev/null &


#paquetes salientes
sysdig $OPTIONS proc.name=sshd and fd.port!=22 and syscall.type=sendto > /var/log/ssh/outconnections.json 2> /dev/null &

#comandos
sysdig $OPTIONS -c spy_users > /var/log/ssh/command.json 2> /dev/null &

#files changes in /root or in /usr/local/bin
#input/output

sysdig $OPTIONS syscall.type in \('open','close','stat','read','write'\) and \(fd.type=file or fd.type=directory\) and fd.name != /dev/ptmx > /var/log/ssh/io.json

