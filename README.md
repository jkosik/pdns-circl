## Passive DNS client for CIRCL PDNS Database - golang implementation
The `pdns-circl` client gets data from [CIRCL PDNS Database](https://www.circl.lu/services/passive-dns/).  
Passive DNS data follows [Passive DNS Common Output Format](https://www.ietf.org/archive/id/draft-dulaunoy-dnsop-passive-dns-cof-01.txt).  

## Installation
Clone this repo and simply use the binary `pdns-circl`. If needed, compile the golang code for your platform: 
```
$ go build pdns-circl.go
```

## Usage
```
$ ./pdns-circl -h
Usage of ./pdns-circl:
  -p string
    	CIRCL PDNS API Password (default "pass")
  -r	Complete raw output for -rrname. Good option for jq processing and filtering. Ignores -rrtype flag.
  -rrname string
    	Domain to lookup, e.g. www.google.com (default "www.circl.lu")
  -rrtype string
    	RR as subfilter, e.g. A, CNAME, AAAA (default "nil")
  -u string
    	CIRCL PDNS API Username (default "user")
```


## Human readable PDNS output for selected domain
```
$ ./pdns-circl -u CIRCL_API_USER -p CIRCL_API_PASSWORD -rrname www.circl.lu 
+++++ Listing CIRCL PDNS records for www.circl.lu +++++

Count = 989255
TimeFirst = 2016-10-07 09:26:02 +0200 CEST
RRType = CNAME
RRName = www.circl.lu
RData = cpab.circl.lu
TimeLast = 2018-10-30 01:56:36 +0100 CET
------------------------------------------
Count = 20426
TimeFirst = 2011-02-22 19:13:37 +0100 CET
RRType = A
RRName = www.circl.lu
RData = 194.154.205.24
TimeLast = 2011-03-04 19:41:17 +0100 CET
------------------------------------------
Count = 23479
TimeFirst = 2011-02-22 19:06:42 +0100 CET
RRType = CNAME
RRName = www.circl.lu
RData = cpa.circl.lu
TimeLast = 2012-02-14 10:31:34 +0100 CET
------------------------------------------
```

## Human readable filtered output for selected Resoure Record.
Subfilter `-rrtype` accepts the following RR types: A, CNAME, AAAA, PTR, SOA, NS, SRV, TXT.  
When `-rrtype` flag is not used or contains any other option, all RR types are listed.
```
$ ./pdns-circl -u CIRCL_API_USER -p CIRCL_API_PASSWORD -rrname www.google.sk -rrtype CNAME
+++++ Listing CIRCL PDNS records for www.google.sk +++++

Count = 4
TimeFirst = 2012-01-19 10:27:27 +0100 CET
RRType = CNAME
RRName = www.google.sk
RData = www-cctld.l.google.com
TimeLast = 2012-01-25 03:02:37 +0100 CET
------------------------------------------
Count = 18
TimeFirst = 2011-09-19 17:11:21 +0200 CEST
RRType = CNAME
RRName = www.google.sk
RData = www.google.com
TimeLast = 2012-01-05 13:17:21 +0100 CET
------------------------------------------
```

## Raw output
`pdns-circl` supports raw data output, when using `-raw` flag. Suitable for automated data processing using external tools, e.g. `jq`.    
`-raw` flag ignores `-rrtype` flag and lists all the RR types.
```
$ ./pdns-circl -u CIRCL_API_USER -p CIRCL_API_PASSWORD -rrname www.google.sk -raw | jq
[{"count": 4, "origin": "https://www.circl.lu/pdns/", "time_first": 1326965247, "rrtype": "CNAME", "rrname": "www.google.sk", "rdata": "www-cctld.l.google.com", "time_last": 1327456957},{"count": 18, "origin": "https://www.circl.lu/pdns/", "time_first": 1316445081, "rrtype": "CNAME", "rrname": "www.google.sk", "rdata": "www.google.com", "time_last": 1325765841},{"count": 2, "origin": "https://www.circl.lu/pdns/", "time_first": 1531249383, "rrtype": "A", "rrname": "www.google.sk", "rdata": "172.217.17.99", "time_last": 1531249383},{"count": 2, "origin": "https://www.circl.lu/pdns/", "time_first": 1527587658, "rrtype": "A", "rrname": "www.google.sk", "rdata": "172.217.2.67", "time_last": 1527587658},{"count": 3, "origin": "https://www.circl.lu/pdns/", "time_first": 1540772271, "rrtype": "A", "rrname": "www.google.sk", "rdata": "216.58.207.67", "time_last": 1540772271},{"count": 19, "origin": "https://www.circl.lu/pdns/", "time_first": 1535562496, "rrtype": "A", "rrname": "www.google.sk", "rdata": "172.217.20.99", "time_last": 1538180401},{"count": 2, "origin": "https://www.circl.lu/pdns/", "time_first": 1528468399, "rrtype": "A", "rrname": "www.google.sk", "rdata": "172.217.8.99", "time_last": 1528468399},{"count": 14, "origin": "https://www.circl.lu/pdns/", "time_first": 1516886195, "rrtype": "A", "rrname": "www.google.sk", "rdata": "172.217.16.67", "time_last": 1517311165},{"count": 6, "origin": "https://www.circl.lu/pdns/", "time_first": 1478091716, "rrtype": "A", "rrname": "www.google.sk", "rdata": "172.217.17.131", "time_last": 1530377664},{"count": 2, "origin": "https://www.circl.lu/pdns/", "time_first": 1539644082, "rrtype": "A", "rrname": "www.google.sk", "rdata": "172.217.17.67", "time_last": 1539644082}]
```

#### JQ processing:
```
$ ./pdns-circl -u CIRCL_API_USER -p CIRCL_API_PASSWORD -rrname www.google.sk -raw | jq
[
  {
    "count": 4,
    "origin": "https://www.circl.lu/pdns/",
    "time_first": 1326965247,
    "rrtype": "CNAME",
    "rrname": "www.google.sk",
    "rdata": "www-cctld.l.google.com",
    "time_last": 1327456957
  },
  {
    "count": 18,
    "origin": "https://www.circl.lu/pdns/",
    "time_first": 1316445081,
    "rrtype": "CNAME",
    "rrname": "www.google.sk",
    "rdata": "www.google.com",
    "time_last": 1325765841
  },
  {
    "count": 2,
    "origin": "https://www.circl.lu/pdns/",
    "time_first": 1531249383,
    "rrtype": "A",
    "rrname": "www.google.sk",
    "rdata": "172.217.17.99",
    "time_last": 1531249383
  },
...snipped...
```




