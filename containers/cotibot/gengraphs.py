#!/usr/bin/env python3

# coding=utf-8
import warnings
warnings.filterwarnings("ignore")
import coloredlogs, logging

import matplotlib
matplotlib.use("Agg")

import matplotlib.pyplot as plt
import numpy as np
import seaborn as sb
import pandas as pd
import networkx as nx
import time, sys, signal
from collections import defaultdict
from itertools import groupby
from datetime import datetime, timedelta
import matplotlib.patches as mpatches
from matplotlib.dates import DateFormatter
from time import sleep
import os, gc, threading, sqlite3
import urllib.request
import untangle
from bs4 import BeautifulSoup as soup
import subprocess

def signal_handler(signal, frame):
        logging.info('Gracefully shuttting down bot')
        sys.exit(0)

days = 'Lun Mar Mie Jue Vie Sab Dom'.split()

colors = ['#a6cee3','#1f78b4','#b2df8a','#33a02c','#fb9a99','#e31a1c','#fdbf6f','#ff7f00','#cab2d6','#6a3d9a','#ffff99','#b15928',
         '#a6cee3','#1f78b4','#b2df8a','#33a02c','#fb9a99','#e31a1c','#fdbf6f','#ff7f00','#cab2d6','#6a3d9a','#ffff99','#b15928',
         '#a6cee3','#1f78b4','#b2df8a','#33a02c','#fb9a99','#e31a1c','#fdbf6f','#ff7f00','#cab2d6','#6a3d9a','#ffff99','#b15928',
         '#a6cee3','#1f78b4','#b2df8a','#33a02c','#fb9a99','#e31a1c','#fdbf6f','#ff7f00','#cab2d6','#6a3d9a','#ffff99','#b15928',
         '#a6cee3','#1f78b4','#b2df8a','#33a02c','#fb9a99','#e31a1c','#fdbf6f','#ff7f00','#cab2d6','#6a3d9a','#ffff99','#b15928'
         ]

def getStartDate(cursor, cid):
    lines =  cursor.execute("""
        SELECT min(date) 
        FROM """+cid+""" 
        GROUP BY strftime('%Y-%m-%d',date)""").fetchall()
    return datetime.strptime(lines[0][0], '%Y-%m-%d %H:%M:%S').date()

def deleteUnused(cursor):
    chats = getManagedChats(cursor)
    for chat in chats:
        users = len(getUsersID(cursor,chat))
        msgs = len(getChat(cursor,chat,120))
        if users < 2 or msgs < 20:
            logging.warning("Chat "+str(chat)+" has only "+str(users)+" users")
            logging.warning("Chat "+str(chat)+" has published just "+str(msgs)+" msgs in the last 120 days")
            cursor.execute("DROP TABLE "+str(chat))
            logging.warning("DELETED "+str(chat))

def getManagedChats(cursor):
    tables =  cursor.execute("""
    SELECT name 
    FROM sqlite_master 
    WHERE type = 'table'""").fetchall()
    return [x[0] for x in tables]
    
def getChat(cursor,cid,delta = None):
    dstr = ""
    if delta != None:
        dstr = """ WHERE date >  date('now','-"""+str(delta+1)+""" day') """
    lines =  cursor.execute("""
    SELECT userid 
    FROM """+cid+dstr).fetchall()
    return [line[0] for line in lines]
        
def getUsersID(cursor,cid):
    idlist = cursor.execute("""
        SELECT DISTINCT userid 
        FROM """+cid).fetchall()
    return [i[0] for i in idlist]

def getUsersName(cursor,cid,uid):
    idlist = cursor.execute("""
        SELECT DISTINCT username 
        FROM %s 
        WHERE userid=%s""",(cid,uid,)).fetchall()
    return idlist[0].replace("$","S")

def getUsers(cursor,cid):
    tlist = cursor.execute("""
        SELECT DISTINCT userid,username 
        FROM """+cid+"""
        GROUP BY userid""").fetchall()
    ids = [x[0] for x in tlist]
    names = [x[1].replace("$","S") for x in tlist]
    return ids, names
            
def getMsgsPerUser(cursor,cid,delta = None):
    dstr = ""
    if delta != None:
        dstr = """ WHERE date >  date('now','-"""+str(delta+1)+""" day') """
    tuples = cursor.execute("""
        SELECT userid, COUNT(*) as count, username 
        FROM """+cid+dstr+""" 
        GROUP BY userid 
        ORDER BY count DESC""").fetchall()
    out = []
    for t in tuples:
        t = list(t)
        t[2] = t[2].replace("$","S")
        out.append(t)
    return out

def getCharsPerUser(cursor,cid,delta = None):
    dstr = ""
    if delta != None:
        dstr = """ WHERE date >  date('now','-"""+str(delta+1)+""" day') """
    tuples = cursor.execute("""
        SELECT userid,SUM(length) as length, username 
        FROM """+cid+dstr+""" 
        GROUP BY userid 
        ORDER BY length DESC""").fetchall()
    out = []
    for t in tuples:
        t = list(t)
        t[2] = t[2].replace("$","S")
        out.append(t)
    return out

def findChats(cursor,uid):
    tables = getManagedChats(cursor)
    tables = [x for x in tables if x.find("G")!=-1]
    out = []
    for table in tables:
        users = getUsersID(cursor,table)
        if uid in users:
            out.append(table)
    return list(set(out))

def lookFor(lines,u1,u2):
    output = 0
    for i in range(0,len(lines)-1):
        if lines[i] == u1 and lines[i+1] == u2:
            output += 1
    return output

#         WHERE date > """+str(t.year)+"-"+str(t.month)+"-"+str(t.day)+"""
def getMessageVectorPerUser(cursor,cid,uid,delta = None):
    dstr = " WHERE "
    if delta != None:
        dstr = """ WHERE date >  date('now','-"""+str(delta+1)+""" day') AND """
    clist = cursor.execute("""SELECT strftime('%Y-%m-%d',date),count(userid) 
        FROM """+cid+dstr+""" 
        userid="""+str(uid)+"""
        GROUP BY strftime('%d', date)""").fetchall()
    return clist

def getPie(cursor,cid):
    mlist = getMsgsPerUser(cursor,cid)
    colorm = [colors[i%len(colors)] for i in range(0, len(mlist))]
    labels = [item[2] for item in mlist]
    sizes = [item[1] for item in mlist]
    fig1, ax1 = plt.subplots()
    ax1.pie(sizes, labels=labels, colors = colorm, autopct='%1.1f%%')
    ax1.axis('equal')  # Equal aspect ratio ensures that pie is drawn as a circle.
    #plt.tight_layout()
    plt.savefig("/tmp/"+cid+"_pie.jpg", format='jpg')
    plt.close()

def getLengthPerUser(cursor,cid):
    data = []
    names = []
    scatterdata = []
    ids,realnames = getUsers(cursor,cid)
    mlist = getMsgsPerUser(cursor,cid)
    clist = getCharsPerUser(cursor,cid)
    for i in range(0,len(ids)):
        for mitem in mlist: 
            if(mitem[0]== ids[i]):
                for citem in clist:
                    if(citem[0]== ids[i]):
                        diff = (1.0*citem[1])/(1.0*mitem[1])
                        data.append(diff)
                        scatterdata.append((mitem[1], diff, realnames[i]))
                        names.append(realnames[i])

    fig, ax = plt.subplots()
    for item in scatterdata:
        ax.scatter(item[0],item[1])
        ax.annotate(item[2], (item[0],item[1]))
    ax.set_xlabel('Messages')
    ax.set_ylabel('Message length')
    plt.savefig("/tmp/"+cid+"_length.jpg", format='jpg')
    plt.close()

def getConnections(cursor,cid):
    users, realnames = getUsers(cursor,cid)
    lines =  cursor.execute("SELECT userid FROM "+cid).fetchall()
    lines = [line[0] for line in lines]

    connections = np.zeros((len(users), len(users)))
    for u1 in range(0,len(users)):
        for u2 in range(0,len(users)):
            n = lookFor(lines,users[u1],users[u2])
            connections[u1][u2] = n

    maxc = np.matrix(connections).max()
    meanc = 100*np.matrix(connections).mean()/maxc
    stdc = 100*np.matrix(connections).std()/maxc
    for u1 in range(0,len(users)):
        for u2 in range(0,len(users)):
            connections[u1][u2] = 100.0*connections[u1][u2]/maxc
            if connections[u1][u2] < meanc or u1==u2:
                connections[u1][u2] = 0
    sb.set()
    pdata = pd.DataFrame(data=connections, index=realnames,columns=realnames)  # 1st row as the column names
    ax = sb.heatmap(pdata, annot=True, linewidths=.2)
    ax.set(xlabel='Origen', ylabel='Destino')
    plt.savefig("/tmp/"+cid+"_conn.jpg", format='jpg')
    plt.close()

def getNet(cursor,cid):
    users, realnames = getUsers(cursor,cid)
    lines =  cursor.execute("SELECT userid FROM "+cid).fetchall()
    lines = [line[0] for line in lines]
    G = nx.Graph()
    connections = np.zeros((len(users), len(users)))
    for u1 in range(0,len(users)):
        for u2 in range(0,len(users)):
            n = lookFor(lines,users[u1],users[u2])
            connections[u1][u2] = n

    maxc = np.matrix(connections).max()
    meanc = 100*np.matrix(connections).mean()/maxc
    stdc = 100*np.matrix(connections).std()/maxc
    for u1 in range(0,len(users)):
        G.add_node(realnames[u1])
        for u2 in range(0,len(users)):
            w = 100.0*connections[u1][u2]/maxc
            if w > (meanc/1.0) and u1!=u2:
                G.add_edge(realnames[u1], realnames[u2], weight=w)
    plt.subplot(111)
    edges,weights = zip(*nx.get_edge_attributes(G,'weight').items())
    pos = nx.circular_layout(G)
    nx.draw(G,
        pos,
        node_color = 'b',
        font_size=9,
        edgelist = edges,
        edge_color = weights,
        width = 5.0,
        edge_cmap = plt.cm.Blues,
        with_labels = False)
    offset = 0.12
    for p in pos:  # raise text positions
        pos[p][1] += offset
    nx.draw_networkx_labels(G, pos)
    plt.autoscale()
    #plt.tight_layout(pad=2, w_pad=0.5, h_pad=1.0)
    plt.tight_layout()
    plt.savefig("/tmp/"+cid+"_net.jpg", format='jpg')
    plt.close()    

def bash_history_map(lines):
    otlist = []
    for strdate in lines:
        strdate = strdate.replace("+00:00","")
        try:
            tm = time.strptime(strdate, '%Y-%m-%d %H:%M:%S')
        except Exception as e:
            logging.error("{} parsing {}".format(e, strdate))
        otlist.append((tm.tm_wday, tm.tm_hour))
    return otlist

def time_count_reduce(time_list):
    h = defaultdict(lambda: 0)
    # return sparse dict of  key - (day, hour) tuples  value - count of tuples
    for k, g in groupby(sorted(time_list)):
        h[k] += sum(1 for _ in g)
    return h

def genWeeklyFreq(cursor,cid,delta = None):
    title = None
    cstr = ""
    title = "Frequency since the beggining."
    if delta != None:
        cstr = """ WHERE date >  date('now','-"""+str(delta+1)+""" day') """
        title = "Frequency last "+str(delta)+" days."
    lines =  cursor.execute("""
        SELECT date 
        FROM """+cid+cstr).fetchall()
    lines = [line[0] for line in lines]
    logging.info(lines[0])
    logging.info(lines[-1])
    h = time_count_reduce(bash_history_map(lines))

    data = [[h[x,y] for y in range(24)] for x in range(7)]

    maxvalue = max(max(i) for i in data)
    if maxvalue==0:
        logging.warning("No traffic on this chat")
        return
    maxrender = 325 # seems good enough
    xs, ys, rs, ss = [], [], [], []
    for y, d in enumerate(data):
        for x, n in enumerate(d):
            xs.append(x)
            ys.append(y)
            linear_scale = float(n)/float(maxvalue) * maxrender
            log_scale = linear_scale**2/maxrender
            ss.append(log_scale)
    # create a figure an axes with the same background color
    # facecolor -> axisbg
    #fig = plt.figure(figsize=(8, title and 3 or 2.5), axisbg='#efefef')
    fig = plt.figure(figsize=(8, title and 3 or 2.5))
    ax = fig.add_subplot('111',facecolor='#efefef')
    ax.set_facecolor('#efefef')
    # make the figure margins smaller
    if title:
        fig.subplots_adjust(left=0.06, bottom=0.04, right=0.98, top=0.95)
        ax.set_title(title, y=0.96).set_color('#333333')
    else:
        fig.subplots_adjust(left=0.06, bottom=0.10, right=0.98, top=0.99)
    # don't display the axes frame
    ax.set_frame_on(False)
    # plot the punch card data
    ax.scatter(xs, ys[::-1], s=ss, c='#333333', edgecolor='#333333')
    # hide the tick lines
    for line in ax.get_xticklines() + ax.get_yticklines():
        line.set_alpha(0.0)
    # draw x and y lines (instead of axes frame)
    dist = -0.8
    ax.plot([dist, 24], [dist, dist], c='#999999')
    ax.plot([dist, dist], [dist, 6.4], c='#999999')
    # select new axis limits
    ax.set_xlim(-0.9, 24.5)
    ax.set_ylim(-0.9, 7)
    # set tick labels and draw them smaller than normal
    ax.set_yticks(range(7))
    for tx in ax.set_yticklabels(days[::-1]):
        tx.set_color('#555555')
        tx.set_size(8)
    ax.set_xticks(range(24))
    t = '12am|1|2|3|4|5|6|7|8|9|10|11|12pm|1|2|3|4|5|6|7|8|9|10|11'.split('|')
    for tx in ax.set_xticklabels(t):
        tx.set_color('#555555')
        tx.set_size(8.5)
    # get equal spacing for days and hours
    ax.set_aspect('equal')
    # Produce an image.
    if delta==None:
        plt.savefig("/tmp/"+cid+"_freq.jpg", format='jpg')
    else:
        plt.savefig("/tmp/"+cid+"_freq_delta"+str(delta)+".jpg", format='jpg')
        
    plt.clf()
    plt.close('all')

def lookForDate(vect,d):
    for i in range(0,len(vect)):
        if str(vect[i][0]) == str(d):
            return vect[i][1]
    return 0

def getMessagesOverTime(cursor,cid,md):
    fig, ax = plt.subplots()
    vects = []

    now = datetime.now().date()
    dates = [now - timedelta(days=x) for x in range(0, md)]
    dates = list(reversed(dates))

    users, realnames = getUsers(cursor,cid)
    for x in range(0,len(users)):
        vect = getMessageVectorPerUser(cursor,cid, users[x],md)
        hits = []
        for d in dates:
            h = lookForDate(vect,d)
            hits.append(h)
        vects.append(hits)

    plt.stackplot(dates,vects,colors=colors[0:len(realnames)])
    print(realnames)
    ax.set(xlabel='Date',ylabel='Messages')
    plt.legend([mpatches.Patch(color=colors[x%len(colors)]) for x in range(0,len(realnames))],realnames)
    ax.legend(loc=2)

    ax.xaxis.set_major_locator(matplotlib.dates.DayLocator(interval=1))  
    ax.xaxis.set_major_formatter(DateFormatter('%d\n%b'))  

    plt.savefig("/tmp/"+cid+"_tl.jpg", format='jpg')
    plt.close()


def genPoster(cursor,cid,timeDelta = None):
    logging.info("Started poster generation")

    f, ax = plt.subplots(4, sharex=False,figsize=(6,16))
    if(timeDelta!=None):
        plt.suptitle('Last '+str(timeDelta)+' days stats',fontsize=25)
    else:
        plt.suptitle('Total stats',fontsize=25)
    mlist = getMsgsPerUser(cursor,cid,timeDelta)
    clist = getCharsPerUser(cursor,cid,timeDelta)
    ids,realnames = getUsers(cursor,cid)

    if len(ids) > 50:
        logging.warning("Group is too big: "+ str(len(ids)))
        return
    if len(ids) < 2:
        logging.warning("Group is too small: "+str(len(ids)))
        return

    colorm = [colors[i%len(colors)] for i in range(0, len(mlist))]

    labels = [item[2] for item in mlist]
    sizes = [item[1] for item in mlist]
    ax[0].pie(sizes, labels=labels, colors = colorm, autopct='%1.1f%%')
    ax[0].axis('equal')  # Equal aspect ratio ensures that pie is drawn as a circle.
    ax[0].set_title('Msgs per user', size=20) # Title

    data = []
    names = []
    scatterdata = []
    for i in range(0,len(ids)):
        for mitem in mlist:
            if(mitem[0]== ids[i]):
                for citem in clist:
                    if(citem[0]== ids[i]):
                        diff = (1.0*citem[1])/(1.0*mitem[1])
                        data.append(diff)
                        scatterdata.append((mitem[1], diff, realnames[i]))
                        names.append(realnames[i])

    for item in scatterdata:
        ax[1].scatter(item[0],item[1])
        ax[1].annotate(item[2], (item[0],item[1]))
    ax[1].set_title('Msgs/Characters per user', size=20,pad=10) # Title
    ax[1].set_xlabel('Messages')
    ax[1].set_ylabel('Characters')

    ax[2].set_title('Friendship', size=20, y=1.08) # Title
    lines = getChat(cursor,cid,timeDelta)
    logging.info("    Chat lines obtained: "+str(len(lines)))
    if len(lines) < 25:
        logging.warning("Not enough traffic: "+str(len(lines)))
        return
    G = nx.Graph()
    connections = np.zeros((len(ids), len(ids)))
    for u1 in range(0,len(ids)):
        for u2 in range(0,len(ids)):
            n = lookFor(lines,ids[u1],ids[u2])
            connections[u1][u2] = n

    maxc = np.matrix(connections).max()
    meanc = 100*np.matrix(connections).mean()/maxc
    stdc = 100*np.matrix(connections).std()/maxc
    for u1 in range(0,len(ids)):
        G.add_node(realnames[u1])
        for u2 in range(0,len(ids)):
            w = 100.0*connections[u1][u2]/maxc
            if w > (meanc/1.0) and u1!=u2:
                G.add_edge(realnames[u1], realnames[u2], weight=w)

    try:
        edges, weights = zip(*nx.get_edge_attributes(G,'weight').items())
    except Exception as e:
        logging.warning("Not significant connections between members")
        return
    pos = nx.circular_layout(G)
    # logging.info(pos)
    # logging.info(edges)
    logging.info(G)
    nx.draw(G,
        pos,
        node_color = colorm,
        ax=ax[2],
        font_size=7,
        edgelist = edges,
        edge_color = weights,
        width = 5.0,
        edge_cmap = plt.cm.Blues,
        with_labels = False
    )
    offset = 0.12
    for p in pos:  # raise text positions
        pos[p][1] += offset
    nx.draw_networkx_labels(G, pos,ax=ax[2])

    #box = ax[2].get_position()
    #ax[2].set_position([box.x0, box.y0, box.width, box.height * 1.1])

    vects = []
    now = datetime.now().date()

    if timeDelta != None :
        dates = [now - timedelta(days=x) for x in range(0, timeDelta)]
        dates = list(reversed(dates))
    else:
        d1 = getStartDate(cursor,cid)
        totald = now - d1
        dates = [(d1 + timedelta(days=i)) for i in range(0,totald.days)]
        #dates = list(reversed(dates))

    for x in range(0,len(ids)):
        vect = getMessageVectorPerUser(cursor,cid, ids[x],timeDelta)
        hits = []
        for d in dates:
            h = lookForDate(vect,d)
            hits.append(h)
        vects.append(hits)

    box = ax[3].get_position()
    ax[3].stackplot(dates,vects,colors=colors[0:len(realnames)])
    ax[3].set(xlabel='Date',ylabel='Messages')
    handles = [mpatches.Patch(color=colors[x%len(colors)]) for x in range(0,len(realnames))]

    ax[3].legend(handles,realnames,loc='upper center', bbox_to_anchor=(0.3, 1.25),ncol=2)
    #ax[3].legend(loc=2)

    ax[3].xaxis.set_major_locator(matplotlib.dates.DayLocator(interval=1))
    ax[3].xaxis.set_major_formatter(DateFormatter('%d\n%b'))
    if timeDelta==None or timeDelta > 30:
        ax[3].xaxis.set_major_locator(matplotlib.dates.MonthLocator(interval=1))
        ax[3].xaxis.set_major_formatter(DateFormatter('%b'))

    ax[3].set_position([box.x0, box.y0, box.width, box.height])

    # pad=1.0, w_pad=1.0, h_pad=1.0 rect=[0.0, 0, 0.95, 0.95]
    plt.tight_layout(pad=1.0, w_pad=1.0, h_pad=1.0, rect=[0.0, 0, 0.95, 0.95])
    #f.subplots_adjust(bottom = 0)
    #f.subplots_adjust(top = 0.95)
    #f.subplots_adjust(right = 0.95)
    #f.subplots_adjust(left = 0)
    if timeDelta==None :
        plt.savefig("/tmp/"+cid+"_all.jpg", format='jpg')
    else:
        plt.savefig("/tmp/"+cid+"_delta"+str(timeDelta)+".jpg", format='jpg')
    plt.clf()
    plt.close('all')

    logging.debug("  Finished Poster generation")

def generateStats():
    conn = sqlite3.connect('tgdatabase.db')
    c = conn.cursor()
    deleteUnused(c)
    conn.commit()
    sleep(1)
    chats = getManagedChats(c)
    chats = [x for x in chats if x.find("G")!=-1]
    i=0
    for cid in chats:
        #try:
        logging.info("["+str(i)+"/"+str(len(chats))+"] Started Async Stats Gen:\t"+cid)
        genWeeklyFreq(c,cid,120)
        genWeeklyFreq(c,cid,15)
        genPoster(c,cid,15)
        genPoster(c,cid,120)
        logging.info("["+str(i)+"/"+str(len(chats))+"] Ended Async Stats Gen:\t"+cid)
        #except Exception as e:
        #    logging.exception("["+str(i)+"/"+str(len(chats))+"] ["+cid+"] : "+str(e))
        i+=1
    conn.close()

logging.basicConfig(
    format='%(asctime)s - %(levelname)s - %(message)s',
    level=logging.INFO,
    filename='/tmp/cotibot.log'
    )
coloredlogs.DEFAULT_LOG_FORMAT = '%(asctime)s - %(levelname)s - %(message)s'
coloredlogs.install()

signal.signal(signal.SIGINT, signal_handler)
conn = sqlite3.connect('tgdatabase.db')
cursor = conn.cursor()
deleteUnused(cursor)
conn.commit()
sleep(1)
chats = getManagedChats(cursor)
chats = [x for x in chats if x.find("G")!=-1]
i=0
for chatID in chats:
    try:
        logging.info("[{0}/{1}] Started Async Stats Gen for chat {2}".format(i,len(chats),chatID))
        genWeeklyFreq(cursor, chatID, 120)
        genWeeklyFreq(cursor, chatID, 15)
        try:
            genPoster(cursor, chatID, 15)
        except Exception as e:
            logging.exception("["+str(i)+"/"+str(len(chats))+"] ["+chatID+"] : "+str(e))
        try:
            genPoster(cursor, chatID, 120)
        except Exception as e:
            logging.exception("["+str(i)+"/"+str(len(chats))+"] ["+chatID+"] : "+str(e))
        logging.info("[{0}/{1}] Ended Async Stats Gen for chat {2}".format(i,len(chats),chatID))
    except Exception as e:
        logging.exception("["+str(i)+"/"+str(len(chats))+"] ["+chatID+"] : "+str(e))
    i+=1
conn.close()
gc.collect()
