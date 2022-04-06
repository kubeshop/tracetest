// JSON sructure to be used when auto generating selectors for a span.
// Semantic groups are named based on the filenames of each group
// as shown at https://github.com/open-telemetry/opentelemetry-specification/tree/main/semantic_conventions/trace
//
// All attributes listed below for a particular group should be checked for existence in
// the selected span that we want to autogenerate the selectors array for.  Create a SelectorItem
// for each attribute that you find a value for in that selected span.

// Note - need to add the following:
// lamda - aws lambda section
// aws-sdk - aws sdk section

export enum SemanticGroupNames {
  Http = 'http',
  Rcp = 'rcp',
  Messaging = 'messaging',
  Faas = 'faas',
  Exception = 'exception',
  General = 'general',
  Compatibility = 'compatibility',
}

export const SemanticGroupNamesToText = {
  [SemanticGroupNames.Http]: 'HTTP Service',
  [SemanticGroupNames.Rcp]: 'RCP Service',
  [SemanticGroupNames.Messaging]: 'Messaging Queue',
  [SemanticGroupNames.Faas]: 'Function as a Service',
  [SemanticGroupNames.Exception]: 'Exception',
  [SemanticGroupNames.General]: 'General',
  [SemanticGroupNames.Compatibility]: 'Compatibility',
};

export const SELECTOR_DEFAULT_ATTRIBUTES = [
  {
    semanticGroup: SemanticGroupNames.Http,
    attributes: ['service.name', 'http.target', 'http.method'],
  },
  {
    semanticGroup: 'database',
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
    semanticGroup: SemanticGroupNames.Rcp,
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
    attributes: ['service.name', 'faas.invoked_name', 'faas.invoked_provider', 'faas.trigger', 'faas.trigger'],
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
