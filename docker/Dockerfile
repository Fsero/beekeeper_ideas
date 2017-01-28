FROM debian:latest
ARG PASSWORD_GENERATED
RUN echo "root:$PASSWORD_GENERATED" | chpasswd
RUN unset PASSWORD_GENERATED
RUN apt-get update
RUN apt-get install -y ssh
CMD /etc/init.d/ssh start && bash -l
ENTRYPOINT /etc/init.d/ssh start && bash -l


