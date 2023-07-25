# certinfo

# Usage
```
echo "104.98.132.228" | go run main.go
cat urls.txt | go run main.go
```

# Output
```
{
  "host": "104.98.132.228",
  "Issued_To": {
    "Common_Name_(CN)": "xapi.stg.xfinity.com",
    "Organization_(O)": "Comcast Corporation"
  },
  "Issued_By": {
    "Common_Name_(CN)": "COMODO RSA Organization Validation Secure Server CA",
    "Organization_(O)": "COMODO CA Limited"
  },
  "Validity_Period": {
    "Expires_On": "2024-06-28T23:59:59Z",
    "Issued_On": "2023-06-29T00:00:00Z"
  },
  "Certificate_Subject_Alternative_Name": [
    "xapi.stg.xfinity.com",
    "cdn.ch2.int.business.comcast.com",
    "cdn.ch2.int.comcast.com",
    "cdn.ch2.stg.business.comcast.com",
    "cdn.ch2.stg.comcast.com",
    "cdn.int.business.comcast.com",
    "cdn.int.comcast.com",
    "cdn.pdc.int.business.comcast.com",
    "cdn.pdc.int.comcast.com",
    "cdn.pdc.stg.business.comcast.com",
    "cdn.pdc.stg.comcast.com",
    "cdn.perf.business.comcast.com",
    "cdn.perf.comcast.com",
    "cdn.stg.business.comcast.com",
    "cdn.stg.comcast.com",
    "cdn.wcdc.int.business.comcast.com",
    "cdn.wcdc.int.comcast.com",
    "cdn.wcdc.perf.business.comcast.com",
    "cdn.wcdc.perf.comcast.com",
    "cdn.wcdc.stg.business.comcast.com",
    "cdn.wcdc.stg.comcast.com",
    "compat.business.int.comcast.com",
    "compat.customer.int.xfinity.com",
    "compat.delivery.int.xfinity.com",
    "compat.www.int.xfinity.com",
    "compat.xapi.int.xfinity.com",
    "delivery.int.xfinity.com",
    "delivery.perf.xfinity.com",
    "delivery.stg.xfinity.com",
    "idm-perf.xfinity.com",
    "login-perf.xfinity.com",
    "oauth-perf.xfinity.com",
    "preview.api.stg.xfinity.com",
    "preview.www.stg.xfinity.com",
    "preview.xapi.stg.xfinity.com",
    "prodtest.business.int.comcast.com",
    "prodtest.customer.int.xfinity.com",
    "prodtest.delivery.int.xfinity.com",
    "prodtest.www.int.xfinity.com",
    "prodtest.xapi.int.xfinity.com",
    "prv.www.int.xfinity.com",
    "prv.www.stg.xfinity.com",
    "services.e2e.xfinity.com",
    "services.int.xfinity.com",
    "services.perf.xfinity.com",
    "services.qa.xfinity.com",
    "services.qa1.xfinity.com",
    "services.qa2.xfinity.com",
    "services.qa3.xfinity.com",
    "services.stg.xfinity.com",
    "www.dev.xfinity.com",
    "www.e2e.xfinity.com",
    "www.int.xfinity.com",
    "www.perf.xfinity.com",
    "www.qa.xfinity.com",
    "www.qa1.xfinity.com",
    "www.qa2.xfinity.com",
    "www.qa3.xfinity.com",
    "www.stg.xfinity.com",
    "www.xapi.stg.xfinity.com",
    "xapi.e2e.xfinity.com",
    "xapi.int.xfinity.com",
    "xapi.perf.xfinity.com",
    "xapi.qa.xfinity.com",
    "xapi.qa1.xfinity.com",
    "xapi.qa2.xfinity.com",
    "xapi.qa3.xfinity.com"
  ]
}
```
