FROM node:20.10.0

FROM mcr.microsoft.com/playwright:focal

WORKDIR /app 

ENV PATH /app/node_modules/.bin:$PATH

COPY package.json /app/
COPY playwright/ /app/playwright/
COPY tsconfig.json /app/
COPY playwright.config.ts /app/

RUN apt-get update && apt-get -y install libnss3 libatk-bridge2.0-0 libdrm-dev libxkbcommon-dev libgbm-dev libasound-dev libatspi2.0-0 libxshmfence-dev

RUN npm i
