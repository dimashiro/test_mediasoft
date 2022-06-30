FROM golang:1.17 as builder
ENV CGO_ENABLED 0
ARG BUILD_REF

COPY . /src

#build binary
WORKDIR /src/app/service
RUN go build

# Run
FROM alpine:3.14
ARG BUILD_DATE
ARG BUILD_REF
COPY --from=builder /src/app/service/service /service/service
WORKDIR /service
EXPOSE 3000
CMD [ "./service" ]