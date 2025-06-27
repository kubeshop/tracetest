import { registerOTel } from '@vercel/otel'
registerOTel({ serviceName: 'next-app' })
