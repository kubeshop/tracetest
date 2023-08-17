import {useCallback, useState} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';
import {ResourceType} from 'types/Resource.type';

const {useLazyGetResourceDefinitionQuery} = TracetestAPI.instance;

const useDefinitionFile = () => {
  const [definition, setDefinition] = useState<string>('');
  const [getResourceDefinition] = useLazyGetResourceDefinitionQuery();

  const loadDefinition = useCallback(
    async (resourceType: ResourceType, resourceId: string, version?: number) => {
      const data = await getResourceDefinition({resourceId, resourceType, version}).unwrap();

      setDefinition(data);
    },
    [getResourceDefinition]
  );

  return {definition, loadDefinition};
};

export default useDefinitionFile;
