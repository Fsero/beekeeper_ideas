FROM debian:8.6
MAINTAINER fsero
RUN apt-get -y update
RUN apt-get -y install sudo apt-utils
RUN apt-get -y install python-setuptools
RUN apt-get -y install python-dev
RUN apt-get -y install libffi-dev
RUN apt-get -y install build-essential
RUN apt-get -y install libssl-dev
RUN apt-get -y install python-pip
RUN apt-get -y install ssh
RUN apt-get -y install rsyslog
RUN apt-get -y install apt-transport-https ca-certificates gnupg2
RUN apt-get -y install apt-transport-https ca-certificates gnupg2
RUN apt-get -y install software-properties-common
RUN add-apt-repository 'deb https://apt.dockerproject.org/repo debian-jessie main'
RUN apt-get -y update
RUN apt-get -y install docker-engine --force-yes
RUN pip install --upgrade cffi
RUN pip install ansible
RUN mkdir -p /root/.ssh_keys/
RUN mkdir -p /root/.ssh/
COPY *.pub /root/.ssh_keys/
RUN chmod 0400 /root/.ssh_keys/*
RUN cat /root/.ssh_keys/* >> /root/.ssh/authorized_keys
RUN chmod 0400 /root/.ssh/*
COPY *.pub /root/.ssh_keys/
RUN chmod 0400 /root/.ssh_keys/*
RUN chmod 0400 /root/.ssh/*
RUN cat /root/.ssh_keys/* >> /root/.ssh/authorized_keys


CMD /etc/init.d/ssh start && bash -l
ENTRYPOINT /etc/init.d/ssh start && bash -l
