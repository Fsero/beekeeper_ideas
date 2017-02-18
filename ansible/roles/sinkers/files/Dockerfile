FROM debian:latest
ARG PASSWORD_GENERATED
RUN useradd -ms /bin/bash jeremy
RUN echo "jeremy:$PASSWORD_GENERATED" | chpasswd
RUN echo "root:$PASSWORD_GENERATED" | chpasswd
RUN unset PASSWORD_GENERATED
RUN apt-get update && \
    apt-get install -y rsyslog sudo ssh vim curl wget dbus && \
    apt-get clean
RUN wget -O /usr/local/sbin/init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64
RUN chmod +x /usr/local/sbin/init
RUN /etc/init.d/rsyslog start
RUN sed -i 's/PermitRootLogin without-password/PermitRootLogin yes/g' /etc/ssh/sshd_config
CMD /etc/init.d/ssh start && sleep infinity
ENTRYPOINT /etc/init.d/ssh start && sleep infinity