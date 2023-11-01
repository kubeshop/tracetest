import {parse, Namespace, ReflectionObject, Service} from 'protobufjs';
import {IRpcValues, ITriggerService} from 'types/Test.types';
import Validator from 'utils/Validator';
import GrpcRequest from 'models/GrpcRequest.model';

interface IRpcTriggerService extends ITriggerService {
  getMethodList(protoFile: string): string[];
}

function isService(ro: ReflectionObject): ro is Service {
  return (ro as Service).methods !== undefined;
}

const RpcTriggerService = (): IRpcTriggerService => ({
  getMethodList(protoFile) {
    const parsedData = parse(protoFile);

    const methodList = parsedData.root.nestedArray.flatMap(aReflection => {
      if (isService(aReflection)) {
        return aReflection.methodsArray;
      }

      return (
        (aReflection as Namespace)?.nestedArray?.flatMap(bReflection => {
          if (isService(bReflection)) {
            return bReflection.methodsArray;
          }

          return [];
        }) ?? []
      );
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

    return GrpcRequest({
      address,
      request,
      auth,
      method,
      metadata: parsedMetadata,
      protobufFile,
    });
  },

  getInitialValues(request) {
    const {address: url, method, metadata, request: message, auth, protobufFile} = request as GrpcRequest;
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
