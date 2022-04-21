// JSON sructure to be used when auto generating selectors for a span.
// Semantic groups are named based on the filenames of each group
// as shown at https://github.com/open-telemetry/opentelemetry-specification/tree/main/semantic_conventions/trace
//
// All attributes listed below for a particular group should be checked for existence in
// the selected span that we want to autogenerate the selectors array for.  Create a SelectorItem
// for each attribute that you find a value for in that selected span.

import {LOCATION_NAME} from '../types';

// Note - need to add the following:
// lamda - aws lambda section
// aws-sdk - aws sdk section

export enum SemanticGroupNames {
  Http = 'http',
  Rpc = 'rpc',
  Messaging = 'messaging',
  Faas = 'faas',
  Exception = 'exception',
  General = 'general',
  Compatibility = 'compatibility',
  Database = 'db',
}

export const SemanticGroupNamesToText = {
  [SemanticGroupNames.Database]: 'Database',
  [SemanticGroupNames.Http]: 'HTTP',
  [SemanticGroupNames.Rpc]: 'RPC',
  [SemanticGroupNames.Messaging]: 'Messaging',
  [SemanticGroupNames.Faas]: 'FASS',
  [SemanticGroupNames.Exception]: 'Exception',
  [SemanticGroupNames.General]: 'General',
  [SemanticGroupNames.Compatibility]: 'Compatibility',
};

export const ResourceSpanAttributeList = ['service.name'];

export const SemanticGroupNameNodeMap: Record<SemanticGroupNames, {primary: string[]; type: string}> = {
  [SemanticGroupNames.Http]: {
    primary: ['http.target', 'service.name'],
    type: '',
  },
  [SemanticGroupNames.Database]: {
    primary: ['db.mongodb.collection', 'db.sql.table', 'db.redis.database_index', 'db.cassandra.table', 'service.name'],
    type: 'db.system',
  },
  [SemanticGroupNames.Rpc]: {
    primary: ['rpc.service', 'service.name'],
    type: '',
  },
  [SemanticGroupNames.Messaging]: {
    primary: ['messaging.destination', 'service.name'],
    type: 'messaging.system',
  },
  [SemanticGroupNames.Faas]: {
    primary: ['faas.invoked_name', 'service.name'],
    type: '',
  },
  [SemanticGroupNames.Exception]: {
    primary: ['exception.type', 'service.name'],
    type: '',
  },
  [SemanticGroupNames.General]: {
    primary: ['service.name', 'service.name'],
    type: '',
  },
  [SemanticGroupNames.Compatibility]: {
    primary: ['opentracing.ref_type', 'service.name'],
    type: '',
  },
};

export const SELECTOR_DEFAULT_ATTRIBUTES = [
  {
    semanticGroup: SemanticGroupNames.Http,
    attributes: ['service.name', 'http.target', 'http.method'],
  },
  {
    semanticGroup: SemanticGroupNames.Database,
    attributes: [
      'service.name',
      'db.system',
      'db.name',
      'db.user',
      'db.operation',
      'db.redis.database_index',
      'db.mongodb.collection',
      'db.sql.table',
      'db.cassandra.table',
    ],
  },
  {
    semanticGroup: SemanticGroupNames.Rpc,
    attributes: ['service.name', 'rpc.system', 'rpc.service', 'rpc.method', 'message.type'],
  },
  {
    semanticGroup: SemanticGroupNames.Messaging,
    attributes: [
      'service.name',
      'messaging.system',
      'messaging.destination',
      'messaging.destination_kind',
      'messaging.operation',
      'messaging.rabbitmq.routing_key',
      'messaging.kafka.consumer_group',
      'messaging.rocketmq.namespace',
      'messaging.rocketmq.client_group',
      'messaging.rocketmq.message_type',
      'messaging.rocketmq.message_keys',
      'messaging.rocketmq.consumption_model',
    ],
  },
  {
    semanticGroup: SemanticGroupNames.Faas,
    attributes: ['service.name', 'faas.invoked_name', 'faas.invoked_provider', 'faas.trigger'],
  },
  {
    semanticGroup: SemanticGroupNames.Exception,
    attributes: ['service.name', 'exception.type', 'exception.message', 'exception.escaped'],
  },
  {
    semanticGroup: SemanticGroupNames.Compatibility,
    attributes: ['service.name', 'opentracing.ref_type'],
  },
  {
    semanticGroup: SemanticGroupNames.General,
    attributes: [
      'service.name',
      'enduser.id',
      'enduser.role',
      'enduser.scope',
      'thread.name',
      'code.function',
      'code.namespace',
      'code.filepath',
    ],
  },
];

type TSemanticGroupSignature = {
  [key in SemanticGroupNames]: {[key2 in LOCATION_NAME]: string[]};
};

export const SemanticGroupsSignature = SELECTOR_DEFAULT_ATTRIBUTES.reduce<TSemanticGroupSignature>(
  (acc, {semanticGroup, attributes}) => ({
    ...acc,
    [semanticGroup]: {
      [LOCATION_NAME.SPAN_ATTRIBUTES]: attributes.filter(
        attribute => !ResourceSpanAttributeList.find(resourceAttribute => resourceAttribute === attribute)
      ),
      [LOCATION_NAME.RESOURCE_ATTRIBUTES]: attributes.filter(attribute =>
        ResourceSpanAttributeList.find(resourceAttribute => resourceAttribute === attribute)
      ),
    },
  }),
  {} as TSemanticGroupSignature
);
