FROM golang:1.7

WORKDIR /go

USER root

# NOTE: Everything must already have been built outside the container
COPY build/picaxe .

RUN \
   addgroup --gid 9000 app \
&& adduser \
  --uid 9000 \
  --home $PWD \
  --ingroup app \
  --disabled-password \
  --no-create-home \
  --disabled-login \
  --gecos '' \
  --quiet \
  app

USER app

EXPOSE 3000

ENTRYPOINT ["./picaxe", "--listen", ":3000"]
