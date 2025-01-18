## certinfo

Scrape domain names from Certificate Subject Alternative Name

## Installation
```
go install github.com/rix4uni/certinfo@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/certinfo/releases/download/v0.0.4/certinfo-linux-amd64-0.0.4.tgz
tar -xvzf certinfo-linux-amd64-0.0.4.tgz
rm -rf certinfo-linux-amd64-0.0.4.tgz
mv certinfo ~/go/bin/certinfo
```
Or download [binary release](https://github.com/rix4uni/certinfo/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/certinfo.git
cd certinfo; go install
```

## Usage
```
Usage of certinfo:
  -c int
        number of concurrent workers (default 50)
  -csv
        output in CSV format
  -expires
        output only the expiration date
  -issued
        output host, port, and certificate expiration date
  -json
        output in JSON format
  -san
        monitor the san certificate details in a simple format
  -silent
        silent mode.
  -timeout string
        connection timeout duration (e.g. 5s, 10m, 1h) (default "3s")
  -today
        filter results to show only certificates issued today (works only with -issued flag)
  -verbose
        enable verbose logging
  -version
        Print the version of the tool and exit.
```

## Usage
Single Target:
```
▶ echo "xapi.stg.xfinity.com" | certinfo -silent
```

Multiple Targets:
```
▶ cat targets.txt
xapi.stg.xfinity.com
207.207.12.80

▶ cat targets.txt | certinfo -silent
```

## Usage Examples
Domain
```
▶ echo "xapi.stg.xfinity.com" | certinfo -silent
xapi.stg.xfinity.com
cdn.ch2.int.business.comcast.com
cdn.ch2.int.comcast.com
cdn.ch2.stg.business.comcast.com
cdn.ch2.stg.comcast.com
cdn.int.business.comcast.com
cdn.int.comcast.com
cdn.pdc.int.business.comcast.com
cdn.pdc.int.comcast.com
cdn.pdc.stg.business.comcast.com
cdn.pdc.stg.comcast.com
cdn.perf.business.comcast.com
cdn.perf.comcast.com
cdn.stg.business.comcast.com
cdn.stg.comcast.com
cdn.wcdc.int.business.comcast.com
cdn.wcdc.int.comcast.com
cdn.wcdc.perf.business.comcast.com
cdn.wcdc.perf.comcast.com
cdn.wcdc.stg.business.comcast.com
cdn.wcdc.stg.comcast.com
compat.business.int.comcast.com
compat.customer.int.xfinity.com
compat.delivery.int.xfinity.com
compat.www.int.xfinity.com
compat.xapi.int.xfinity.com
delivery.int.xfinity.com
delivery.perf.xfinity.com
delivery.stg.xfinity.com
idm-perf.xfinity.com
login-perf.xfinity.com
oauth-perf.xfinity.com
preview.api.stg.xfinity.com
preview.www.stg.xfinity.com
preview.xapi.stg.xfinity.com
prodtest.business.int.comcast.com
prodtest.customer.int.xfinity.com
prodtest.delivery.int.xfinity.com
prodtest.www.int.xfinity.com
prodtest.xapi.int.xfinity.com
prv.www.int.xfinity.com
prv.www.stg.xfinity.com
services.e2e.xfinity.com
services.int.xfinity.com
services.perf.xfinity.com
services.qa.xfinity.com
services.qa1.xfinity.com
services.qa2.xfinity.com
services.qa3.xfinity.com
services.stg.xfinity.com
ts43-stage-waf.ecs.xm.comcast.com
ts43-stage.ecs.xm.comcast.com
www.dev.xfinity.com
www.e2e.xfinity.com
www.int.xfinity.com
www.perf.xfinity.com
www.qa.xfinity.com
www.qa1.xfinity.com
www.qa2.xfinity.com
www.qa3.xfinity.com
www.stg.xfinity.com
www.xapi.stg.xfinity.com
xapi.e2e.xfinity.com
xapi.int.xfinity.com
xapi.perf.xfinity.com
xapi.qa.xfinity.com
xapi.qa1.xfinity.com
xapi.qa2.xfinity.com
xapi.qa3.xfinity.com
```

IPv4
```
▶ echo "207.207.12.80" | certinfo -silent
wwwmicrolb.informatica.com
trust.informatica.com
diaku.com
careers.informatica.com
```

Idea got from `https://kaeferjaeger.gay/sni-ip-ranges/google/ipv4_merged_sni.txt`
```
▶ echo "207.207.12.80" | certinfo -silent -san
207.207.12.80:443 [wwwmicrolb.informatica.com, trust.informatica.com, diaku.com, careers.informatica.com]
```