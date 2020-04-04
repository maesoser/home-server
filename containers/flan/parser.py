import sys
import json
import urllib.request as urllib
import os
import xmltodict
import argparse
import markdown
import nmap_parser as np

def create_json (nmap_command, start_date, output_file, ip_file, nmap_raw):
    json_obj = {
        "nmap_cmd" : nmap_command,
        "datetime" : start_date,
        "output" : output_file,
        "nmap_raw" :nmap_raw,
        "vulnerable_services" : [],
        "non_vulnerable_services" : [],
        "targets" : []
    }
    if vulnerable_services:
        for s in vulnerable_services:
            vulns = results[s]['vulns']
            locations = results[s]['locations']
            service = {
                "name" : s,
                "endpoints" : [],
                "vulns" : []
            }
            for i, v in enumerate(vulns):
                vuln_obj = {
                    "name" : v['name'],
                    "url" : "https://nvd.nist.gov/vuln/detail/" + v['name'],
                    "severity" : np.convert_severity(v['severity']),
                    "score" : v['severity'],
                    "description" : np.get_description(v['name'], v['type'])
                }
                service["vulns"].append(vuln_obj)
            for addr in locations.keys():
                for port in locations[addr]:
                    service["endpoints"].append('{0}:{1}'.format(addr,port))
            json_obj["vulnerable_services"].append(service)

    non_vuln_services = list(set(results.keys()) - set(vulnerable_services))
    if non_vuln_services:
        for ns in non_vuln_services:
            service = {
                "name" : ns,
                "endpoints" : [],
            }
            for addr in locations.keys():
                for port in locations[addr]:
                    service["endpoints"].append('{0}:{1}'.format(addr,port))
            json_obj["non_vulnerable_services"].append(service)

    f = open(ip_file)
    for line in f:
        if line != "":
            json_obj["targets"].append(line.replace("\n",""))
    f.close()

    return json_obj

def create_latex(nmap_command, start_date, output_file, ip_file):
    f = open('./latex_header.tex')
    write_buffer = f.read()
    f.close()

    write_buffer += "Flan Scan ran a network vulnerability scan with the following Nmap command on " \
                 + start_date \
                 + "UTC.\n\\begin{lstlisting}\n" \
                 + nmap_command \
                 + "\n\end{lstlisting}\nTo find out what IPs were scanned see the end of this report.\n"
    write_buffer += "\section*{Services with Vulnerabilities}"
    if vulnerable_services:
        write_buffer += """\\begin{enumerate}[wide, labelwidth=!, labelindent=0pt,
                        label=\\textbf{\large \\arabic{enumi} \large}]\n"""
        for s in vulnerable_services:
            write_buffer += '\item \\textbf{\large ' + s + ' \large}'
            vulns = results[s]['vulns']
            locations = results[s]['locations']
            num_vulns = len(vulns)

            for i, v in enumerate(vulns):
                write_buffer += '\\begin{figure}[h!]\n'
                severity_name = convert_severity(v['severity'])
                write_buffer += '\\begin{tabular}{|p{16cm}|}\\rowcolor[HTML]{' \
                         + colors[severity_name] \
                         + """} \\begin{tabular}{@{}p{15cm}>{\\raggedleft\\arraybackslash}
                           p{0.5cm}@{}}\\textbf{""" \
                         + v['name'] + ' ' + severity_name + ' (' \
                         + str(v['severity']) \
                         + ')} & \href{https://nvd.nist.gov/vuln/detail/' \
                         + v['name'] + '}{\large \\faicon{link}}' \
                         + '\end{tabular}\\\\\n Summary:' \
                         + get_description(v['name'], v['type']) \
                         + '\\\\ \hline \end{tabular}  '
                write_buffer += '\end{figure}\n'

            write_buffer += '\FloatBarrier\n\\textbf{The above ' \
                         + str(num_vulns) \
                         + """ vulnerabilities apply to these network locations:}\n
                         \\begin{itemize}\n"""
            for addr in locations.keys():
                write_buffer += '\item ' + addr + ' Ports: ' + str(locations[addr])+ '\n'
            write_buffer += '\\\\ \\\\ \n \end{itemize}\n'
        write_buffer += '\end{enumerate}\n'

    non_vuln_services = list(set(results.keys()) - set(vulnerable_services))
    write_buffer += '\section*{Services With No Known Vulnerabilities}'

    if non_vuln_services:
        write_buffer += """\\begin{enumerate}[wide, labelwidth=!, labelindent=0pt,
        label=\\textbf{\large \\arabic{enumi} \large}]\n"""
        for ns in non_vuln_services:
            write_buffer += '\item \\textbf{\large ' + ns \
                            + ' \large}\n\\begin{itemize}\n'
            locations = results[ns]['locations']
            for addr in locations.keys():
                write_buffer += '\item ' + addr + ' Ports: ' + str(locations[addr])+ '\n'
            write_buffer += '\end{itemize}\n'
        write_buffer += '\end{enumerate}\n'

    write_buffer += '\section*{List of IPs Scanned}'
    write_buffer += '\\begin{itemize}\n'
    f = open(ip_file)
    for line in f:
        write_buffer += '\item ' + line + '\n'
    f.close()
    write_buffer += '\end{itemize}\n'

    write_buffer += '\end{document}'

    return write_buffer

def save_file(write_buffer, filename):
    out_file = open(filename, "w+")
    out_file.write(write_buffer)
    out_file.close()

def create_markdown(nmap_command, start_date, output_file, ip_file, nmap_raw):
    write_buffer = ""

    write_buffer += "Flan Scan ran a network vulnerability scan with the following Nmap command on " \
                 + start_date \
                 + " UTC.\n```\n" \
                 + nmap_command \
                 + "\n```\nTo find out what IPs were scanned see the end of this report.\n"

    write_buffer += "## Services with Vulnerabilities\n"
    if vulnerable_services:
        for s in vulnerable_services:
            write_buffer += '### {0} \n'.format(s)
            vulns = results[s]['vulns']
            locations = results[s]['locations']
            num_vulns = len(vulns)

            for i, v in enumerate(vulns):
                severity_name = convert_severity(v['severity'])
                vuln_url = "https://nvd.nist.gov/vuln/detail/" + v['name']
                write_buffer += '\n | [{0}]({1}) **{2}** ({3}) | \n'.format(
			v['name'], vuln_url, severity_name, v['severity'])
                write_buffer += '|---|\n'
                write_buffer += '| {0} | \n'.format(get_description(v['name'], v['type']))

            write_buffer += '\nThe above {0} vulnerabilities apply to these network locations:\n'.format(num_vulns)
            for addr in locations.keys():
                for port in locations[addr]:
                    write_buffer += ' 1. `{0}:{1}`\n'.format(addr,port)

    non_vuln_services = list(set(results.keys()) - set(vulnerable_services))
    write_buffer += '\n## Services With No Known Vulnerabilities\n'
    if non_vuln_services:
        for ns in non_vuln_services:
            write_buffer += '### {0}\n'.format(ns)
            locations = results[ns]['locations']
            for addr in locations.keys():
                for port in locations[addr]:
                    write_buffer += ' - `{0}:{1}`\n'.format(addr, port)

    write_buffer += '\n### List of IPs Scanned\n'
    f = open(ip_file)
    for line in f:
        if line != "":
            write_buffer += ' - {0}\n'.format(line.replace("\n",""))
    f.close()

    write_buffer += '\n### Nmap output\n'
    write_buffer += '```\n'
    write_buffer += nmap_raw
    write_buffer += '\n```\n'
    return write_buffer

def create_html(text):
   write_buffer = """<!DOCTYPE html>
   <html>
     <head>
       <meta charset="utf-8">
       <link href="default.css" type="text/css" rel="stylesheet" />
       <title>FlanScan</title>
     </head>
   <body>
   """
   write_buffer += markdown.markdown(text, extensions=['tables', 'fenced_code'])
   write_buffer += """</body></html>"""
   return write_buffer
