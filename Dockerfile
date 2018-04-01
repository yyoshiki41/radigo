FROM golang:1.10-alpine AS build

LABEL maintainer="yyoshiki41@gmail.com"

ARG PROJECT_PATH=/go/src/github.com/yyoshiki41/radigo

# Set timezone
ENV TZ "Asia/Tokyo"

# Install tools required to build the project
RUN apk add --no-cache ca-certificates \
  curl \
  git \
  make \
  rtmpdump \
  tzdata \
  ffmpeg

WORKDIR ${PROJECT_PATH}
COPY . ${PROJECT_PATH}/

RUN curl https://glide.sh/get | sh
RUN make build-4-docker


# This results in a single layer image
FROM alpine:latest

# Set timezone
ENV TZ "Asia/Tokyo"

RUN apk add --no-cache ca-certificates tzdata ffmpeg rtmpdump

COPY --from=build /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=build /usr/bin/rtmpdump /usr/bin/rtmpdump
COPY --from=build /bin/radigo /bin/radigo
ENTRYPOINT ["/bin/radigo"]
CMD ["--help"]
