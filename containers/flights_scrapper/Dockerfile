FROM debian

RUN apt update \
  && apt upgrade -y \
  && apt install -y python3 python3-pip python3-pandas && pip3 install --no-cache selenium

RUN mkdir /app
WORKDIR /app/
COPY expedia.py /app/expedia.py

ENTRYPOINT ["/usr/bin/python3", "/app/expedia.py"]
