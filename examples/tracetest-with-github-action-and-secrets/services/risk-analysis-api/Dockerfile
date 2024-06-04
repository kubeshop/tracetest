FROM python:3.10.1-slim

WORKDIR /opt/app
COPY . .

RUN pip install --no-cache-dir -r requirements.txt
RUN opentelemetry-bootstrap -a install

CMD [ "python", "app.py" ]
