import {parse, NamespaceBase, Service} from 'protobufjs';
import {IRpcValues, ITriggerService, TGRPCRequest} from 'types/Test.types';
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
    const parsedMetadata = metadata.filter(({key}) => key);

    return {
      address,
      request,
      auth,
      method,
      metadata: parsedMetadata,
      protobufFile,
    };
  },

  getInitialValues(request) {
    const {address: url, method, metadata, request: message, auth, protobufFile} = request as TGRPCRequest;
    const protoFile = new File([protobufFile], 'file.proto');

    return {
      url,
      auth,
      method,
      message,
      metadata,
      protoFile,
    };
  },
});

export default RpcTriggerService();
