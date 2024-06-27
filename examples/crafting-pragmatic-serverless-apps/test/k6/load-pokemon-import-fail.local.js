import http from 'k6/http'
import { sleep } from 'k6'

export const options = {
  vus: 10,
  duration: '5s',
}

export default function() {
  http.post('http://localhost:8787/api/pokemon?id=99999')
  sleep(1)
}
