## certinfo

Scrape domain names from SSL certificates of arbitrary hosts

## Installation
```
go install github.com/rix4uni/certinfo@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/certinfo/releases/download/v0.0.5/certinfo-linux-amd64-0.0.5.tgz
tar -xvzf certinfo-linux-amd64-0.0.5.tgz
rm -rf certinfo-linux-amd64-0.0.5.tgz
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
Domain:
```
▶ echo "xapi.stg.xfinity.com" | certinfo -silent
xapi.stg.xfinity.com
cdn.ch2.int.business.comcast.com
cdn.ch2.int.comcast.com
...
...
...
xapi.qa3.xfinity.com
```

IPv4:
```
▶ echo "207.207.12.80" | certinfo -silent
wwwmicrolb.informatica.com
trust.informatica.com
diaku.com
careers.informatica.com
```

I got idea from `https://kaeferjaeger.gay/sni-ip-ranges/google/ipv4_merged_sni.txt`:
```
▶ echo "207.207.12.80" | certinfo -silent -san
207.207.12.80:443 [wwwmicrolb.informatica.com, trust.informatica.com, diaku.com, careers.informatica.com]
```

Get those certificates issued today:
```
gungnir -r inscope_wildcards.txt | unew | certinfo -silent -issued -today
www.www.internal.moveit.qms.grab.com:443 [2025-01-18T10:58:38Z]
img-ru.shein.com:443 [2025-01-18T10:10:00Z]
37081b66-60bf-4872-ac7e-f23446bd4d23.unifi-hosting.ui.com:443 [2025-01-18T10:12:41Z]
```

More:
```
gungnir -r inscope_wildcards.txt | unew | certinfo -silent -issued -today | awk '{print $1}' | nuclei
```
