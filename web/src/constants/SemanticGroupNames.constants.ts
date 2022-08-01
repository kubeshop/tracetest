import {Attributes} from './SpanAttribute.constants';

export enum SemanticGroupNames {
  General = 'general',
  Http = 'http',
  Database = 'database',
  Rpc = 'rpc',
  Messaging = 'messaging',
  Faas = 'faas',
  Exception = 'exception',
  Compatibility = 'compatibility',
}

export const SemanticGroupNamesToText = {
  [SemanticGroupNames.General]: 'General',
  [SemanticGroupNames.Http]: 'HTTP',
  [SemanticGroupNames.Database]: 'Database',
  [SemanticGroupNames.Rpc]: 'RPC',
  [SemanticGroupNames.Messaging]: 'Messaging',
  [SemanticGroupNames.Faas]: 'FaaS',
  [SemanticGroupNames.Exception]: 'Exception',
  [SemanticGroupNames.Compatibility]: 'Compatibility',
} as const;

export const SemanticGroupNamesToColor = {
  [SemanticGroupNames.General]: '#FFBB96',
  [SemanticGroupNames.Http]: '#C1E095',
  [SemanticGroupNames.Database]: '#EFDBFF',
  [SemanticGroupNames.Rpc]: '#87E8DE',
  [SemanticGroupNames.Messaging]: '#91D5FF',
  [SemanticGroupNames.Faas]: '#FFD591',
  [SemanticGroupNames.Exception]: '#FFFB8F',
  [SemanticGroupNames.Compatibility]: '#ADC6FF',
} as const;

export const SemanticGroupNamesToLightColor = {
  [SemanticGroupNames.General]: 'rgba(255, 187, 150, 0.3)',
  [SemanticGroupNames.Http]: 'rgba(193, 224, 149, 0.3)',
  [SemanticGroupNames.Database]: 'rgba(239, 219, 255, 0.3)',
  [SemanticGroupNames.Rpc]: 'rgba(135, 232, 222, 0.3)',
  [SemanticGroupNames.Messaging]: 'rgba(145, 213, 255, 0.3)',
  [SemanticGroupNames.Faas]: 'rgba(255, 213, 145, 0.3)',
  [SemanticGroupNames.Exception]: 'rgba(255, 251, 143, 0.3)',
  [SemanticGroupNames.Compatibility]: 'rgba(173, 198, 255, 0.3)',
} as const;

export const SemanticGroupNamesToSystem = {
  [SemanticGroupNames.General]: '',
  [SemanticGroupNames.Http]: '',
  [SemanticGroupNames.Database]: Attributes.DB_SYSTEM,
  [SemanticGroupNames.Rpc]: Attributes.RPC_SYSTEM,
  [SemanticGroupNames.Messaging]: Attributes.MESSAGING_SYSTEM,
  [SemanticGroupNames.Faas]: Attributes.CLOUD_PROVIDER,
  [SemanticGroupNames.Exception]: '',
  [SemanticGroupNames.Compatibility]: '',
} as const;

export const BASE_ATTRIBUTES = [Attributes.TRACETEST_SPAN_TYPE, Attributes.SERVICE_NAME, Attributes.NAME];

export const SELECTOR_DEFAULT_ATTRIBUTES = [
  {
    semanticGroup: SemanticGroupNames.Http,
    attributes: [...BASE_ATTRIBUTES, Attributes.HTTP_TARGET, Attributes.HTTP_METHOD],
  },
  {
    semanticGroup: SemanticGroupNames.Database,
    attributes: [
      ...BASE_ATTRIBUTES,
      Attributes.DB_SYSTEM,
      Attributes.DB_NAME,
      Attributes.DB_USER,
      Attributes.DB_OPERATION,
      Attributes.DB_MONGODB_COLLECTION,
      Attributes.DB_REDIS_DATABASE_INDEX,
      Attributes.DB_SQL_TABLE,
      Attributes.DB_CASSANDRA_TABLE,
    ],
  },
  {
    semanticGroup: SemanticGroupNames.Rpc,
    attributes: [
      ...BASE_ATTRIBUTES,
      Attributes.RPC_SYSTEM,
      Attributes.RPC_METHOD,
      Attributes.RPC_SERVICE,
      Attributes.MESSAGE_TYPE,
    ],
  },
  {
    semanticGroup: SemanticGroupNames.Messaging,
    attributes: [
      ...BASE_ATTRIBUTES,
      Attributes.MESSAGING_SYSTEM,
      Attributes.MESSAGING_DESTINATION,
      Attributes.MESSAGING_DESTINATION_KIND,
      Attributes.MESSAGING_OPERATION,
      Attributes.MESSAGING_RABBITMQ_ROUTING_KEY,
      Attributes.MESSAGING_KAFKA_CONSUMER_GROUP,
    ],
  },
  {
    semanticGroup: SemanticGroupNames.Faas,
    attributes: [
      ...BASE_ATTRIBUTES,
      Attributes.FAAS_INVOKED_NAME,
      Attributes.FAAS_INVOKED_PROVIDER,
      Attributes.FAAS_TRIGGER,
    ],
  },
  {
    semanticGroup: SemanticGroupNames.Exception,
    attributes: [
      ...BASE_ATTRIBUTES,
      Attributes.EXCEPTION_TYPE,
      Attributes.EXCEPTION_MESSAGE,
      Attributes.EXCEPTION_ESCAPED,
    ],
  },
  {
    semanticGroup: SemanticGroupNames.Compatibility,
    attributes: BASE_ATTRIBUTES,
  },
  {
    semanticGroup: SemanticGroupNames.General,
    attributes: [
      ...BASE_ATTRIBUTES,
      Attributes.ENDUSER_ID,
      Attributes.ENDUSER_ROLE,
      Attributes.ENDUSER_SCOPE,
      Attributes.THREAD_NAME,
      Attributes.CODE_FUNCTION,
      Attributes.CODE_NAMESPACE,
      Attributes.CODE_FILEPATH,
    ],
  },
];

type TSemanticGroupSignature = {
  [key in SemanticGroupNames]: string[];
};

export const SemanticGroupsSignature = SELECTOR_DEFAULT_ATTRIBUTES.reduce<TSemanticGroupSignature>(
  (acc, {semanticGroup, attributes}) => ({
    ...acc,
    [semanticGroup]: attributes,
  }),
  {} as TSemanticGroupSignature
);
