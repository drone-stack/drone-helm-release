FROM ysicing/debian AS gethelm

RUN curl -s -L https://dl.ysicing.me/helm-v3.8.0-linux-amd64.tar.gz -o /tmp/helm-linux-amd64.tar.gz && \
    mkdir -p /tmp/helm && \
    tar xzf /tmp/helm-linux-amd64.tar.gz -C /tmp/helm  --strip-components=1 

FROM ysicing/shell

COPY --from=gethelm /tmp/helm/helm /usr/local/bin/helm

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x  /usr/local/bin/helm /entrypoint.sh && helm plugin install https://gitee.com/ysbot/helm-push

ENTRYPOINT ["/entrypoint.sh"]