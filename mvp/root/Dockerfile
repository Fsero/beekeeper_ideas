FROM rastasheep/ubuntu-sshd
RUN passwd -d root
RUN sed -i 's/PermitEmptyPasswords no/PermitEmptyPasswords yes/g' /etc/ssh/sshd_config
