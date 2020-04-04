FROM python:3.7.2-alpine3.9

WORKDIR /app

RUN python3 -m pip install feedparser pickledb httplib2 requests

COPY torrent.py .
CMD [ "python3", "/app/torrent.py", "/data" ]
