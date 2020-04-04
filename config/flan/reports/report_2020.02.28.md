Flan Scan ran a network vulnerability scan with the following Nmap command on Fri Feb 28 14:30:20 2020 UTC.
```
nmap -sV -oX <output-file> -oN - -v1 --script=vulners/vulners.nse
```
To find out what IPs were scanned see the end of this report.
## Services with Vulnerabilities
### MiniUPnP 1.9 (cpe:/a:miniupnp_project:miniupnpd:1.9)  

 | [CVE-2017-8798](https://nvd.nist.gov/vuln/detail/CVE-2017-8798) **High** (7.5) | 
|---|
| Integer signedness error in MiniUPnP MiniUPnPc v1.4.20101221 through v2.0 allows remote attackers to cause a denial of service or possibly have unspecified other impact. | 

 | [CVE-2014-3985](https://nvd.nist.gov/vuln/detail/CVE-2014-3985) **Medium** (5.0) | 
|---|
| The getHTTPResponse function in miniwget.c in MiniUPnP 1.9 allows remote attackers to cause a denial of service (crash) via crafted headers that trigger an out-of-bounds read. | 

 | [CVE-2017-1000494](https://nvd.nist.gov/vuln/detail/CVE-2017-1000494) **Medium** (4.6) | 
|---|
| Uninitialized stack variable vulnerability in NameValueParserEndElt (upnpreplyparse.c) in miniupnpd < 2.0 allows an attacker to cause Denial of Service (Segmentation fault and Memory Corruption) or possibly have unspecified other impact | 

The above 3 vulnerabilities apply to these network locations:
 1. `192.168.0.1:5000`
### MySQL 5.5.5-10.1.43-MariaDB-0ubuntu0.18.04.1 (cpe:/a:mysql:mysql:5.5.5-10.1.43-mariadb-0ubuntu0.18.04.1)  

 | [NODEJS:602](https://nvd.nist.gov/vuln/detail/NODEJS:602) **Low** (0.0) | 
|---|
|  | 

The above 1 vulnerabilities apply to these network locations:
 1. `192.168.0.10:3306`
### lighttpd 1.4.45 (cpe:/a:lighttpd:lighttpd:1.4.45)  

 | [CVE-2018-19052](https://nvd.nist.gov/vuln/detail/CVE-2018-19052) **Medium** (5.0) | 
|---|
| An issue was discovered in mod_alias_physical_handler in mod_alias.c in lighttpd before 1.4.50. There is potential ../ path traversal of a single directory above an alias target, with a specific mod_alias configuration where the matched alias lacks a trailing '/' character, but the alias target filesystem path does have a trailing '/' character. | 

The above 1 vulnerabilities apply to these network locations:
 1. `192.168.0.10:8080`
### Boa HTTPd 0.94.14rc21 (cpe:/a:boa:boa:0.94.14rc21)  

 | [CVE-2009-4496](https://nvd.nist.gov/vuln/detail/CVE-2009-4496) **Medium** (5.0) | 
|---|
| Boa 0.94.14rc21 writes data to a log file without sanitizing non-printable characters, which might allow remote attackers to modify a window's title, or possibly execute arbitrary commands or overwrite files, via an HTTP request containing an escape sequence for a terminal emulator. | 

The above 1 vulnerabilities apply to these network locations:
 1. `192.168.0.11:80`

## Services With No Known Vulnerabilities
### blackice-icecap 
 - `192.168.0.1:8081`
### http 
 - `192.168.0.10:80`
### lighttpd (cpe:/a:lighttpd:lighttpd) 
 - `192.168.0.1:80`
 - `192.168.0.1:443`
### jetdirect 
 - `192.168.0.10:9100`
### tcpwrapped 
 - `192.168.0.10:443`
 - `192.168.0.10:5900`
 - `192.168.0.54:62078`
 - `192.168.0.56:62078`
### ppp 
 - `192.168.0.10:3000`
### blackice-alerts 
 - `192.168.0.1:8082`
### dnsmasq pi-hole-2.80 (cpe:/a:thekelleys:dnsmasq:pi-hole-2.80) 
 - `192.168.0.10:53`
### Golang net/http server (cpe:/a:protocol_labs:go-ipfs) 
 - `192.168.0.10:9090`
### Postfix smtpd (cpe:/a:postfix:postfix) 
 - `192.168.0.10:25`
### axhttpd/1.4.9 
 - `192.168.0.11:443`
### Jetty 9.4.8.v20180619 (cpe:/a:mortbay:jetty:9.4.8.v20180619) 
 - `192.168.0.10:4444`
### nginx 1.16.1 (cpe:/a:igor_sysoev:nginx:1.16.1) 
 - `192.168.0.10:9080`

### List of IPs Scanned
 - 192.168.0.1/24

### Nmap output
```
# Nmap 7.80 scan initiated Fri Feb 28 14:30:20 2020 as: nmap -sV -oX /shared/xml_files/2020.02.28/192.168.0.1-24.xml -oN - -v1 --script=vulners/vulners.nse 192.168.0.1/24
Nmap scan report for 192.168.0.0 [host down]
Nmap scan report for 192.168.0.2 [host down]
Nmap scan report for 192.168.0.3 [host down]
Nmap scan report for 192.168.0.4 [host down]
Nmap scan report for 192.168.0.5 [host down]
Nmap scan report for 192.168.0.6 [host down]
Nmap scan report for 192.168.0.7 [host down]
Nmap scan report for 192.168.0.8 [host down]
Nmap scan report for 192.168.0.9 [host down]
Nmap scan report for 192.168.0.12 [host down]
Nmap scan report for 192.168.0.14 [host down]
Nmap scan report for 192.168.0.15 [host down]
Nmap scan report for 192.168.0.16 [host down]
Nmap scan report for 192.168.0.17 [host down]
Nmap scan report for 192.168.0.18 [host down]
Nmap scan report for 192.168.0.19 [host down]
Nmap scan report for 192.168.0.20 [host down]
Nmap scan report for 192.168.0.21 [host down]
Nmap scan report for 192.168.0.22 [host down]
Nmap scan report for 192.168.0.23 [host down]
Nmap scan report for 192.168.0.24 [host down]
Nmap scan report for 192.168.0.25 [host down]
Nmap scan report for 192.168.0.26 [host down]
Nmap scan report for 192.168.0.27 [host down]
Nmap scan report for 192.168.0.28 [host down]
Nmap scan report for 192.168.0.29 [host down]
Nmap scan report for 192.168.0.30 [host down]
Nmap scan report for 192.168.0.31 [host down]
Nmap scan report for 192.168.0.32 [host down]
Nmap scan report for 192.168.0.33 [host down]
Nmap scan report for 192.168.0.34 [host down]
Nmap scan report for 192.168.0.35 [host down]
Nmap scan report for 192.168.0.36 [host down]
Nmap scan report for 192.168.0.37 [host down]
Nmap scan report for 192.168.0.38 [host down]
Nmap scan report for 192.168.0.39 [host down]
Nmap scan report for 192.168.0.40 [host down]
Nmap scan report for 192.168.0.41 [host down]
Nmap scan report for 192.168.0.42 [host down]
Nmap scan report for 192.168.0.43 [host down]
Nmap scan report for 192.168.0.44 [host down]
Nmap scan report for 192.168.0.45 [host down]
Nmap scan report for 192.168.0.46 [host down]
Nmap scan report for 192.168.0.47 [host down]
Nmap scan report for 192.168.0.48 [host down]
Nmap scan report for 192.168.0.49 [host down]
Nmap scan report for 192.168.0.50 [host down]
Nmap scan report for 192.168.0.51 [host down]
Nmap scan report for 192.168.0.52 [host down]
Nmap scan report for 192.168.0.53 [host down]
Nmap scan report for 192.168.0.55 [host down]
Nmap scan report for 192.168.0.57 [host down]
Nmap scan report for 192.168.0.58 [host down]
Nmap scan report for 192.168.0.59 [host down]
Nmap scan report for 192.168.0.60 [host down]
Nmap scan report for 192.168.0.61 [host down]
Nmap scan report for 192.168.0.62 [host down]
Nmap scan report for 192.168.0.63 [host down]
Nmap scan report for 192.168.0.64 [host down]
Nmap scan report for 192.168.0.65 [host down]
Nmap scan report for 192.168.0.66 [host down]
Nmap scan report for 192.168.0.67 [host down]
Nmap scan report for 192.168.0.68 [host down]
Nmap scan report for 192.168.0.69 [host down]
Nmap scan report for 192.168.0.70 [host down]
Nmap scan report for 192.168.0.71 [host down]
Nmap scan report for 192.168.0.72 [host down]
Nmap scan report for 192.168.0.73 [host down]
Nmap scan report for 192.168.0.74 [host down]
Nmap scan report for 192.168.0.75 [host down]
Nmap scan report for 192.168.0.76 [host down]
Nmap scan report for 192.168.0.77 [host down]
Nmap scan report for 192.168.0.78 [host down]
Nmap scan report for 192.168.0.79 [host down]
Nmap scan report for 192.168.0.80 [host down]
Nmap scan report for 192.168.0.81 [host down]
Nmap scan report for 192.168.0.82 [host down]
Nmap scan report for 192.168.0.83 [host down]
Nmap scan report for 192.168.0.84 [host down]
Nmap scan report for 192.168.0.85 [host down]
Nmap scan report for 192.168.0.86 [host down]
Nmap scan report for 192.168.0.87 [host down]
Nmap scan report for 192.168.0.88 [host down]
Nmap scan report for 192.168.0.89 [host down]
Nmap scan report for 192.168.0.90 [host down]
Nmap scan report for 192.168.0.91 [host down]
Nmap scan report for 192.168.0.92 [host down]
Nmap scan report for 192.168.0.93 [host down]
Nmap scan report for 192.168.0.94 [host down]
Nmap scan report for 192.168.0.95 [host down]
Nmap scan report for 192.168.0.96 [host down]
Nmap scan report for 192.168.0.97 [host down]
Nmap scan report for 192.168.0.98 [host down]
Nmap scan report for 192.168.0.99 [host down]
Nmap scan report for 192.168.0.100 [host down]
Nmap scan report for 192.168.0.101 [host down]
Nmap scan report for 192.168.0.102 [host down]
Nmap scan report for 192.168.0.103 [host down]
Nmap scan report for 192.168.0.104 [host down]
Nmap scan report for 192.168.0.105 [host down]
Nmap scan report for 192.168.0.106 [host down]
Nmap scan report for 192.168.0.107 [host down]
Nmap scan report for 192.168.0.108 [host down]
Nmap scan report for 192.168.0.109 [host down]
Nmap scan report for 192.168.0.110 [host down]
Nmap scan report for 192.168.0.111 [host down]
Nmap scan report for 192.168.0.112 [host down]
Nmap scan report for 192.168.0.113 [host down]
Nmap scan report for 192.168.0.114 [host down]
Nmap scan report for 192.168.0.115 [host down]
Nmap scan report for 192.168.0.116 [host down]
Nmap scan report for 192.168.0.117 [host down]
Nmap scan report for 192.168.0.118 [host down]
Nmap scan report for 192.168.0.119 [host down]
Nmap scan report for 192.168.0.120 [host down]
Nmap scan report for 192.168.0.121 [host down]
Nmap scan report for 192.168.0.122 [host down]
Nmap scan report for 192.168.0.123 [host down]
Nmap scan report for 192.168.0.124 [host down]
Nmap scan report for 192.168.0.125 [host down]
Nmap scan report for 192.168.0.126 [host down]
Nmap scan report for 192.168.0.127 [host down]
Nmap scan report for 192.168.0.128 [host down]
Nmap scan report for 192.168.0.129 [host down]
Nmap scan report for 192.168.0.130 [host down]
Nmap scan report for 192.168.0.131 [host down]
Nmap scan report for 192.168.0.132 [host down]
Nmap scan report for 192.168.0.133 [host down]
Nmap scan report for 192.168.0.134 [host down]
Nmap scan report for 192.168.0.135 [host down]
Nmap scan report for 192.168.0.136 [host down]
Nmap scan report for 192.168.0.137 [host down]
Nmap scan report for 192.168.0.138 [host down]
Nmap scan report for 192.168.0.139 [host down]
Nmap scan report for 192.168.0.140 [host down]
Nmap scan report for 192.168.0.141 [host down]
Nmap scan report for 192.168.0.142 [host down]
Nmap scan report for 192.168.0.143 [host down]
Nmap scan report for 192.168.0.144 [host down]
Nmap scan report for 192.168.0.145 [host down]
Nmap scan report for 192.168.0.146 [host down]
Nmap scan report for 192.168.0.147 [host down]
Nmap scan report for 192.168.0.148 [host down]
Nmap scan report for 192.168.0.149 [host down]
Nmap scan report for 192.168.0.150 [host down]
Nmap scan report for 192.168.0.151 [host down]
Nmap scan report for 192.168.0.152 [host down]
Nmap scan report for 192.168.0.153 [host down]
Nmap scan report for 192.168.0.154 [host down]
Nmap scan report for 192.168.0.155 [host down]
Nmap scan report for 192.168.0.156 [host down]
Nmap scan report for 192.168.0.157 [host down]
Nmap scan report for 192.168.0.158 [host down]
Nmap scan report for 192.168.0.159 [host down]
Nmap scan report for 192.168.0.160 [host down]
Nmap scan report for 192.168.0.161 [host down]
Nmap scan report for 192.168.0.162 [host down]
Nmap scan report for 192.168.0.163 [host down]
Nmap scan report for 192.168.0.164 [host down]
Nmap scan report for 192.168.0.165 [host down]
Nmap scan report for 192.168.0.166 [host down]
Nmap scan report for 192.168.0.167 [host down]
Nmap scan report for 192.168.0.168 [host down]
Nmap scan report for 192.168.0.169 [host down]
Nmap scan report for 192.168.0.170 [host down]
Nmap scan report for 192.168.0.171 [host down]
Nmap scan report for 192.168.0.172 [host down]
Nmap scan report for 192.168.0.173 [host down]
Nmap scan report for 192.168.0.174 [host down]
Nmap scan report for 192.168.0.175 [host down]
Nmap scan report for 192.168.0.176 [host down]
Nmap scan report for 192.168.0.177 [host down]
Nmap scan report for 192.168.0.178 [host down]
Nmap scan report for 192.168.0.179 [host down]
Nmap scan report for 192.168.0.180 [host down]
Nmap scan report for 192.168.0.181 [host down]
Nmap scan report for 192.168.0.182 [host down]
Nmap scan report for 192.168.0.183 [host down]
Nmap scan report for 192.168.0.184 [host down]
Nmap scan report for 192.168.0.185 [host down]
Nmap scan report for 192.168.0.186 [host down]
Nmap scan report for 192.168.0.187 [host down]
Nmap scan report for 192.168.0.188 [host down]
Nmap scan report for 192.168.0.189 [host down]
Nmap scan report for 192.168.0.190 [host down]
Nmap scan report for 192.168.0.191 [host down]
Nmap scan report for 192.168.0.192 [host down]
Nmap scan report for 192.168.0.193 [host down]
Nmap scan report for 192.168.0.194 [host down]
Nmap scan report for 192.168.0.195 [host down]
Nmap scan report for 192.168.0.196 [host down]
Nmap scan report for 192.168.0.197 [host down]
Nmap scan report for 192.168.0.198 [host down]
Nmap scan report for 192.168.0.199 [host down]
Nmap scan report for 192.168.0.200 [host down]
Nmap scan report for 192.168.0.201 [host down]
Nmap scan report for 192.168.0.202 [host down]
Nmap scan report for 192.168.0.203 [host down]
Nmap scan report for 192.168.0.204 [host down]
Nmap scan report for 192.168.0.205 [host down]
Nmap scan report for 192.168.0.206 [host down]
Nmap scan report for 192.168.0.207 [host down]
Nmap scan report for 192.168.0.208 [host down]
Nmap scan report for 192.168.0.209 [host down]
Nmap scan report for 192.168.0.210 [host down]
Nmap scan report for 192.168.0.211 [host down]
Nmap scan report for 192.168.0.212 [host down]
Nmap scan report for 192.168.0.213 [host down]
Nmap scan report for 192.168.0.214 [host down]
Nmap scan report for 192.168.0.215 [host down]
Nmap scan report for 192.168.0.216 [host down]
Nmap scan report for 192.168.0.217 [host down]
Nmap scan report for 192.168.0.218 [host down]
Nmap scan report for 192.168.0.219 [host down]
Nmap scan report for 192.168.0.220 [host down]
Nmap scan report for 192.168.0.221 [host down]
Nmap scan report for 192.168.0.222 [host down]
Nmap scan report for 192.168.0.223 [host down]
Nmap scan report for 192.168.0.224 [host down]
Nmap scan report for 192.168.0.225 [host down]
Nmap scan report for 192.168.0.226 [host down]
Nmap scan report for 192.168.0.227 [host down]
Nmap scan report for 192.168.0.228 [host down]
Nmap scan report for 192.168.0.229 [host down]
Nmap scan report for 192.168.0.230 [host down]
Nmap scan report for 192.168.0.231 [host down]
Nmap scan report for 192.168.0.232 [host down]
Nmap scan report for 192.168.0.233 [host down]
Nmap scan report for 192.168.0.234 [host down]
Nmap scan report for 192.168.0.235 [host down]
Nmap scan report for 192.168.0.236 [host down]
Nmap scan report for 192.168.0.237 [host down]
Nmap scan report for 192.168.0.238 [host down]
Nmap scan report for 192.168.0.239 [host down]
Nmap scan report for 192.168.0.240 [host down]
Nmap scan report for 192.168.0.241 [host down]
Nmap scan report for 192.168.0.242 [host down]
Nmap scan report for 192.168.0.243 [host down]
Nmap scan report for 192.168.0.244 [host down]
Nmap scan report for 192.168.0.245 [host down]
Nmap scan report for 192.168.0.246 [host down]
Nmap scan report for 192.168.0.247 [host down]
Nmap scan report for 192.168.0.248 [host down]
Nmap scan report for 192.168.0.249 [host down]
Nmap scan report for 192.168.0.250 [host down]
Nmap scan report for 192.168.0.251 [host down]
Nmap scan report for 192.168.0.252 [host down]
Nmap scan report for 192.168.0.253 [host down]
Nmap scan report for 192.168.0.254 [host down]
Nmap scan report for 192.168.0.255 [host down]
Increasing send delay for 192.168.0.1 from 0 to 5 due to 20 out of 66 dropped probes since last increase.
Increasing send delay for 192.168.0.54 from 0 to 5 due to 22 out of 73 dropped probes since last increase.
Increasing send delay for 192.168.0.56 from 0 to 5 due to 22 out of 73 dropped probes since last increase.
Increasing send delay for 192.168.0.1 from 10 to 20 due to 17 out of 56 dropped probes since last increase.
Increasing send delay for 192.168.0.56 from 40 to 80 due to 17 out of 56 dropped probes since last increase.
Increasing send delay for 192.168.0.56 from 80 to 160 due to 11 out of 11 dropped probes since last increase.
Increasing send delay for 192.168.0.56 from 160 to 320 due to 11 out of 11 dropped probes since last increase.
Increasing send delay for 192.168.0.54 from 160 to 320 due to 11 out of 24 dropped probes since last increase.
Increasing send delay for 192.168.0.54 from 320 to 640 due to 11 out of 11 dropped probes since last increase.
Warning: 192.168.0.54 giving up on port because retransmission cap hit (10).
Warning: 192.168.0.56 giving up on port because retransmission cap hit (10).
WARNING: Service 192.168.0.1:5000 had already soft-matched upnp, but now soft-matched rtsp; ignoring second valueWARNING: Service 192.168.0.1:5000 had already soft-matched upnp, but now soft-matched rtsp; ignoring second value

WARNING: Service 192.168.0.1:5000 had already soft-matched upnp, but now soft-matched sip; ignoring second valueWARNING: Service 192.168.0.1:5000 had already soft-matched upnp, but now soft-matched sip; ignoring second value

Nmap scan report for 192.168.0.1
Host is up (0.0028s latency).
Not shown: 995 closed ports
PORT     STATE    SERVICE         VERSION
80/tcp   open     http            lighttpd
443/tcp  open     ssl/http        lighttpd
5000/tcp open     upnp            MiniUPnP 1.9 (UPnP 1.1)
| fingerprint-strings: 
|   FourOhFourRequest, GetRequest: 
|     HTTP/1.0 404 Not Found
|     Content-Type: text/html
|     Connection: close
|     Content-Length: 134
|     Server: RedHatEnterpriseServer/6.10 UPnP/1.1 MiniUPnPd/1.9
|     Ext:
|     <HTML><HEAD><TITLE>404 Not Found</TITLE></HEAD><BODY><H1>Not Found</H1>The requested URL was not found on this server.</BODY></HTML>
|   GenericLines: 
|     501 Not Implemented
|     Content-Type: text/html
|     Connection: close
|     Content-Length: 149
|     Server: RedHatEnterpriseServer/6.10 UPnP/1.1 MiniUPnPd/1.9
|     Ext:
|     <HTML><HEAD><TITLE>501 Not Implemented</TITLE></HEAD><BODY><H1>Not Implemented</H1>The HTTP Method is not implemented by this server.</BODY></HTML>
|   HTTPOptions: 
|     HTTP/1.0 501 Not Implemented
|     Content-Type: text/html
|     Connection: close
|     Content-Length: 149
|     Server: RedHatEnterpriseServer/6.10 UPnP/1.1 MiniUPnPd/1.9
|     Ext:
|     <HTML><HEAD><TITLE>501 Not Implemented</TITLE></HEAD><BODY><H1>Not Implemented</H1>The HTTP Method is not implemented by this server.</BODY></HTML>
|   RTSPRequest: 
|     RTSP/1.0 501 Not Implemented
|     Content-Type: text/html
|     Connection: close
|     Content-Length: 149
|     Server: RedHatEnterpriseServer/6.10 UPnP/1.1 MiniUPnPd/1.9
|     Ext:
|_    <HTML><HEAD><TITLE>501 Not Implemented</TITLE></HEAD><BODY><H1>Not Implemented</H1>The HTTP Method is not implemented by this server.</BODY></HTML>
| vulners: 
|   cpe:/a:miniupnp_project:miniupnpd:1.9: 
|     	CVE-2017-8798	7.5	https://vulners.com/cve/CVE-2017-8798
|     	CVE-2014-3985	5.0	https://vulners.com/cve/CVE-2014-3985
|_    	CVE-2017-1000494	4.6	https://vulners.com/cve/CVE-2017-1000494
8081/tcp filtered blackice-icecap
8082/tcp filtered blackice-alerts
1 service unrecognized despite returning data. If you know the service/version, please submit the following fingerprint at https://nmap.org/cgi-bin/submit.cgi?new-service :
SF-Port5000-TCP:V=7.80%I=7%D=2/28%Time=5E59310C%P=armv6-alpine-linux-musle
SF:abihf%r(GenericLines,130,"\x20501\x20Not\x20Implemented\r\nContent-Type
SF::\x20text/html\r\nConnection:\x20close\r\nContent-Length:\x20149\r\nSer
SF:ver:\x20RedHatEnterpriseServer/6\.10\x20UPnP/1\.1\x20MiniUPnPd/1\.9\r\n
SF:Ext:\r\n\r\n<HTML><HEAD><TITLE>501\x20Not\x20Implemented</TITLE></HEAD>
SF:<BODY><H1>Not\x20Implemented</H1>The\x20HTTP\x20Method\x20is\x20not\x20
SF:implemented\x20by\x20this\x20server\.</BODY></HTML>\r\n")%r(GetRequest,
SF:123,"HTTP/1\.0\x20404\x20Not\x20Found\r\nContent-Type:\x20text/html\r\n
SF:Connection:\x20close\r\nContent-Length:\x20134\r\nServer:\x20RedHatEnte
SF:rpriseServer/6\.10\x20UPnP/1\.1\x20MiniUPnPd/1\.9\r\nExt:\r\n\r\n<HTML>
SF:<HEAD><TITLE>404\x20Not\x20Found</TITLE></HEAD><BODY><H1>Not\x20Found</
SF:H1>The\x20requested\x20URL\x20was\x20not\x20found\x20on\x20this\x20serv
SF:er\.</BODY></HTML>\r\n")%r(RTSPRequest,138,"RTSP/1\.0\x20501\x20Not\x20
SF:Implemented\r\nContent-Type:\x20text/html\r\nConnection:\x20close\r\nCo
SF:ntent-Length:\x20149\r\nServer:\x20RedHatEnterpriseServer/6\.10\x20UPnP
SF:/1\.1\x20MiniUPnPd/1\.9\r\nExt:\r\n\r\n<HTML><HEAD><TITLE>501\x20Not\x2
SF:0Implemented</TITLE></HEAD><BODY><H1>Not\x20Implemented</H1>The\x20HTTP
SF:\x20Method\x20is\x20not\x20implemented\x20by\x20this\x20server\.</BODY>
SF:</HTML>\r\n")%r(HTTPOptions,138,"HTTP/1\.0\x20501\x20Not\x20Implemented
SF:\r\nContent-Type:\x20text/html\r\nConnection:\x20close\r\nContent-Lengt
SF:h:\x20149\r\nServer:\x20RedHatEnterpriseServer/6\.10\x20UPnP/1\.1\x20Mi
SF:niUPnPd/1\.9\r\nExt:\r\n\r\n<HTML><HEAD><TITLE>501\x20Not\x20Implemente
SF:d</TITLE></HEAD><BODY><H1>Not\x20Implemented</H1>The\x20HTTP\x20Method\
SF:x20is\x20not\x20implemented\x20by\x20this\x20server\.</BODY></HTML>\r\n
SF:")%r(FourOhFourRequest,123,"HTTP/1\.0\x20404\x20Not\x20Found\r\nContent
SF:-Type:\x20text/html\r\nConnection:\x20close\r\nContent-Length:\x20134\r
SF:\nServer:\x20RedHatEnterpriseServer/6\.10\x20UPnP/1\.1\x20MiniUPnPd/1\.
SF:9\r\nExt:\r\n\r\n<HTML><HEAD><TITLE>404\x20Not\x20Found</TITLE></HEAD><
SF:BODY><H1>Not\x20Found</H1>The\x20requested\x20URL\x20was\x20not\x20foun
SF:d\x20on\x20this\x20server\.</BODY></HTML>\r\n");

Nmap scan report for 192.168.0.10
Host is up (0.000045s latency).
Not shown: 988 closed ports
PORT     STATE SERVICE    VERSION
25/tcp   open  smtp       Postfix smtpd
53/tcp   open  domain     dnsmasq pi-hole-2.80
80/tcp   open  http
| fingerprint-strings: 
|   GetRequest, HTTPOptions: 
|     HTTP/1.0 200 OK
|     Date: Fri, 28 Feb 2020 15:26:04 GMT
|     Content-Length: 1890
|     Content-Type: text/html; charset=utf-8
|     <!DOCTYPE html>
|     <html>
|     <meta name="viewport" content="width=device-width, initial-scale=1.0"> 
|     <head>
|     <title>//</title>
|     <style>
|     body {
|     background-color: #242424;
|     font-family: 'Roboto', sans-serif
|     text-decoration: none;
|     color: #e2e2e2
|     a:hover {
|     color: #fff
|     table {
|     color: #999;
|     text-align: center;
|     margin: auto;
|     padding-top: 15vh
|     padding: 10px 35px
|     padding: 30px
|     color: #999;
|     position: absolute;
|     bottom: 0;
|     right: 0;
|     padding: 5px;
|     margin: 0
|     input {
|     width: 100%!;(MISSING)
|     height: 100%!;(MISSING)
|     color: #999;
|     background-color: #242424;
|     border: none;
|     font-size: 16px;
|     font-family: 'Roboto', sans-serif;
|     text-align: center
|_    </style>
443/tcp  open  tcpwrapped
3000/tcp open  ppp?
| fingerprint-strings: 
|   GenericLines, Help: 
|     HTTP/1.1 400 Bad Request
|     Content-Type: text/plain; charset=utf-8
|     Connection: close
|     Request
|   GetRequest: 
|     HTTP/1.0 302 Found
|     Cache-Control: no-cache
|     Content-Type: text/html; charset=utf-8
|     Expires: -1
|     Location: /login
|     Pragma: no-cache
|     Set-Cookie: redirect_to=%2F; Path=/; HttpOnly; SameSite=Lax
|     X-Frame-Options: deny
|     Date: Fri, 28 Feb 2020 15:26:04 GMT
|     Content-Length: 29
|     href="/login">Found</a>.
|   HTTPOptions: 
|     HTTP/1.0 404 Not Found
|     Cache-Control: no-cache
|     Content-Type: text/html; charset=UTF-8
|     Expires: -1
|     Pragma: no-cache
|     X-Frame-Options: deny
|     Date: Fri, 28 Feb 2020 15:26:09 GMT
|     <!DOCTYPE html>
|     <html lang="en">
|     <head>
|     <script>
|     !(function() {
|     ('PerformanceLongTaskTiming' in window) {
|     (window.__tti = { e: [] });
|     PerformanceObserver(function(l) {
|     g.e.concat(l.getEntries());
|     g.o.observe({ entryTypes: ['longtask'] });
|     })();
|     </script>
|     <meta charset="utf-8" />
|     <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
|     <meta name="viewport" content="width=device-width" />
|     <meta name="theme-color" content="#000" />
|     <title>Grafana</title>
|     <base href="/" />
|     <link
|     rel="preload"
|_    href="public/fonts/roboto/RxZJdnzeo3R5zSe
3306/tcp open  mysql      MySQL 5.5.5-10.1.43-MariaDB-0ubuntu0.18.04.1
| vulners: 
|   MySQL 5.5.5-10.1.43-MariaDB-0ubuntu0.18.04.1: 
|_    	NODEJS:602	0.0	https://vulners.com/nodejs/NODEJS:602
4444/tcp open  http       Jetty 9.4.8.v20180619
|_http-server-header: Jetty(9.4.8.v20180619)
5900/tcp open  tcpwrapped
8080/tcp open  http       lighttpd 1.4.45
|_http-server-header: lighttpd/1.4.45
| vulners: 
|   cpe:/a:lighttpd:lighttpd:1.4.45: 
|_    	CVE-2018-19052	5.0	https://vulners.com/cve/CVE-2018-19052
9080/tcp open  http       nginx 1.16.1
|_http-server-header: nginx/1.16.1
9090/tcp open  http       Golang net/http server (Go-IPFS json-rpc or InfluxDB API)
9100/tcp open  jetdirect?
2 services unrecognized despite returning data. If you know the service/version, please submit the following fingerprints at https://nmap.org/cgi-bin/submit.cgi?new-service :
==============NEXT SERVICE FINGERPRINT (SUBMIT INDIVIDUALLY)==============
SF-Port80-TCP:V=7.80%I=7%D=2/28%Time=5E59310C%P=armv6-alpine-linux-musleab
SF:ihf%r(GetRequest,7D8,"HTTP/1\.0\x20200\x20OK\r\nDate:\x20Fri,\x2028\x20
SF:Feb\x202020\x2015:26:04\x20GMT\r\nContent-Length:\x201890\r\nContent-Ty
SF:pe:\x20text/html;\x20charset=utf-8\r\n\r\n<!DOCTYPE\x20html>\n<html>\n<
SF:meta\x20name=\"viewport\"\x20content=\"width=device-width,\x20initial-s
SF:cale=1\.0\">\x20\n<head>\n\t<title>//</title>\n\t<style>\n\t\tbody\x20{
SF:\n\t\t\tbackground-color:\x20#242424;\n\t\t\tfont-family:\x20'Roboto',\
SF:x20sans-serif\n\t\t}\n\n\t\ta\x20{\n\t\t\ttext-decoration:\x20none;\n\t
SF:\t\tcolor:\x20#e2e2e2\n\t\t}\n\n\t\ta:hover\x20{\n\t\t\tcolor:\x20#fff\
SF:n\t\t}\n\n\t\ttable\x20{\n\t\t\tcolor:\x20#999;\n\t\t\ttext-align:\x20c
SF:enter;\n\t\t\tmargin:\x20auto;\n\t\t\tpadding-top:\x2015vh\n\t\t}\n\n\t
SF:\ttd\x20{\n\t\t\tpadding:\x2010px\x2035px\n\t\t}\n\n\t\tth\x20{\n\t\t\t
SF:padding:\x2030px\n\t\t}\n\n\t\tp\x20{\n\t\t\tcolor:\x20#999;\n\t\t\tpos
SF:ition:\x20absolute;\n\t\t\tbottom:\x200;\n\t\t\tright:\x200;\n\t\t\tpad
SF:ding:\x205px;\n\t\t\tmargin:\x200\n\t\t}\n\n\t\tinput\x20{\n\t\t\twidth
SF::\x20100%!;\(MISSING\)\n\t\t\theight:\x20100%!;\(MISSING\)\n\t\t\tcolor
SF::\x20#999;\n\t\t\tbackground-color:\x20#242424;\n\t\t\tborder:\x20none;
SF:\n\t\t\tfont-size:\x2016px;\n\t\t\tfont-family:\x20'Roboto',\x20sans-se
SF:rif;\n\t\t\ttext-align:\x20center\n\t\t}\n\t</style>\n</")%r(HTTPOption
SF:s,7D8,"HTTP/1\.0\x20200\x20OK\r\nDate:\x20Fri,\x2028\x20Feb\x202020\x20
SF:15:26:04\x20GMT\r\nContent-Length:\x201890\r\nContent-Type:\x20text/htm
SF:l;\x20charset=utf-8\r\n\r\n<!DOCTYPE\x20html>\n<html>\n<meta\x20name=\"
SF:viewport\"\x20content=\"width=device-width,\x20initial-scale=1\.0\">\x2
SF:0\n<head>\n\t<title>//</title>\n\t<style>\n\t\tbody\x20{\n\t\t\tbackgro
SF:und-color:\x20#242424;\n\t\t\tfont-family:\x20'Roboto',\x20sans-serif\n
SF:\t\t}\n\n\t\ta\x20{\n\t\t\ttext-decoration:\x20none;\n\t\t\tcolor:\x20#
SF:e2e2e2\n\t\t}\n\n\t\ta:hover\x20{\n\t\t\tcolor:\x20#fff\n\t\t}\n\n\t\tt
SF:able\x20{\n\t\t\tcolor:\x20#999;\n\t\t\ttext-align:\x20center;\n\t\t\tm
SF:argin:\x20auto;\n\t\t\tpadding-top:\x2015vh\n\t\t}\n\n\t\ttd\x20{\n\t\t
SF:\tpadding:\x2010px\x2035px\n\t\t}\n\n\t\tth\x20{\n\t\t\tpadding:\x2030p
SF:x\n\t\t}\n\n\t\tp\x20{\n\t\t\tcolor:\x20#999;\n\t\t\tposition:\x20absol
SF:ute;\n\t\t\tbottom:\x200;\n\t\t\tright:\x200;\n\t\t\tpadding:\x205px;\n
SF:\t\t\tmargin:\x200\n\t\t}\n\n\t\tinput\x20{\n\t\t\twidth:\x20100%!;\(MI
SF:SSING\)\n\t\t\theight:\x20100%!;\(MISSING\)\n\t\t\tcolor:\x20#999;\n\t\
SF:t\tbackground-color:\x20#242424;\n\t\t\tborder:\x20none;\n\t\t\tfont-si
SF:ze:\x2016px;\n\t\t\tfont-family:\x20'Roboto',\x20sans-serif;\n\t\t\ttex
SF:t-align:\x20center\n\t\t}\n\t</style>\n</");
==============NEXT SERVICE FINGERPRINT (SUBMIT INDIVIDUALLY)==============
SF-Port3000-TCP:V=7.80%I=7%D=2/28%Time=5E59310C%P=armv6-alpine-linux-musle
SF:abihf%r(GenericLines,67,"HTTP/1\.1\x20400\x20Bad\x20Request\r\nContent-
SF:Type:\x20text/plain;\x20charset=utf-8\r\nConnection:\x20close\r\n\r\n40
SF:0\x20Bad\x20Request")%r(GetRequest,132,"HTTP/1\.0\x20302\x20Found\r\nCa
SF:che-Control:\x20no-cache\r\nContent-Type:\x20text/html;\x20charset=utf-
SF:8\r\nExpires:\x20-1\r\nLocation:\x20/login\r\nPragma:\x20no-cache\r\nSe
SF:t-Cookie:\x20redirect_to=%2F;\x20Path=/;\x20HttpOnly;\x20SameSite=Lax\r
SF:\nX-Frame-Options:\x20deny\r\nDate:\x20Fri,\x2028\x20Feb\x202020\x2015:
SF:26:04\x20GMT\r\nContent-Length:\x2029\r\n\r\n<a\x20href=\"/login\">Foun
SF:d</a>\.\n\n")%r(Help,67,"HTTP/1\.1\x20400\x20Bad\x20Request\r\nContent-
SF:Type:\x20text/plain;\x20charset=utf-8\r\nConnection:\x20close\r\n\r\n40
SF:0\x20Bad\x20Request")%r(HTTPOptions,66C1,"HTTP/1\.0\x20404\x20Not\x20Fo
SF:und\r\nCache-Control:\x20no-cache\r\nContent-Type:\x20text/html;\x20cha
SF:rset=UTF-8\r\nExpires:\x20-1\r\nPragma:\x20no-cache\r\nX-Frame-Options:
SF:\x20deny\r\nDate:\x20Fri,\x2028\x20Feb\x202020\x2015:26:09\x20GMT\r\n\r
SF:\n<!DOCTYPE\x20html>\n<html\x20lang=\"en\">\n\x20\x20<head>\n\x20\x20\x
SF:20\x20<script>\n\x20\x20\x20\x20\x20\x20\n\x20\x20\x20\x20\x20\x20!\(fu
SF:nction\(\)\x20{\n\x20\x20\x20\x20\x20\x20\x20\x20if\x20\('PerformanceLo
SF:ngTaskTiming'\x20in\x20window\)\x20{\n\x20\x20\x20\x20\x20\x20\x20\x20\
SF:x20\x20var\x20g\x20=\x20\(window\.__tti\x20=\x20{\x20e:\x20\[\]\x20}\);
SF:\n\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20g\.o\x20=\x20new\x20Performan
SF:ceObserver\(function\(l\)\x20{\n\x20\x20\x20\x20\x20\x20\x20\x20\x20\x2
SF:0\x20\x20g\.e\x20=\x20g\.e\.concat\(l\.getEntries\(\)\);\n\x20\x20\x20\
SF:x20\x20\x20\x20\x20\x20\x20}\);\n\x20\x20\x20\x20\x20\x20\x20\x20\x20\x
SF:20g\.o\.observe\({\x20entryTypes:\x20\['longtask'\]\x20}\);\n\x20\x20\x
SF:20\x20\x20\x20\x20\x20}\n\x20\x20\x20\x20\x20\x20}\)\(\);\n\n\x20\x20\x
SF:20\x20</script>\n\x20\x20\x20\x20<meta\x20charset=\"utf-8\"\x20/>\n\x20
SF:\x20\x20\x20<meta\x20http-equiv=\"X-UA-Compatible\"\x20content=\"IE=edg
SF:e,chrome=1\"\x20/>\n\x20\x20\x20\x20<meta\x20name=\"viewport\"\x20conte
SF:nt=\"width=device-width\"\x20/>\n\x20\x20\x20\x20<meta\x20name=\"theme-
SF:color\"\x20content=\"#000\"\x20/>\n\n\x20\x20\x20\x20<title>Grafana</ti
SF:tle>\n\n\x20\x20\x20\x20<base\x20href=\"/\"\x20/>\n\n\x20\x20\x20\x20<l
SF:ink\n\x20\x20\x20\x20\x20\x20rel=\"preload\"\n\x20\x20\x20\x20\x20\x20h
SF:ref=\"public/fonts/roboto/RxZJdnzeo3R5zSe");
Service Info: Host:  hecuba.home

Nmap scan report for 192.168.0.11
Host is up (0.0068s latency).
Not shown: 998 closed ports
PORT    STATE SERVICE   VERSION
80/tcp  open  http      Boa HTTPd 0.94.14rc21
|_http-server-header: Boa/0.94.14rc21
| vulners: 
|   cpe:/a:boa:boa:0.94.14rc21: 
|_    	CVE-2009-4496	5.0	https://vulners.com/cve/CVE-2009-4496
443/tcp open  ssl/https axhttpd/1.4.9
| fingerprint-strings: 
|   FourOhFourRequest: 
|     HTTP/1.1 200 OK
|     Content-Type: text/html
|     <html><head><meta HTTP-EQUIV="REFRESH" content="0; url=/index.html"></head>
|     <body>
|     <title></title>
|     </body></html>
|   GenericLines, RTSPRequest: 
|     HTTP/1.1 200 OK
|     Server: axhttpd/1.4.9
|     Content-Type: text/html
|     Content-Length: 64
|     Date: Sat, 03 Jan 1970 04:06:15 GMT
|     Last-Modified: Thu, 01 Jan 1970 00:00:19 GMT
|     Expires: Sat, 03 Jan 1970 04:11:15 GMT
|     <meta http-equiv="refresh" content="0;url=http://belkin.range">
|   GetRequest, HTTPOptions: 
|     HTTP/1.1 200 OK
|     Server: axhttpd/1.4.9
|     Content-Type: text/html
|     Content-Length: 64
|     Date: Sat, 03 Jan 1970 04:06:10 GMT
|     Last-Modified: Thu, 01 Jan 1970 00:00:19 GMT
|     Expires: Sat, 03 Jan 1970 04:11:10 GMT
|     <meta http-equiv="refresh" content="0;url=http://belkin.range">
|   SIPOptions: 
|     HTTP/1.1 200 OK
|     Server: axhttpd/1.4.9
|     Content-Type: text/html
|     Content-Length: 64
|     Date: Sat, 03 Jan 1970 04:07:23 GMT
|     Last-Modified: Thu, 01 Jan 1970 00:00:19 GMT
|     Expires: Sat, 03 Jan 1970 04:12:23 GMT
|_    <meta http-equiv="refresh" content="0;url=http://belkin.range">
|_http-server-header: axhttpd/1.4.9
1 service unrecognized despite returning data. If you know the service/version, please submit the following fingerprint at https://nmap.org/cgi-bin/submit.cgi?new-service :
SF-Port443-TCP:V=7.80%T=SSL%I=7%D=2/28%Time=5E593112%P=armv6-alpine-linux-
SF:musleabihf%r(GetRequest,10A,"HTTP/1\.1\x20200\x20OK\nServer:\x20axhttpd
SF:/1\.4\.9\nContent-Type:\x20text/html\nContent-Length:\x2064\nDate:\x20S
SF:at,\x2003\x20Jan\x201970\x2004:06:10\x20GMT\nLast-Modified:\x20Thu,\x20
SF:01\x20Jan\x201970\x2000:00:19\x20GMT\nExpires:\x20Sat,\x2003\x20Jan\x20
SF:1970\x2004:11:10\x20GMT\n\n<meta\x20http-equiv=\"refresh\"\x20content=\
SF:"0;url=http://belkin\.range\">\n")%r(HTTPOptions,10A,"HTTP/1\.1\x20200\
SF:x20OK\nServer:\x20axhttpd/1\.4\.9\nContent-Type:\x20text/html\nContent-
SF:Length:\x2064\nDate:\x20Sat,\x2003\x20Jan\x201970\x2004:06:10\x20GMT\nL
SF:ast-Modified:\x20Thu,\x2001\x20Jan\x201970\x2000:00:19\x20GMT\nExpires:
SF:\x20Sat,\x2003\x20Jan\x201970\x2004:11:10\x20GMT\n\n<meta\x20http-equiv
SF:=\"refresh\"\x20content=\"0;url=http://belkin\.range\">\n")%r(FourOhFou
SF:rRequest,9B,"HTTP/1\.1\x20200\x20OK\nContent-Type:\x20text/html\n\n<htm
SF:l><head><meta\x20HTTP-EQUIV=\"REFRESH\"\x20content=\"0;\x20url=/index\.
SF:html\"></head>\n<body>\n<title></title>\n</body></html>\n")%r(GenericLi
SF:nes,10A,"HTTP/1\.1\x20200\x20OK\nServer:\x20axhttpd/1\.4\.9\nContent-Ty
SF:pe:\x20text/html\nContent-Length:\x2064\nDate:\x20Sat,\x2003\x20Jan\x20
SF:1970\x2004:06:15\x20GMT\nLast-Modified:\x20Thu,\x2001\x20Jan\x201970\x2
SF:000:00:19\x20GMT\nExpires:\x20Sat,\x2003\x20Jan\x201970\x2004:11:15\x20
SF:GMT\n\n<meta\x20http-equiv=\"refresh\"\x20content=\"0;url=http://belkin
SF:\.range\">\n")%r(RTSPRequest,10A,"HTTP/1\.1\x20200\x20OK\nServer:\x20ax
SF:httpd/1\.4\.9\nContent-Type:\x20text/html\nContent-Length:\x2064\nDate:
SF:\x20Sat,\x2003\x20Jan\x201970\x2004:06:15\x20GMT\nLast-Modified:\x20Thu
SF:,\x2001\x20Jan\x201970\x2000:00:19\x20GMT\nExpires:\x20Sat,\x2003\x20Ja
SF:n\x201970\x2004:11:15\x20GMT\n\n<meta\x20http-equiv=\"refresh\"\x20cont
SF:ent=\"0;url=http://belkin\.range\">\n")%r(SIPOptions,10A,"HTTP/1\.1\x20
SF:200\x20OK\nServer:\x20axhttpd/1\.4\.9\nContent-Type:\x20text/html\nCont
SF:ent-Length:\x2064\nDate:\x20Sat,\x2003\x20Jan\x201970\x2004:07:23\x20GM
SF:T\nLast-Modified:\x20Thu,\x2001\x20Jan\x201970\x2000:00:19\x20GMT\nExpi
SF:res:\x20Sat,\x2003\x20Jan\x201970\x2004:12:23\x20GMT\n\n<meta\x20http-e
SF:quiv=\"refresh\"\x20content=\"0;url=http://belkin\.range\">\n");

Nmap scan report for 192.168.0.13
Host is up (0.0094s latency).
All 1000 scanned ports on 192.168.0.13 are closed

Nmap scan report for 192.168.0.54
Host is up (0.0097s latency).
Not shown: 842 closed ports, 157 filtered ports
PORT      STATE SERVICE    VERSION
62078/tcp open  tcpwrapped

Nmap scan report for 192.168.0.56
Host is up (0.092s latency).
Not shown: 764 closed ports, 235 filtered ports
PORT      STATE SERVICE    VERSION
62078/tcp open  tcpwrapped

Read data files from: /usr/bin/../share/nmap
Service detection performed. Please report any incorrect results at https://nmap.org/submit/ .
# Nmap done at Fri Feb 28 15:28:25 2020 -- 256 IP addresses (6 hosts up) scanned in 3484.97 seconds

```
