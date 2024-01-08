FROM node:20-alpine AS base
FROM base AS builder
WORKDIR /app
RUN apk add --no-cache g++ make py3-pip

COPY package.json package-lock.json ./
RUN npm i

COPY app ./app
COPY pages ./pages
COPY public ./public
COPY next.config.js .
COPY tsconfig.json .
COPY instrumentation.ts .
COPY instrumentation.node.ts .

RUN npm run build

# Step 2. Production image, copy all the files and run next
FROM base AS runner
WORKDIR /app
# Don't run production as root
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
USER nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
CMD ["node", "server.js"]
