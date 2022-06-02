import {Attributes} from './SpanAttribute.constants';

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

export const SemanticGroupNameNodeMap: Record<SemanticGroupNames, {primary: string[]; type: string}> = {
  [SemanticGroupNames.Http]: {
    primary: [Attributes.HTTP_TARGET, Attributes.SERVICE_NAME],
    type: '',
  },
  [SemanticGroupNames.Database]: {
    primary: [
      Attributes.DB_MONGODB_COLLECTION,
      Attributes.DB_SQL_TABLE,
      Attributes.DB_CASSANDRA_TABLE,
      Attributes.SERVICE_NAME,
    ],
    type: Attributes.DB_SYSTEM,
  },
  [SemanticGroupNames.Rpc]: {
    primary: [Attributes.RPC_SYSTEM, Attributes.SERVICE_NAME],
    type: Attributes.RPC_SYSTEM,
  },
  [SemanticGroupNames.Messaging]: {
    primary: [Attributes.MESSAGING_DESTINATION, Attributes.SERVICE_NAME],
    type: Attributes.MESSAGING_SYSTEM,
  },
  [SemanticGroupNames.Faas]: {
    primary: [Attributes.NAME, Attributes.SERVICE_NAME],
    type: '',
  },
  [SemanticGroupNames.Exception]: {
    primary: [Attributes.EXCEPTION_TYPE, Attributes.SERVICE_NAME],
    type: '',
  },
  [SemanticGroupNames.General]: {
    primary: [Attributes.NAME, Attributes.SERVICE_NAME],
    type: '',
  },
  [SemanticGroupNames.Compatibility]: {
    primary: [Attributes.NAME, Attributes.SERVICE_NAME],
    type: '',
  },
};

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
