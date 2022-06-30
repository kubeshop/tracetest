import {uniq} from 'lodash';
import {SemanticGroupNames, BASE_ATTRIBUTES, SemanticGroupsSignature} from './SemanticGroupNames.constants';
import {Attributes, SemanticAttributes} from './SpanAttribute.constants';

type TValueOfAttributes = typeof Attributes[keyof typeof Attributes];

export enum SectionNames {
  request = 'request',
  response = 'response',
  metadata = 'metadata',
  operation = 'operation',
  consumer = 'consumer',
  producer = 'producer',
  custom = 'custom',
  all = 'all',
}

export const SelectorAttributesBlackList = [SemanticAttributes.DB_STATEMENT, SemanticAttributes.DB_CONNECTION_STRING];

export const SpanAttributeSections: Record<SemanticGroupNames, Record<string, TValueOfAttributes[]>> = {
  [SemanticGroupNames.Http]: {
    [SectionNames.request]: [
      Attributes.HTTP_URL,
      Attributes.HTTP_METHOD,
      Attributes.HTTP_ROUTE,
      Attributes.HTTP_TARGET,
      Attributes.HTTP_FLAVOR,
      Attributes.HTTP_HOST,
      Attributes.HTTP_CLIENT_IP,
      Attributes.HTTP_SCHEME,
      Attributes.HTTP_REQUEST_CONTENT_LENGTH,
      Attributes.HTTP_REQUEST_CONTENT_LENGTH_UNCOMPRESSED,
      Attributes.HTTP_USER_AGENT,
      Attributes.HTTP_REQUEST_HEADER,
    ],
    [SectionNames.response]: [
      Attributes.HTTP_STATUS_CODE,
      Attributes.TRACETEST_RESPONSE_BODY,
      Attributes.TRACETEST_RESPONSE_HEADERS,
      Attributes.HTTP_RESPONSE_CONTENT_LENGTH,
      Attributes.HTTP_RESPONSE_CONTENT_LENGTH_UNCOMPRESSED,
      Attributes.HTTP_RESPONSE_HEADER,
    ],
  },
  [SemanticGroupNames.Database]: {
    [SectionNames.metadata]: [
      Attributes.DB_NAME,
      Attributes.DB_SYSTEM,
      Attributes.DB_USER,
      Attributes.DB_CONNECTION_STRING,
      Attributes.DB_MONGODB_COLLECTION,
      Attributes.DB_REDIS_DATABASE_INDEX,
      Attributes.DB_SQL_TABLE,
      Attributes.DB_CASSANDRA_TABLE,
      Attributes.DB_CASSANDRA_PAGE_SIZE,
      Attributes.DB_CASSANDRA_CONSISTENCY_LEVEL,
      Attributes.DB_CASSANDRA_IDEMPOTENCE,
      Attributes.DB_CASSANDRA_SPECULATIVE_EXECUTION_COUNT,
      Attributes.DB_CASSANDRA_COORDINATOR_ID,
      Attributes.DB_CASSANDRA_COORDINATOR_DC,
      Attributes.DB_OPERATION,
      Attributes.DB_STATEMENT,
    ],
  },
  [SemanticGroupNames.Messaging]: {
    [SectionNames.metadata]: [
      Attributes.MESSAGING_SYSTEM,
      Attributes.MESSAGING_URL,
      Attributes.MESSAGING_PROTOCOL,
      Attributes.MESSAGING_RABBITMQ_ROUTING_KEY,
      Attributes.MESSAGING_KAFKA_CLIENT_ID,
      Attributes.MESSAGING_KAFKA_PARTITION,
      Attributes.MESSAGING_KAFKA_TOMBSTONE,
      Attributes.MESSAGING_DESTINATION,
      Attributes.MESSAGING_DESTINATION_KIND,
      Attributes.MESSAGING_TEMP_DESTINATION,
      Attributes.MESSAGING_CONVERSATION_ID,
      Attributes.MESSAGING_RABBITMQ_ROUTING_KEY,
      Attributes.MESSAGING_KAFKA_MESSAGE_KEY,
      Attributes.MESSAGING_OPERATION,
      Attributes.MESSAGING_CONSUMER_ID,
      Attributes.MESSAGING_KAFKA_CONSUMER_GROUP,
    ],
  },
  [SemanticGroupNames.Rpc]: {
    [SectionNames.metadata]: SemanticGroupsSignature.rpc,
  },
  [SemanticGroupNames.Exception]: {
    [SectionNames.metadata]: SemanticGroupsSignature.exception,
  },
  [SemanticGroupNames.General]: {
    [SectionNames.metadata]: SemanticGroupsSignature.general,
  },
  [SemanticGroupNames.Compatibility]: {
    [SectionNames.metadata]: SemanticGroupsSignature.compatibility,
  },
  [SemanticGroupNames.Faas]: {
    [SectionNames.metadata]: SemanticGroupsSignature.faas,
  },
};

export const SelectorAttributesWhiteList = uniq([
  ...BASE_ATTRIBUTES,
  ...Object.values(SpanAttributeSections).flatMap(section => {
    const sectionAttributes = Object.values(section).flatMap(attributes => attributes);

    return sectionAttributes;
  }),
]);
