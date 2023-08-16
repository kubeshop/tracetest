import {Model, TKafkaSchemas} from '../types/Common.types';
import {TRawKafkaMessageHeader} from '../types/Test.types';

export type TRawKafkaRequest = TKafkaSchemas['KafkaRequest'];

type KafkaRequest = Model<
  TRawKafkaRequest,
  {
    headers: Model<TRawKafkaMessageHeader, {}>[];
  }
>;

const KafkaRequest = ({
  brokerUrls = [],
  topic = '',
  authentication = {},
  sslVerification = false,
  headers = [],
  messageKey = '',
  messageValue = '',
}: TRawKafkaRequest): KafkaRequest => {
  if (!headers) {
    headers = []; // guard clause to avoid null header
  }

  return {
    brokerUrls,
    topic,
    authentication,
    sslVerification,
    messageKey,
    messageValue,
    headers: headers.map(({key = '', value = ''}) => ({
      key,
      value,
    })),
  };
};

export default KafkaRequest;
