FROM debian

RUN apt update \
  && apt upgrade -y \
  && apt install -y python3 python3-pip python3-pandas python3-numpy python3-matplotlib python3-seaborn python3-networkx ffmpeg \
  && pip3 install --no-cache youtube-dl python-telegram-bot BeautifulSoup4 coloredlogs untangle

RUN mkdir /app
WORKDIR /app/
COPY . .

ENTRYPOINT ["/usr/bin/python3 cotibot.py"]
