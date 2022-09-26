FROM cypress/included:10.7.0

WORKDIR /app
COPY ./web/package.json ./
COPY ./web/package-lock.json ./
RUN npm ci --silent
COPY ./web ./

ENV CYPRESS_BASE_URL http://tracetest:11633

CMD "./node_modules/.bin/cypress run --config baseUrl=${CYPRESS_BASE_URL}"

