import {parse, NamespaceBase, Service} from 'protobufjs';

const RpcService = () => ({
  getMethodList: async (protoFile: File): Promise<string[]> => {
    const data = await protoFile.text();
    const parsedData = parse(data);

    const methodList = parsedData.root.nestedArray.flatMap(a =>
      (a as NamespaceBase).nestedArray.flatMap(b => {
        const service = b as Service;
        return service.methods ? service.methodsArray : [];
      })
    );

    return methodList.reduce<string[]>(
      (list, {requestStream, responseStream, name}) => (!requestStream && !responseStream ? list.concat(name) : list),
      []
    );
  },
});

export default RpcService();
