FROM python:3.12.4-slim

WORKDIR /opt/app
COPY requirements.*.txt ./

RUN pip install --no-cache-dir -r requirements.llm.txt
RUN pip install --no-cache-dir -r requirements.app.txt
RUN pip install --no-cache-dir -r requirements.telemetry.txt
RUN opentelemetry-bootstrap -a install

COPY . .
