FROM golang:1.7

WORKDIR /go

# NOTE: Everything must already have been built outside the container
COPY build/picaxe su-exec ./

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

EXPOSE 3000

ENTRYPOINT ["./su-exec", "app", "./picaxe", "--listen", ":3000"]
