import http from 'k6/http'
import { sleep } from 'k6'

export const options = {
  vus: 10,
  duration: '5s',
}

export default function() {
  // or
  http.get('http://localhost:8787/')
  sleep(1)
}
