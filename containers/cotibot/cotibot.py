#!/usr/bin/env python3
from __future__ import unicode_literals

# coding=utf-8

import youtube_dl

import warnings
warnings.filterwarnings("ignore")

from datetime import datetime, timedelta
from telegram.ext import Updater
from telegram.ext import CommandHandler
import coloredlogs, logging
from telegram.ext import MessageHandler, Filters
from telegram.ext.dispatcher import run_async
from telegram import ParseMode as pmode

import time, sys, signal
from collections import defaultdict
from itertools import groupby
from time import sleep
import os, gc, threading, sqlite3
import urllib.request
import untangle
import subprocess

import argparse

#import madevents

DATABASE_PATH = "tgdatabase.db"
TMP_PATH = "/tmp/"

from telegram.error import (TelegramError, Unauthorized, BadRequest,
                            TimedOut, ChatMigrated, NetworkError)

def signal_handler(signal, frame):
        logging.info('Gracefully shuttting down bot')
        sys.exit(0)

def error_callback(bot, update, error):
    try:
        raise error
    except BadRequest:
        logging.error("telegram-bot: BAD REQUEST")
    except TimedOut:
        logging.error("telegram-bot: TIMEOUT")
    except NetworkError:
        logging.error("telegram-bot: NETWORK ERROR")
    except Exception as e:
        logging.error("telegram-bot: "+ str(e))

def findChats(cursor,uid):
    tables = getManagedChats(cursor)
    tables = [x for x in tables if x.find("G")!=-1]
    out = []
    for table in tables:
        users = getUsersID(cursor,table)
        if uid in users:
            out.append(table)
    return list(set(out))

def getUsersID(cursor,cid):
    idlist = cursor.execute("""
        SELECT DISTINCT userid
        FROM """+cid).fetchall()
    return [i[0] for i in idlist]

def getManagedChats(cursor):
    tables =  cursor.execute("""
    SELECT name
    FROM sqlite_master
    WHERE type = 'table'""").fetchall()
    return [x[0] for x in tables]

def memory():
    with open('/proc/meminfo', 'r') as mem:
        ret = {}
        tmp = 0
        for i in mem:
            sline = i.split()
            if str(sline[0]) == 'MemTotal:':
                ret['total'] = int(int(sline[1])/1024)
            elif str(sline[0]) in ('MemFree:', 'Buffers:', 'Cached:'):
                tmp += int(sline[1])
        ret['free'] = int(tmp/1024)
        ret['used'] = int(ret['total']) - int(ret['free'])
    return ret

def start(bot, update):
    bot.send_message(chat_id=update.message.chat_id, text="""Hi! I am CotiBot and I will take a look to your interactions in this group!""")
    bot.send_message(chat_id=update.message.chat_id, text="Stay cool, I won't keep ANYTHING of what you talked about, trust me.")
    bot.send_message(chat_id=update.message.chat_id, text="You can take a look at my graphics by using \gstats")

def echo(bot, update):
    conn = sqlite3.connect(DATABASE_PATH)
    c = conn.cursor()
    chats = getManagedChats(c)
    cid = str(update.message.chat_id).replace("-","G")
    if cid.find('G')==-1:
        cid = 'U'+str(cid)
    if not cid in chats:
        c.execute('''CREATE TABLE IF NOT EXISTS '''+cid+
            '''(
            id integer PRIMARY KEY,
            date text NOT NULL,
            length integer NOT NULL,
            userid integer NOT NULL,
            username text NOT NULL,
            reply integer NOT NULL,
            emoticon text)'''
        )
        conn.commit()

    message = update.message
    reply = "0"
    name = str(message.from_user.id)
    emoji = ''
    msjlen = 0
    if message.sticker != None and message.sticker.emoji!=None:
        emoji = message.sticker.emoji
    if message.text!=None:
        msjlen = len(message.text)
    if message.from_user.username != None:
        name = message.from_user.username
    elif message.from_user.first_name != None:
        name = message.from_user.first_name
    if(message.reply_to_message!=None):
        reply = message.reply_to_message.message_id
    line = [
        str(message.message_id),
        str(message.date),
        str(msjlen),
        str(message.from_user.id),
        str(name),
        str(reply),
        str(emoji)
        ]
    try:
        c.execute("INSERT INTO "+cid+" VALUES (?,?,?,?,?,?,?)",line)
        conn.commit()
    except Exception as err:
        logging.error(err)
    conn.close()

def graphScheduler(mins):
    logging.info("Asynchronous graph generation set to "+str(mins)+" mins")
    sleep(10)
    while(True):
        logging.info("Calling gengraphs script...")
        subprocess.run(["./gengraphs.py", ""])
        sleep(mins*60)

def sendStats(bot, fromChat, toChat):
    cid = str(fromChat).replace("-","G")

    if os.path.isfile(TMP_PATH+cid+"_delta120.jpg") == False:
        logging.warning("Maybe there is no graphs for group "+str(cid))
    try:
        bot.send_document(chat_id=toChat,text="Overall stats",document=open("/tmp/"+cid+"_delta120.jpg","rb"))
    except Exception as e:
        logging.exception("poster 120"+str(e))
    try:
        bot.send_document(chat_id=toChat,text="Last 15 days stats",document=open("/tmp/"+cid+"_delta15.jpg","rb"))
    except Exception as e:
        logging.exception("poster 15"+str(e))
    try:
        bot.send_photo(chat_id=toChat,text="Overall frequency",photo=open("/tmp/"+cid+"_freq_delta120.jpg","rb"))
    except Exception as e:
        logging.exception("freq 120"+str(e))
    try:
        bot.send_photo(chat_id=toChat,text="Last 15 days frequency",photo=open("/tmp/"+cid+"_freq_delta15.jpg","rb"))
    except Exception as e:
        logging.exception("freq 15"+str(e))

@run_async
def stats(bot,update):
    logging.info("Sending graphs to user "+update.message.from_user.first_name+"\t"+str(update.message.from_user.id))
    bot.send_message(chat_id=update.message.chat_id, text="Let me do some numbers, "+update.message.from_user.first_name)
    if update.message.chat.type=='private':
        conn = sqlite3.connect(DATABASE_PATH)
        c = conn.cursor()
        chats = findChats(c, update.message.from_user.id)
        conn.close()
        for chat in chats:
            sendStats(bot,chat,update.message.chat_id)
        if len(chats)==0:
            bot.send_message(chat_id=update.message.chat_id, text="I'm not monitoring any of your groups!")

    else:
        sendStats(bot,update.message.chat_id,update.message.chat_id)

def aliveresp(bot,update):

    conn = sqlite3.connect(DATABASE_PATH)
    c = conn.cursor()
    chatlen = len(getManagedChats(c))
    conn.close()

    texttosend = "Ey, I\'m here! And I\'m managing {0} group chats".format(chatlen)
    bot.send_message(chat_id=update.message.chat_id, text=texttosend)

def music(url,uid):
    class MyLogger(object):
        def debug(self, msg):
            pass
        def warning(self, msg):
            pass
        def error(self, msg):
            logging.info(str(msg))

    def my_hook(d):
        if d['status'] == 'finished':
            logging.info('Done downloading, now converting ' + d['filename'])

    dirname = TMP_PATH + "yd_" + str(int(time.time())) + str(uid)
    if not os.path.exists(dirname):
        os.mkdir(dirname)
        logging.info("Directory " + dirname + " Created ")
    else:
        logging.info("Directory " + dirname + " already exists")

    ydl_opts = {
        'format': 'bestaudio/best',
        'outtmpl': dirname + '/%(title)s.%(ext)s',
        'restrictfilenames' : True,
        'noplaylist' : True,
        'postprocessors': [
            {
              'key': 'FFmpegExtractAudio',
              'preferredcodec': 'mp3',
              'preferredquality': '192',
            },
            {
              'key' : 'FFmpegMetadata',
            }
        ],
        'logger': MyLogger(),
        'progress_hooks': [my_hook],
    }
    with youtube_dl.YoutubeDL(ydl_opts) as ydl:
        ydl.download([url])

    mp3files = os.listdir(dirname)
    logging.info(mp3files)
    for mp3file in mp3files:
        name = mp3file.replace(".mp3","").replace("_","")
    return dirname+"/"+mp3files[0]

def asyncmusic(bot, url, chatid):
    audiofile = music(url, chatid)
    bot.send_audio(chat_id=chatid, audio=open(audiofile, 'rb'))
    logging.info("Music file sent")
    if os.path.exists(audiofile):
        os.remove(audiofile)

def getmusic(bot, update):
    url = update.message.text[6:].replace(" ","")
    if url.find("http") == -1:
        bot.send_message(chat_id=update.message.chat_id, text="This does not seems to be an url...")
    else:
        bot.send_message(chat_id=update.message.chat_id, text="Ok! I'll send it to you when I finish!!")
        t = threading.Thread(target=asyncmusic,name='Music',args=(bot, url,update.message.chat_id))
        t.start()

def get_args():
    parser = argparse.ArgumentParser(description="This bot makes statistics with the messages sent on a telegra group")
    parser.add_argument('--tgtoken', required=False, help='Telegram token', default=os.getenv('TG_TOKEN', ""))
    parser.add_argument('--log', required=False, help='Default logging file', default=os.getenv('LOG_PATH', "/tmp/cotibot.log"))

    args = parser.parse_args()
    return args

if __name__ == "__main__":

    args = get_args()

    logging.basicConfig(
        format='%(asctime)s - %(levelname)s - %(message)s',
        level=logging.INFO,
        filename=args.log
        )
    coloredlogs.DEFAULT_LOG_FORMAT = '%(asctime)s - %(levelname)s - %(message)s'
    coloredlogs.install()

    logging.info("Initiating bot...")

    signal.signal(signal.SIGINT, signal_handler)

    t = threading.Thread(target=graphScheduler,name='GraphSchd',args=(180,))
    t.start()

    updater = Updater(token=args.tgtoken, workers=2,request_kwargs={'read_timeout': 6, 'connect_timeout': 7})
    dispatcher = updater.dispatcher

    msg_handler = MessageHandler((~Filters.command), echo)
    dispatcher.add_handler(msg_handler)

    stat_handler = CommandHandler('gstats', stats)
    dispatcher.add_handler(stat_handler)

    music_handler = CommandHandler('music', getmusic)
    dispatcher.add_handler(music_handler)

    alive_handler = CommandHandler('alive', aliveresp)
    dispatcher.add_handler(alive_handler)

    logging.info("Bot started")
    dispatcher.add_error_handler(error_callback)

    updater.start_polling()
    updater.idle()
