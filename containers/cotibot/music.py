#!/usr/bin/env python3
from __future__ import unicode_literals

# coding=utf-8

# youtube-dl --extract-audio --audio-format mp3 -o "%(title)s.%(ext)s" $argv[1]

import os
import time
import youtube_dl

def music(url,uid):
    class MyLogger(object):
        def debug(self, msg):
            pass

        def warning(self, msg):
            pass

        def error(self, msg):
            print(msg)


    def my_hook(d):
        if d['status'] == 'finished':
            print('Done downloading, now converting ',d['filename'])


    dirname = "/tmp/yd_" + str(int(time.time())) + str(uid)
    if not os.path.exists(dirname):
        os.mkdir(dirname)
        print("Directory " , dirname ,  " Created ")
    else:    
        print("Directory " , dirname ,  " already exists")

    ydl_opts = {
        'format': 'bestaudio/best',
        'outtmpl': dirname + '/%(title)s.%(ext)s',
        'restrictfilenames' : True,
        'postprocessors': [{
            'key': 'FFmpegExtractAudio',
            'preferredcodec': 'mp3',
            'preferredquality': '192',
        }],
        'logger': MyLogger(),
        'progress_hooks': [my_hook],
    }
    with youtube_dl.YoutubeDL(ydl_opts) as ydl:
        ydl.download([url])

    return dirname+"/"+os.listdir(dirname)[0]

filename = music("https://www.youtube.com/watch?v=BaW_jenozKc", "test")
print(filename)
