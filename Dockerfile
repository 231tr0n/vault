FROM archlinux:latest

RUN pacman -Syu --noconfirm
RUN pacman -Syu go --noconfirm

COPY . /root/vault

WORKDIR /root/vault

RUN go run github.com/go-task/task/v3/cmd/task@latest build

ENTRYPOINT ["./bin/vault"]

CMD ["-help"]
