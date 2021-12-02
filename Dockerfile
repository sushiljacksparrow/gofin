FROM golang:alpine

# Add Maintainer Info
LABEL maintainer="Sushil Singh"
RUN apk add --no-cache --virtual .build-deps \
        bzip2 \
        curl \
        g++ \
        gcc \
        bash \
      cmake \
      sudo \
		libssh2 libssh2-dev\
		git

RUN mkdir -p /app/pkg
RUN mkdir -p /app/resources
COPY ./pkg /app/pkg/
#COPY ./resources /app/resources/
# COPY ./<GOOGLE_CREDS_FIlE> /app/
RUN ls /app/pkg/
RUN ls /app/resources

WORKDIR /app/
#ENV GOPATH=$GOPATH:/go/
RUN echo $GOPATH
# Copy go mod and sum files
COPY go.mod /app/

RUN go build -o gofin pkg/main/server.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./gofin"]