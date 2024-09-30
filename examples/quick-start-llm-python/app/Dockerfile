FROM python:3.12.4-slim

WORKDIR /opt/app
COPY requirements.txt ./

RUN pip install --no-cache-dir -r requirements.txt
RUN opentelemetry-bootstrap -a install

COPY . .
