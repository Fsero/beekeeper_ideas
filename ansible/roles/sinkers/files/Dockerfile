FROM servicebase:0.0.1
ARG PASSWORD_GENERATED
RUN echo "jeremy:$PASSWORD_GENERATED" | chpasswd
RUN echo "root:$PASSWORD_GENERATED" | chpasswd
RUN unset PASSWORD_GENERATED
CMD /etc/init.d/ssh start && sleep infinity
ENTRYPOINT /etc/init.d/ssh start && sleep infinity
