#! /usr/bin/env python3

import feedparser
import requests
import pickledb
import httplib2
from sys import argv, exit
import json, sys
import time, datetime

def doFilter(f, title):
    for item in f:
        item = item.lower()
        if title.find(item) == -1:
            return False
    return True

def filter(filters, title):
    for filter in filters:
        title = title.lower()
        if doFilter(filter, title):
            return True
    return False

def loadFilter(filepath):
    filters = []
    with open(filepath) as f:
        for line in f:
            if line[0] != "#":
                filters.append(line.replace("\n","").split(","))
    return filters

def getTorrents(url, filters, dbpath):
    url = url.replace("\n","")
    added = 0
    db = pickledb.load(dbpath, False)
    feedparser.USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:53.0) Gecko/20100101 Firefox/53.0"
    NewsFeed = feedparser.parse(url)
    for entry in NewsFeed.entries:
        infohash = entry.get('infohash')
        if infohash == None:
            infohash = entry.get("id").split("/")[-1]
        if infohash == None:
            infohash = entry.get("torrent:infoHash")

        uri = entry.get("magneturi")
        if uri == None:
            uri = entry.get("torrent:magnetURI")
        if uri == None:
            uri = entry.get("link")

        title = entry["title"]
        if not db.exists(infohash):
            object = {
                'title' : title,
                'uri' : uri
            }
            db.set(infohash, object)
            added +=1
            if filter(filters, title):
                f = open("/data/downloads.csv", "a+")
                f.write(title+";"+uri+"\n")
                f.close()
    db.dump()
    return added

if len(sys.argv) != 2:
    print("Usage: {0} [workdir]".format(sys.argv[0]))
    exit(1)

path = sys.argv[1]
delayHours = 6

while True:
   filters = loadFilter(path + "/filters")
   download = 0
   with open(path + '/sources') as f:
      for line in f:
          if line[0] != "#":
             d = datetime.date.today()
             dbname = 'torrents-{0:02d}-{1}.db'.format(d.month,d.year)
             try:
                print("Downloading from: {0}".format(line))
                d = getTorrents(line, filters, path +"/"+dbname)
                download += d
             except Exception as e:
                print(e)
   print("{0} new torrents".format(download))
   time.sleep(3600*delayHours)
