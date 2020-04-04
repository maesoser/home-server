FROM python:3.5-alpine

RUN apk add --no-cache nmap nmap-scripts git

RUN pip install --no-cache-dir Flask markdown xmltodict google-cloud-storage boto3

RUN git clone https://github.com/vulnersCom/nmap-vulners /usr/share/nmap/scripts/vulners && nmap --script-updatedb
RUN mkdir /shared

COPY run.sh output_report.py latex_header.tex gcp_push.py aws_push.py app.py /
COPY flan /
RUN chmod +x /run.sh

CMD ["flask", "run", "--host", "0.0.0.0", "--port", "9999"]
