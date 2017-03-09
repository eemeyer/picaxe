FROM us.gcr.io/t11e-platform/base-go

RUN mkdir -p /go/src/github.com/t11e/picaxe \
  && chown -R app:app /go/src/github.com/t11e/picaxe

USER app
WORKDIR /go/src/github.com/t11e/picaxe

COPY glide.yaml glide.lock ./
RUN glide install

COPY . ./
RUN \
     go build -o /srv/picaxe github.com/t11e/picaxe \
  && go test $(go list github.com/t11e/picaxe/... | fgrep -v /vendor)

USER root
RUN rm -rf /go && chown root:root /srv/*

WORKDIR /srv
EXPOSE 3000
ENTRYPOINT ["su-exec", "app", "/srv/picaxe", "--listen", ":3000"]
