FROM python:3.12.4-slim

WORKDIR /opt/app
COPY . .

RUN pip install --no-cache-dir -r requirements.txt

CMD [ "streamlit", "run", "streamlit_app.py", "--server.port=8051", "--server.address=0.0.0.0" ]
