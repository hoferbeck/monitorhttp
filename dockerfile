FROM python:3
ADD MonitorHTTP.py /
RUN pip install tcp_latency
RUN pip install influxdb
CMD [ "python", "./MonitorHTTP.py" ]