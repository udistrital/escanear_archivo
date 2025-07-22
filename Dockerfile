FROM alpine:latest

RUN yum install -y clamav clamav-update clamav-lib && \
    yum clean all && \
    rm -rf /var/cache/yum

RUN freshclam || true

COPY --from=builder /app/main /main
COPY entrypoint.sh /entrypoint.sh
COPY conf/app.conf /conf/app.conf

RUN chmod +x /main /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]