# Load golang image to build
FROM golang:1.14 as builder
ENV APP_USER app
ENV APP_HOME /go/src/simple-restful

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME

WORKDIR $APP_HOME
USER $APP_USER
COPY ./ .

RUN go build -o app ./cmd/apis/user-api


# Deploy execute file to simple linux server
FROM debian:buster
ENV APP_USER app
ENV APP_HOME /go/src/simple-restful

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

COPY conf/ conf/
COPY --chown=0:0 --from=builder $APP_HOME/app $APP_HOME

EXPOSE 8080
USER $APP_USER
CMD ["./app", "-cf=./conf/app.yaml"]