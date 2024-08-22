FROM python:3-alpine

WORKDIR /opt/app
COPY . .

RUN pip install --no-cache-dir -r requirements.txt
RUN opentelemetry-bootstrap -a install

CMD [ "python", "app.py" ]
