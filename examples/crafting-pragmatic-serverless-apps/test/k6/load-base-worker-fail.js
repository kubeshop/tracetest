import http from 'k6/http'
import { sleep } from 'k6'

export const options = {
  vus: 100,
  duration: '5s',
}

export default function() {
  http.get('https://crafting-pragmatic-serverless-apps.xxxx.workers.dev/fail')
  sleep(1)
}
