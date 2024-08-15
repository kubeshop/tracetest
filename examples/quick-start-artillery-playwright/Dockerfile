FROM mcr.microsoft.com/playwright:v1.45.3

WORKDIR /app

COPY . .

RUN npm install

RUN npm install -g artillery artillery-engine-playwright && \
  npm cache clean --force && \
  rm -rf /root/.cache && \
  rm -rf /ms-playwright/firefox* && \
  rm -rf /ms-playwright/webkit*
