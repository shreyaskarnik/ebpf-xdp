FROM ubuntu:22.04
ENV GO_VERSION=1.18.3
RUN apt-get update && apt-get install -y build-essential curl llvm libelf-dev clang
# download and install go
RUN curl -O https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
RUN rm go${GO_VERSION}.linux-amd64.tar.gz
RUN echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
ENV PATH $PATH:/usr/local/go/bin
# install dependencies for ebpf
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive && \
    apt-get install --no-install-recommends -y \
    autoconf bison cmake dkms flex gawk gcc python3 rsync \
    libiberty-dev libncurses-dev libpci-dev libssl-dev libudev-dev \
    && curl https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-5.13.tar.gz | tar -xz \
    && make -C /linux-5.13 headers_install INSTALL_HDR_PATH=/usr \
    && make -C /linux-5.13/tools/lib/bpf install INSTALL_HDR_PATH=/usr \
    && make -C /linux-5.13/tools/bpf/bpftool install \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /linux-5.13

RUN bpftool --version
RUN go version
