FROM golang:1.10-alpine AS build

LABEL maintainer="yyoshiki41@gmail.com"

# Install tools required to build the project
RUN apk add --no-cache ca-certificates \
  curl \
  git \
  make \
  rtmpdump \
  tzdata \
  ffmpeg

# Set timezone
ENV TZ "Asia/Tokyo"

WORKDIR /go/src/github.com/yyoshiki41/radigo/
COPY Makefile /go/src/github.com/yyoshiki41/radigo/
RUN curl https://glide.sh/get | sh

# These layers are only re-built when glide files are updated
COPY glide.lock glide.yaml /go/src/github.com/yyoshiki41/radigo/
# Install library dependencies
RUN make installdeps

# Copy all project and build it
# This layer is rebuilt when ever a file has changed in the project directory
COPY . /go/src/github.com/yyoshiki41/radigo/
RUN make build-4-docker


# This results in a single layer image
FROM alpine:latest AS release

# Set timezone
RUN apk add --no-cache ca-certificates tzdata ffmpeg rtmpdump
ENV TZ "Asia/Tokyo"

COPY --from=build /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=build /usr/bin/rtmpdump /usr/bin/rtmpdump
COPY --from=build /bin/radigo /bin/radigo
ENTRYPOINT ["/bin/radigo"]
CMD ["--help"]
