FROM debian:latest
ARG PASSWORD_GENERATED
RUN useradd -ms /bin/bash jeremy
RUN echo "jeremy:$PASSWORD_GENERATED" | chpasswd
RUN echo "root:$PASSWORD_GENERATED" | chpasswd
RUN unset PASSWORD_GENERATED
RUN apt-get update
RUN apt-get install -y rsyslog
RUN apt-get install -y sudo
RUN apt-get install -y ssh
RUN export HOSTNAME="$(sort -R /usr/share/dict/words | head -1 | tr -d \' )"
RUN sed -i 's/PermitRootLogin without-password/PermitRootLogin yes/g' /etc/ssh/sshd_config
CMD /etc/init.d/ssh start && sleep infinity
ENTRYPOINT /etc/init.d/ssh start && sleep infinity


