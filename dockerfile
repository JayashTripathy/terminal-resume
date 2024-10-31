FROM ubuntu:latest

ENV USERNAME="ssh-resume"
ENV GOVERSION="1.23.2"

RUN apt-get update && apt-get install -y \
    openssh-server \
    git \
    curl \
    wget \
    vim \
    sudo \
    && rm -rf /var/lib/apt/lists/*

RUN echo "HostKeyAlgorithms +ssh-rsa" >> /etc/ssh/sshd_config && \
    echo "PubkeyAcceptedKeyTypes +ssh-rsa" >> /etc/ssh/sshd_config && \
    echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config && \
    echo "PermitRootLogin yes" >> /etc/ssh/sshd_config

# Create SSH configuration directory and file
RUN mkdir -p /root/.ssh && \
    echo "Host localhost" >> /root/.ssh/config && \
    echo "    UserKnownHostsFile /dev/null" >> /root/.ssh/config

# Install Go
RUN wget https://golang.org/dl/go${GOVERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GOVERSION}.linux-amd64.tar.gz \
    && rm go${GOVERSION}.linux-amd64.tar.gz

ENV PATH=/usr/local/go/bin:${PATH}
ENV GOPATH=/root/go
ENV PATH=${GOPATH}/bin:${PATH}

WORKDIR /root/app

COPY . /root/app

# Generate host keys if they don't exist
RUN mkdir -p /var/run/sshd && \
    ssh-keygen -A

EXPOSE 22

# Start SSH service and run the Go application with proper error handling
CMD ["/bin/bash", "-c", "service ssh start && if [ -f go.mod ]; then go mod download && go run .; else go run *.go; fi"]