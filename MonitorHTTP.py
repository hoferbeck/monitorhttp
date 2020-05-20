from tcp_latency import measure_latency
import requests
from influxdb import InfluxDBClient
import re
import urllib3
import time
#test

urllib3.disable_warnings()

databasename = "monitoring"
checkedurls = ["http://192.168.1.4", "https://cloud.hoferbeck.at/", "https://192.168.1.103:8443", "https://192.168.1.104:7443",
               "https://192.168.1.100:444", "http://192.168.1.104:8081", "http://192.168.1.100:32400/web/", "http://192.168.1.1", "http://192.168.0.1", "http://10.0.0.1"]

influxclient = InfluxDBClient(host='192.168.1.100', port=8086)
influxclient.create_database(databasename)
influxclient.switch_database(databasename)

#latency = measure_latency(host="192.168.1.100", port=444)
#print(latency[0])

while True:
    for checkedurl in checkedurls:
        domainwithoutprefix = re.sub("^https?://", "", checkedurl)
        try:
            r = requests.head(checkedurl, verify=False, allow_redirects=True)
            print(r.status_code)
            status_code = r.status_code
            respTime = round(r.elapsed.total_seconds(), 3)
            # prints the int of the status code. Find more at httpstatusrappers.com :)
        except requests.ConnectionError:
            print("failed to connect")
            status_code = 999
            respTime = 999.0

        json_body = [
            {
                "measurement": "HTTP Status Code",
                "tags": {
                    "Url": domainwithoutprefix,
                },
                "fields": {
                    "HTTP Status Code": status_code,
                    "Respone Time_float": respTime
                }
            }
        ]
        try:
            influxclient.write_points(json_body)
        except:
            pass

    time.sleep(60)
