#Stage 1
FROM node:20-alpine as builder
WORKDIR /app
COPY package*.json .
RUN npm install
COPY . .

ARG REACT_APP_APP_HOST
ENV REACT_APP_APP_HOST=$REACT_APP_APP_HOST
ARG REACT_APP_OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
ENV REACT_APP_OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=$REACT_APP_OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
RUN npm run build

#Stage 2
FROM nginx:1.19.0
WORKDIR /usr/share/nginx/html
RUN rm -rf ./*
COPY --from=builder /app/build .
ENTRYPOINT ["nginx", "-g", "daemon off;"]
