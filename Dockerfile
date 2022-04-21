FROM ysicing/god AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src/

COPY go.mod go.mod

COPY go.sum go.sum

RUN go mod download

COPY . .

ARG GOOS=linux

ARG GOARCH=amd64

ARG CGO_ENABLED=0

RUN go build -o release/linux/amd64/plugin

FROM ysicing/debian AS gethelm

RUN curl -s -L https://get.helm.sh/helm-v3.8.2-linux-amd64.tar.gz -o /tmp/helm-linux-amd64.tar.gz && \
	mkdir -p /tmp/helm && \
	tar xzf /tmp/helm-linux-amd64.tar.gz -C /tmp/helm  --strip-components=1

FROM ysicing/shell

COPY --from=builder /go/src/release/linux/amd64/plugin /bin/drone-plugin

COPY --from=gethelm /tmp/helm/helm /usr/local/bin/helm

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh /bin/drone-plugin /usr/local/bin/helm && helm plugin install https://github.com/chartmuseum/helm-push

ENTRYPOINT ["/entrypoint.sh"]

CMD [ "/bin/drone-plugin" ]
