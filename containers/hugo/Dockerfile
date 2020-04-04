FROM arm32v7/debian

RUN apt-get update && apt-get install -y wget libsass1 &&\
    wget -O /tmp/hugo.deb https://github.com/gohugoio/hugo/releases/download/v0.67.1/hugo_0.67.1_Linux-ARM.deb &&\
    dpkg -i /tmp/hugo.deb

COPY run.sh /run.sh
RUN chmod 0777 /run.sh
WORKDIR /
CMD ["/run.sh"]
