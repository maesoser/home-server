FROM python

ENV PYTHONUNBUFFERED=0
RUN pip install --no-cache-dir requests
WORKDIR /
COPY logpull.py .
ENTRYPOINT ["/usr/local/bin/python3", "-u", "logpull.py"]
