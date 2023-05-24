FROM archlinux:latest

RUN pacman -Sy go --needed --noconfirm

COPY . /root/vault

WORKDIR /root/vault

RUN go build -o /usr/local/bin/vault cmd/vault/*.go

RUN go clean -r -modcache -cache -i -x -testcache -fuzzcache

WORKDIR /root

RUN rm -rf /root/vault
RUN rm -rf /root/go

RUN pacman -Rnsu go --noconfirm
RUN pacman -Scc

ENTRYPOINT ["vault"]

CMD ["-help"]
