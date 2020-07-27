#! /usr/bin/env python3

import argparse, json, random
import dns.resolver
from datetime import datetime
import multiprocessing.dummy as mp

def dig(target, nameserver):
    job = {
        "success": False,
        "nameserver": nameserver["domain"],
        "response": [],
        "error" : ""
    }
    dnsresolver = dns.resolver.Resolver()
    dnsresolver.nameservers = [ nameserver["addr"] ]
    dnsresolver.timeout = 3
    dnsresolver.lifetime = 3
    try:
        result = dnsresolver.query(target, 'a')
    except dns.resolver.NXDOMAIN:
        job["error"] = "Error: No such domain"
        return job
    except dns.resolver.Timeout:
        job["error"] = "Error: Timed out"
        return job
    except dns.exception.DNSException as e:
        job["error"] = "Error: {}".format(e)
        return job
    job["response"] = [x.to_text() for x in result]
    job["response"].sort()
    job["success"] = True
    if job["response"] is None:
        job["success"] = False
    return job

def get_nameservers(country, filename):
    out = []
    netf = open(filename, 'r')
    while True:
        line = netf.readline()
        if not line:
            break
        obj = json.loads(line)
        if obj["country_iso"] == country:
            out.append(obj)
    netf.close()
    return out

def get_countries(filename):
    out = []
    netf = open(filename, 'r')
    while True:
        line = netf.readline()
        if not line:
            break
        obj = json.loads(line)
        out.append("{} ({})".format(obj["country_iso"],obj["country_name"]))
    netf.close()
    out = list(set(out))
    return out

def process_net(domain_entry):
    job = dig(args.target, domain_entry)
    domain_entry["job"] = job
    return domain_entry

def get_args():
    parser = argparse.ArgumentParser(description="Performs DNS resolution on Namservers that belongs to an specific country")
    parser.add_argument('--country', required=False, help='Country (Empty lists available countries)', default="")
    parser.add_argument('--target', required=False, help='DNS target', default="wikipedia.org")
    parser.add_argument('--jobs', required=False, help='Number of concurrent queries', default=64)
    parser.add_argument('--verbose', dest='verbose', help='Diplays verbose output', action='store_true')
    parser.add_argument('--json', dest='jsonout', help='Display JSON output', action='store_true')
    parser.add_argument('--db', required=False, help='Database file', default="ripe_domains.json")
    args = parser.parse_args()
    return args

if __name__ == '__main__':

    args = get_args()

    if args.country == "":
        print("Available countries :")
        countries = get_countries(args.db)
        for c in countries:
            print("\t- {0}".format(c))
        exit(0)

    networks = get_nameservers(args.country, args.db)
    print("{} networks found at {}".format(len(networks), args.country))

    p = mp.Pool(args.jobs)
    results = p.map(process_net,networks)
    p.close()
    p.join()

    print("")
    success = 0
    if args.verbose:
        for result in results:
            if result["job"]["success"] == True:
                print("{} - {} \t({})".format(result["as_num"], result["as_org"], result["domain"]))
                for item in result["job"]["response"]:
                    print("   - {0}".format(item))
    elif args.jsonout:
        print(json.dumps(results, indent=2))
    else:
        resultd = {}
        for result in results:
            if result["job"]["success"] == True:
                success += 1
                key = ';'.join(result["job"]["response"])
                if resultd.get(key) == None:
                    resultd[key] = [result]
                else:
                    resultd[key].append(result)
        for key in resultd:
            print("{0} networks returned the following addresses".format(len(resultd[key])))
            for n in resultd[key][0]["job"]["response"]:
                print("\t{}".format(n))
    print("{} networks did not answered correctly".format(len(results) - success))
