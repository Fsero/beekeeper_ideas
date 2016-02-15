FROM krlmlr/debian-ssh
ENV PASSWORD_GENERATED 123456
RUN echo "docker:$PASSWORD_GENERATED" | chpasswd
RUN unset PASSWORD_GENERATED
COPY ./entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod a+x /usr/local/bin/entrypoint.sh
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
