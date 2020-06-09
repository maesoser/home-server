FROM python

ENV PYTHONUNBUFFERED=0
RUN pip install --no-cache-dir requests
WORKDIR /
COPY grafbck.py .
ENTRYPOINT ["/usr/local/bin/python3", "-u", "grafbck.py"]
