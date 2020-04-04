#!/usr/bin/env bash
rm tgdatabase.db
for i in *.log; do ./txt2db.py "$i"; done
du -h tgdatabase.db
