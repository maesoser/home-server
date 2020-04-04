import sys
import json
import urllib.request as urllib
import os
import xmltodict
import argparse
import markdown

def parse_vuln(ip_addr, port, app_name, vuln):
    vuln_name = ''
    severity = ''
    type = ''
    for field in vuln:
        if field['@key'] == 'cvss':
            severity = float(field['#text'])
        elif field['@key'] == 'id':
            vuln_name = field['#text']
        elif field['@key'] == 'type':
            type = field['#text']
    if 'vulns'in results[app_name].keys():
        results[app_name]['vulns'].append({'name': vuln_name,
                                           'type': type,
                                           'severity': severity})
    else:
        results[app_name]['vulns'] = [{'name': vuln_name,
                                       'type': type,
                                       'severity': severity}]


def parse_script(ip_addr, port, app_name, script):
    if 'table' in script.keys():
        vulnerable_services.append(app_name)
        script_table = script['table']['table']
        if isinstance(script_table, list):
            for vuln in script_table:
                parse_vuln(ip_addr, port, app_name, vuln['elem'])
        else:
            parse_vuln(ip_addr, port, app_name, script_table['elem'])
    else:
        print('ERROR in script: ' + script['@output'] + " at location: " + ip_addr + " port: " + port + " app: " + app_name)


def get_app_name(service):
    app_name = ''
    if '@product' in service.keys():
        app_name += service['@product'] + " "
        if '@version' in service.keys():
            app_name += service['@version'] + " "
    elif '@name' in service.keys():
        app_name += service['@name'] + " "

    if('cpe' in service.keys()):
        if isinstance(service['cpe'], list):
            for cpe in service['cpe']:
                app_name += '(' + cpe + ") "
        else:
            app_name += '(' + service['cpe'] + ") "
    return app_name


def parse_port(ip_addr, port):
    if port['state']['@state'] == 'closed':
        return
    app_name = get_app_name(port['service'])

    port_num = port['@portid']

    if app_name in results.keys():
        if ip_addr in results[app_name]['locations'].keys():
            results[app_name]['locations'][ip_addr].append(port_num)
        else:
            results[app_name]['locations'][ip_addr] = [port_num]
    else:
        results[app_name] = {'locations': {ip_addr: [port_num]}}
        if 'script' in port.keys():
            scripts = port['script']
            if isinstance(scripts, list):
                for s in scripts:
                    if s['@id'] == 'vulners':
                        parse_script(ip_addr, port_num, app_name, s)
            else:
                if scripts['@id'] == 'vulners':
                    parse_script(ip_addr, port_num, app_name, scripts)


def parse_host(host):
    addresses = host['address']
    if isinstance(addresses, list):
        for addr in addresses:
            if "ip" in addr['@addrtype']:
                ip_addr = addr['@addr']
    else:
        ip_addr = addresses['@addr']

    if host['status']['@state'] == 'up' and 'port' in host['ports'].keys():
        ports = host['ports']['port']
        if isinstance(ports, list):
            for p in ports:
                parse_port(ip_addr, p)
        else:
            parse_port(ip_addr, ports)


def parse_results(data):
    if 'host' in data['nmaprun'].keys():
        hosts = data['nmaprun']['host']

        if isinstance(hosts, list):
            for h in hosts:
                parse_host(h)
        else:
            parse_host(hosts)


def convert_severity(sev):
    if sev < 4:
        return 'Low'
    elif sev < 7:
        return 'Medium'
    else:
        return 'High'


def get_description(vuln, type):
    if type == 'cve':
        year = vuln[4:8]
        section = vuln[9:-3] + 'xxx'
        url = """https://raw.githubusercontent.com/CVEProject/cvelist/master/{}/{}/{}.json""".format(year, section, vuln)
        cve_json = json.loads(urllib.urlopen(url).read().decode("utf-8"))
        return cve_json["description"]["description_data"][0]["value"]
    else:
        return ''

def parse_nmap_command(raw_command):
    nmap_split = raw_command.split()[:-1] #remove last element, ip address
    nmap_split[3] = "<output-file>"
    return " ".join(nmap_split)
