#!/bin/env python3

import requests, argparse, os
import time, json
from datetime import datetime

def blog(txt):
    now = datetime.now()
    dt_string = now.strftime("%d/%m/%Y %H:%M:%S")
    print("[{0}] {1}".format(dt_string, txt))

def get_args():
    parser = argparse.ArgumentParser(description="Makes a backup of the dashboards that are currently on grafana")
    parser.add_argument('-host', required=False, help='Host', default=os.environ.get('GF-HOST'))
    parser.add_argument('-token', required=False, help='API Token', default=os.environ.get('GF-TOKEN'))
    parser.add_argument('-interval', required=False, help='Interval in hours', default=os.getenv('GF-INTERVAL', 24))
    parser.add_argument('-output', required=False, help='Output folder', default="dashboards")
    args = parser.parse_args()
    args.interval = int(args.interval)
    if not args.token:
        exit(parser.print_usage())
    return args

def getDashboardList(host, token):
    url = "{}/api/search?query=&".format(host)
    r = requests.get(
        url,
        headers = { "Authorization": "Bearer {}".format(token) },
        allow_redirects=True,
    )
    if r.status_code != 200:
        blog("[{}] Error: {}", r.status_code, r.text)
        exit(1)
    return [x for x in r.json() if x["type"] == "dash-db"]

def saveDashboard(host, token, uid, folder):
    url = "{}/api/dashboards/uid/{}".format(host, uid)
    r = requests.get(
        url,
        headers = { "Authorization": "Bearer {}".format(token) },
        allow_redirects=True,
    )
    if r.status_code != 200:
        blog("Error getting dashboard: {}".format(r.text))
        return False
    obj = r.json()
    try:
        filename = "{}/{}_{}.json".format(args.output, uid, obj["meta"]["slug"])
        logfile = open(filename, "w+")
        logfile.write(json.dumps(obj, indent=4))
        logfile.close()
    except Exception as e:
        blog("Error writting to {}: {}".format(filename, e))
        return False
    return True
            
if __name__ == '__main__':
    args = get_args()
    blog("Getting dashboards from {} every {} hours".format(args.host, args.interval))
    while True:
        dashboards = getDashboardList(args.host, args.token)
        for dash in dashboards:
            x = saveDashboard(args.host, args.token, dash["uid"], args.output)
        blog("{0} dashboards copied".format(len(dashboards)))
        time.sleep(3600*args.interval)
