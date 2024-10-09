FROM ubuntu:latest

ENV USERNAME "ssh-resume"
ENV GOVERSION="1.21.3"


RUN apt-get update && apt-get install -y \
    openssh-server \
    git \
    curl \
    wget \
    vim \
    sudo \
    && rm -rf /var/lib/apt/lists/*

RUN echo "HostKeyAlgorithms +ssh-rsa" >> /etc/ssh/sshd_config
RUN echo "PubkeyAcceptedKeyTypes +ssh-rsa" >> /etc/ssh/sshd_config
RUN echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config
RUN echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
    

# Configure SSH
RUN echo "HostKeyAlgorithms +ssh-rsa" >> /etc/ssh/sshd_config
RUN echo "PubkeyAcceptedKeyTypes +ssh-rsa" >> /etc/ssh/sshd_config

# Install Go
RUN wget https://golang.org/dl/go${GOVERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GOVERSION}.linux-amd64.tar.gz \
    && rm go${GOVERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/root/go"
ENV PATH="${GOPATH}/bin:${PATH}"

WORKDIR "/root/app"

COPY . /root/app

EXPOSE 22
CMD go run *.go

