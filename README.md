## certinfo

**certinfo** is a powerful SSL certificate scraping tool that extracts domain names from SSL certificates of arbitrary hosts. It supports both basic certificate data extraction and recursive subdomain enumeration through Certificate Subject Alternative Names (SANs).

## Features

- 🚀 **High Performance**: Multi-threaded concurrent certificate processing (default: 50 workers)
- 🔄 **Recursive Enumeration**: Automatically discover subdomains through certificate SANs
- 📊 **Multiple Output Formats**: JSON, CSV, or plain text output
- 🎯 **Flexible Input**: Supports domains, IPs, and custom ports
- ⚡ **Real-time Output**: Stream results as they're discovered
- 🔍 **Detailed Information**: Extract issuer, validity, CN, Organization, and SANs

## Installation

### Install via Go
```
go install github.com/rix4uni/certinfo@latest
```

### Download Prebuilt Binaries
```
wget https://github.com/rix4uni/certinfo/releases/download/v0.0.6/certinfo-linux-amd64-0.0.6.tgz
tar -xvzf certinfo-linux-amd64-0.0.6.tgz
rm -rf certinfo-linux-amd64-0.0.6.tgz
mv certinfo ~/go/bin/certinfo
```

Or download the [latest release](https://github.com/rix4uni/certinfo/releases) for your platform.

### Compile from Source
```
git clone --depth 1 https://github.com/rix4uni/certinfo.git
cd certinfo; go install
```

## Usage

### Basic Usage

**Single Target:**
```yaml
echo "google.com" | certinfo -silent
```

**Multiple Targets:**
```yaml
cat targets.txt | certinfo -silent
```

**With Custom Port:**
```yaml
echo "example.com:8443" | certinfo -silent
```

## Command-Line Options

```yaml
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
  -recursive
        recursive subdomain enumeration from certificate SANs
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

## Output Formats

### Default Output
By default, certinfo prints all SAN domains found in certificates:
```yaml
echo "google.com" | certinfo -silent
google.com
*.google.com
*.google.co.uk
*.google.fr
*.accounts.google.com
...
```

### JSON Output
Get structured JSON output:
```yaml
echo "google.com" | certinfo -silent -json
{
  "host": "google.com:443",
  "Issued_To": {
    "Common_Name_(CN)": "*.google.com",
    "Organization_(O)": "Google LLC"
  },
  "Issued_By": {
    "Common_Name_(CN)": "GTS CA 1C3",
    "Organization_(O)": "Google Trust Services"
  },
  "Validity_Period": {
    "Issued_On": "2024-01-15T10:00:00Z",
    "Expires_On": "2024-04-15T10:00:00Z"
  },
  "Certificate_Subject_Alternative_Name": [
    "*.google.com",
    "*.google.co.uk",
    "*.accounts.google.com",
    ...
  ]
}
```

### CSV Output
Get comma-separated values:
```yaml
echo "google.com" | certinfo -silent -csv
Host,IssuedTo_CommonName,IssuedTo_Organization,IssuedBy_CommonName,IssuedBy_Organization,IssuedOn,ExpiresOn,SubjectAlternativeNames
google.com:443,*.google.com,Google LLC,GTS CA 1C3,Google Trust Services,2024-01-15T10:00:00Z,2024-04-15T10:00:00Z,"*.google.com, *.google.co.uk, *.accounts.google.com"
```

### Specialized Outputs

**SAN Details:**
```yaml
echo "207.207.12.80" | certinfo -silent -san
207.207.12.80:443 [wwwmicrolb.informatica.com, trust.informatica.com, diaku.com, careers.informatica.com]
```

**Issue Date:**
```yaml
echo "google.com" | certinfo -silent -issued
google.com:443 [2024-01-15T10:00:00Z]
```

**Expiration Date:**
```yaml
echo "google.com" | certinfo -silent -expires
google.com:443 [2024-04-15T10:00:00Z]
```

**Issue + Expiration + SAN:**
```yaml
echo "google.com" | certinfo -silent -issued -expires -san
google.com:443 [2024-01-15T10:00:00Z] [2024-04-15T10:00:00Z] [*.google.com, *.google.co.uk, *.accounts.google.com]
```

## Recursive Subdomain Enumeration

The `-recursive` flag enables powerful recursive subdomain discovery by automatically processing domains found in certificate SANs.

### How It Works

1. **Initial Phase**: Process the input domain(s) and extract SANs from their certificates
2. **Recursion Phase**: For each discovered SAN domain:
   - Strip wildcard prefixes (`*.` → ``)
   - Fetch the certificate for the cleaned domain
   - Extract new SANs from that certificate
3. **Iteration**: Repeat until no new domains are discovered
4. **Real-time Output**: Results are printed immediately as they're discovered

### Example: Recursive vs Non-Recursive

**Non-recursive** (138 domains found):
```yaml
echo "google.com" | certinfo -silent | unew | wc -l
138
```

**Recursive** (663 domains found - 4.8x more!):
```yaml
echo "google.com" | certinfo -silent -recursive | unew | wc -l
663
```


## Advanced Use Cases

### Certificate Monitoring

**Find certificates issued today:**
```yaml
gungnir -r inscope_wildcards.txt | unew | certinfo -silent -issued -today
www.www.internal.moveit.qms.grab.com:443 [2025-01-18T10:58:38Z]
img-ru.shein.com:443 [2025-01-18T10:10:00Z]
```

**Pipe to Nuclei for vulnerability scanning:**
```yaml
gungnir -r inscope_wildcards.txt | unew | certinfo -silent -issued -today | awk '{print $1}' | nuclei
```

### Network Discovery

**From IP address:**
```yaml
echo "207.207.12.80" | certinfo -silent
wwwmicrolb.informatica.com
trust.informatica.com
diaku.com
careers.informatica.com
```

**Recursive discovery on IP:**
```yaml
echo "207.207.12.80" | certinfo -silent -recursive
# Discovers all domains hosted on that IP through their certificates
```

### Bulk Processing

**Process large target lists:**
```yaml
cat large_target_list.txt | certinfo -silent -c 100 -timeout 5s -recursive | tee results.txt
```

**With custom concurrency:**
```yaml
cat targets.txt | certinfo -silent -c 200 -timeout 10s
```

### Analysis & Reporting

**Save as CSV for analysis:**
```yaml
cat targets.txt | certinfo -silent -csv > certificate_analysis.csv
```

**Filter expired certificates:**
```yaml
cat targets.txt | certinfo -silent -expires | grep "2024-01"
```

## Performance

certinfo is optimized for speed and efficiency:

- **Concurrent Processing**: Default 50 workers, configurable up to hundreds
- **Connection Pooling**: Efficient TCP connection reuse
- **Low Memory Footprint**: Streaming processing without buffering
- **Network Optimization**: Configurable timeouts and keep-alives

**Benchmark Example:**
```yaml
time echo "google.com" | certinfo -silent -recursive -c 100 | wc -l
663 domains in ~5-10 seconds
```