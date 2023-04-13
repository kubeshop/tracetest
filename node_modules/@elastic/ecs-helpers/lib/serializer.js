// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

'use strict'

const build = require('fast-json-stringify')

const string = { type: 'string' }
const number = { type: 'number' }

const stringify = build({
  type: 'object',
  properties: {
    '@timestamp': string,
    'log.level': string,
    log: {
      type: 'object',
      properties: {
        logger: string
      }
    },
    message: string,
    ecs: {
      type: 'object',
      properties: {
        version: string
      }
    },
    // https://www.elastic.co/guide/en/ecs/current/ecs-event.html
    event: {
      type: 'object',
      properties: {
        dataset: string,
        id: string
      },
      additionalProperties: true
    },
    // https://www.elastic.co/guide/en/ecs/current/ecs-http.html
    http: {
      type: 'object',
      properties: {
        version: string,
        request: {
          type: 'object',
          properties: {
            method: string,
            headers: {
              type: 'object',
              additionalProperties: true
            },
            body: {
              type: 'object',
              properties: {
                bytes: number
              }
            }
          }
        },
        response: {
          type: 'object',
          properties: {
            status_code: number,
            headers: {
              type: 'object',
              additionalProperties: true
            },
            body: {
              type: 'object',
              properties: {
                bytes: number
              }
            }
          }
        }
      }
    },
    // https://www.elastic.co/guide/en/ecs/current/ecs-url.html
    url: {
      type: 'object',
      properties: {
        path: string,
        domain: string,
        port: number,
        query: string,
        fragment: string,
        full: string
      }
    },
    // https://www.elastic.co/guide/en/ecs/current/ecs-client.html
    client: {
      type: 'object',
      properties: {
        address: string,
        ip: string,
        port: number
      }
    },
    // https://www.elastic.co/guide/en/ecs/current/ecs-user_agent.html
    user_agent: {
      type: 'object',
      properties: {
        original: string
      }
    },
    // https://www.elastic.co/guide/en/ecs/current/ecs-tracing.html
    trace: {
      type: 'object',
      properties: {
        id: string
      }
    },
    transaction: {
      type: 'object',
      properties: {
        id: string
      }
    },
    span: {
      type: 'object',
      properties: {
        id: string
      }
    },
    // https://www.elastic.co/guide/en/ecs/current/ecs-service.html
    service: {
      type: 'object',
      properties: {
        name: string
      },
      additionalProperties: true
    },
    // https://www.elastic.co/guide/en/ecs/current/ecs-error.html
    error: {
      type: 'object',
      properties: {
        type: string,
        message: string,
        stack_trace: string
      },
      additionalProperties: true
    }
  },
  additionalProperties: true
})

module.exports = stringify
