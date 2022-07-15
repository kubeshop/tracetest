import {parse, NamespaceBase, Service} from 'protobufjs';
import {IRpcValues, ITriggerService} from 'types/Test.types';
import Validator from 'utils/Validator';

interface IRpcTriggerService extends ITriggerService {
  getMethodList(protoFile: string): string[];
}

const RpcTriggerService = (): IRpcTriggerService => ({
  getMethodList(protoFile) {
    const parsedData = parse(protoFile);

    const methodList = parsedData.root.nestedArray.flatMap(a => {
      const namespace = a as NamespaceBase;

      return namespace.nestedArray.flatMap(b => {
        const service = b as Service;
        return service.methods ? service.methodsArray : [];
      });
    });

    return methodList.reduce<string[]>(
      (list, {requestStream, responseStream, fullName}) =>
        !requestStream && !responseStream ? list.concat(fullName.slice(1, fullName.length)) : list,
      []
    );
  },
  async validateDraft(draft) {
    const {protoFile, method, url} = draft as IRpcValues;

    const isValid = Validator.required(url) && Validator.required(method) && Validator.required(protoFile);

    return isValid;
  },
  async getRequest(values) {
    const {protoFile, message: request, metadata, method, auth, url: address} = values as IRpcValues;
    const protobufFile = await protoFile.text();

    return {
      address,
      request,
      auth,
      method,
      metadata,
      protobufFile,
    };
  },
});

export default RpcTriggerService();
