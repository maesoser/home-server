#!/bin/env python3

import requests, argparse, os
import time
from datetime import date
from datetime import datetime

'''
Quick and dirty script to extract logs from Cloudflare using its LogPull API periodically
'''

logFields = "botScore,CacheCacheStatus,CacheResponseBytes,CacheResponseStatus,CacheTieredFill,ClientASN,ClientCountry,ClientDeviceType,ClientIP,ClientIPClass,ClientRequestBytes,ClientRequestHost,ClientRequestMethod,ClientRequestProtocol,ClientRequestReferer,ClientRequestURI,ClientRequestUserAgent,ClientSSLCipher,ClientSSLProtocol,ClientSrcPort,EdgeColoID,EdgeEndTimestamp,EdgePathingOp,EdgePathingSrc,EdgePathingStatus,EdgeRateLimitAction,EdgeRateLimitID,EdgeRequestHost,EdgeResponseBytes,EdgeResponseCompressionRatio,EdgeResponseContentType,EdgeResponseStatus,EdgeServerIP,EdgeStartTimestamp,OriginIP,OriginResponseBytes,OriginResponseHTTPExpires,OriginResponseHTTPLastModified,OriginResponseStatus,OriginResponseTime,OriginSSLProtocol,ParentRayID,RayID,SecurityLevel,WAFAction,WAFFlags,WAFMatchedVar,WAFProfile,WAFRuleID,WAFRuleMessage,WorkerCPUTime,WorkerStatus,WorkerSubrequest,WorkerSubrequestCount,ZoneID"

def blog(txt):
    now = datetime.now()
    dt_string = now.strftime("%d/%m/%Y %H:%M:%S")
    print("[{0}] {1}".format(dt_string, txt))

def get_args():
    parser = argparse.ArgumentParser(description="Downloads logs from Clouflare using logpull")
    parser.add_argument('-email', required=False, help='User Email', default=os.environ.get('CF-EMAIL'))
    parser.add_argument('-token', required=False, help='User Token', default=os.environ.get('CF-TOKEN'))
    parser.add_argument('-zone', required=False, help='Zone ID', default=os.environ.get('CF-ZONE'))
    parser.add_argument('-interval', required=False, help='Interval in minutes', default=os.getenv('CF-INTERVAL', 30))
    parser.add_argument('-fields', required=False, help='Log Fields', default=logFields)
    parser.add_argument('-output', required=False, help='Output folder', default="logs")
    args = parser.parse_args()
    args.interval = int(args.interval)
    if not args.email or not args.token or not args.zone:
        exit(parser.print_usage())
    return args

def getLogs(zoneid, email, token, interval, fields):
    stop = int(time.time()) - 5*60
    start = stop - interval*60
    url = "https://api.cloudflare.com/client/v4/zones/{}/logs/received?start={}&end={}&fields={}".format(zoneid, start, stop, fields)
    r = requests.get(
        url,
        headers = { "X-Auth-Email": email, "X-Auth-Key": token },
        allow_redirects=True,
    )
    return r

if __name__ == '__main__':
    args = get_args()
    blog("Getting last {}m from zone {}".format(args.interval, args.zone))
    while True:
        r = getLogs(args.zone, args.email, args.token, args.interval, args.fields)
        if r.status_code != 200:
            blog("[{}] Error: {}", r.status_code, r.text)
            exit(1)
        else:
            blog("[{}] Received {} bytes".format(r.status_code, len(r.text)))
            try:
                filename = "{}/{}.log".format(args.output, date.today().strftime("%d%m%Y"))
                logfile = open(filename, "a+")
                logfile.write(r.text)
                logfile.close()
            except Exception as e:
                blog("Error writting to {}: {}".format(filename, e))
                exit(1)
        time.sleep(60*args.interval)
