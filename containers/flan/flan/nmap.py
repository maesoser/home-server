import sys
import json
import os
import random, string
from datetime import date
import subprocess

results = {}
vulnerable_services = []

def do_scan(targets):
    today = date.today()
    scanID = "{0}_{1}_{2}".format(
        ''.join(random.choice(string.ascii_uppercase + string.ascii_lowercase + string.digits) for _ in range(8)),
        today.strftime("%d%m%Y")
        today.strftime("%H:%M:%S")
        )
    print("Starting scan {0}".format(scanID))
    for target in targets:
        print("[{0}] Scanning {1}".format(scanID, target))
        do_nmap(scanID, target)

# Un NMAP que se inicia tiene 3 variables:
#   -targets
#   - nmap commands
#   - start_date

def do_nmap(scanID, target):
    # nmap -sV -oX "${filepath}.xml" -oN - -v1 $@ --script=vulners/vulners.nse $line > "${filepath}.txt" 2>&1
    nmap_command = [
        "nmap",
        "-sV",
        "-oX /reports/{0}/{1}.xml".format(scanID, target),
        "-oN",
        "-v1 $@",
        "--script=vulners/vulners.nse",
        "{0}".format(target)
    ]
    print("Executing {0}".format(' '.join(nmap_command)))
    process = subprocess.Popen(
        nmap_command,
        stdout=subprocess.PIPE,
        universal_newlines=True
    )
    raw_output = ""
    return_code = 0
    while True:
        output = process.stdout.readline()
        raw_output += output
        print(output.strip())
        return_code = process.poll()
        if return_code is not None:
            print('RETURN CODE', return_code)
            for output in process.stdout.readlines():
                raw_output += output
                print(output.strip())
            break
    out_file = open("/reports/{0}/{1}.txt".format(scanID, target), "w+")
    out_file.write(raw_output)
    out_file.close()
    return return_code
