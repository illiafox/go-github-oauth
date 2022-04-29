FROM golang:1.18-alpine as build-env
RUN apk --no-cache add git
ADD . /server
RUN cd /server/cmd/app && go build -o bin

FROM alpine
WORKDIR /server
COPY --from=build-env /server/sql /server/sql
VOLUME /sql

FROM alpine
WORKDIR /server
COPY --from=build-env /server/cmd /server/cmd
COPY --from=build-env /server/web /server/web

WORKDIR /server/cmd/app
ENTRYPOINT ./bin $ARGS