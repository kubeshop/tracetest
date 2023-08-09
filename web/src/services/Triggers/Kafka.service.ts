import {IKafkaValues, ITriggerService} from 'types/Test.types';
import Validator from 'utils/Validator';
import KafkaRequest from 'models/KafkaRequest.model';

const KafkaTriggerService = (): ITriggerService => ({
  async validateDraft(draft) {
    const { brokerUrls, topic, messageValue } = draft as IKafkaValues;

    const isValid = Validator.required(brokerUrls) && Validator.required(topic) && Validator.required(messageValue);

    return isValid;
  },
  async getRequest(values) {
    const { brokerUrls, topic, authentication, sslVerification, headers, messageKey, messageValue } = values as IKafkaValues;
    const parsedHeaders = headers.filter(({key}) => key);

    return KafkaRequest({
      brokerUrls,
      topic,
      authentication,
      sslVerification,
      headers: parsedHeaders,
      messageKey,
      messageValue
    });
  },

  getInitialValues(request) {
    console.log(request);
    const { brokerUrls, topic, authentication, sslVerification, headers, messageKey, messageValue } = request as KafkaRequest;

    return {
      brokerUrls,
      topic,
      authentication,
      sslVerification,
      headers,
      messageKey,
      messageValue
    };
  },
});

export default KafkaTriggerService();
