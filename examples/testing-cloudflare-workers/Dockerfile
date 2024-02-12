FROM node
WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm i

COPY src ./src
COPY test ./test
COPY wrangler.toml .
COPY tsconfig.json .
COPY schema.sql .

RUN npx wrangler d1 execute testing-cloudflare-workers --local --file=./schema.sql

CMD ["npm", "run", "docker"]
