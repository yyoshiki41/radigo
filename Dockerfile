FROM yyoshiki41/ubuntu-ffmpeg-v3.3

LABEL maintainer="yyoshiki41@gmail.com"

ENV VERSION "v0.5.1"

RUN apt update
RUN apt install -y tzdata wget unzip

# set timezone
ENV TZ "Asia/Tokyo"
RUN echo $TZ > /etc/timezone
RUN dpkg-reconfigure --frontend noninteractive tzdata

# download radigo
WORKDIR /tmp
RUN wget https://github.com/yyoshiki41/radigo/releases/download/${VERSION}/radigo_${VERSION}_linux_amd64.zip
RUN unzip ./radigo_${VERSION}_linux_amd64.zip -d /usr/local/bin

RUN mkdir -p /tmp/radigo/output

CMD ["/bin/bash"]
