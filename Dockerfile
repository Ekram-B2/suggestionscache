FROM alpine:3.10.3

RUN apk add tzdata ca-certificates
RUN addgroup -S gvlcache && adduser -s /bin/bash -D -S -G gvlcache gvlcache
COPY . /gvlcache
RUN chown root:gvlcache /gvlcache && \
  chmod 750 /gvlcache

USER gvlcache

ENTRYPOINT [ "/gvlcache/build/gvlcache"]