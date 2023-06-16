FROM node:slim

WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install

COPY . .

CMD [ "npm", "run", "with-grpc-tracer" ]
