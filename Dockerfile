FROM golang:1.14-alpine AS build

LABEL maintainer="yyoshiki41@gmail.com"

# Set timezone
ENV TZ "Asia/Tokyo"

# Install tools required to build the project
RUN apk add --no-cache ca-certificates \
  curl \
  ffmpeg \
  git \
  make \
  rtmpdump \
  tzdata

ARG PROJECT_PATH=/go/src/github.com/yyoshiki41/radigo
WORKDIR ${PROJECT_PATH}
COPY . ${PROJECT_PATH}/

# Install deps
RUN make installdeps
# Build the project binary
RUN make build-4-docker


# This results in a single layer image
FROM alpine:latest

# Set timezone
ENV TZ "Asia/Tokyo"
# Set default output dir
VOLUME ["/output"]

RUN apk add --no-cache ca-certificates ffmpeg rtmpdump tzdata

COPY --from=build /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=build /usr/bin/rtmpdump /usr/bin/rtmpdump
COPY --from=build /bin/radigo /bin/radigo

ENTRYPOINT ["/bin/radigo"]
CMD ["--help"]
