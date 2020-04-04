#!/usr/bin/env python3

import sqlite3
import sys

if len(sys.argv)!=2:
    print("Usage: "+sys.argv[0]+" [txt_file]")
    exit()

fname = sys.argv[1]
print("Importing:",fname)

conn = sqlite3.connect('tgdatabase.db')

c = conn.cursor()
# 8997, 2017-12-14 21:35:50,    -1,     383639205,      Behruz_0098,    0,      😐
# message_id message.date message.length message.from_user.id   name reply, emoticono
tablename = fname.replace(".txt","").replace(".log","")
if tablename.find('G')==-1:
     tablename = 'U'+str(tablename)
c.execute('''CREATE TABLE IF NOT EXISTS '''+tablename+
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
n = 0
e = 0
with open(fname, 'r') as infile:
    for line in infile:
        line = line.replace("\n","").split(",")
        if len(line)== 6:
            line.append(" ")
        #print(line)
        try:
          c.executemany("INSERT INTO "+tablename+" VALUES (?,?,?,?,?,?,?)",[line])
        except Exception as err:
             print(err)
             print("on line:")
             print(line)
             e +=1
        n+=1
print(str(n)+" messages imported,"+str(e)+" errors.")
# Save (commit) the changes
conn.commit()
conn.close()
