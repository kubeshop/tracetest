import {useCallback, useState} from 'react';
import {useLazyGetResourceDefinitionQuery, useLazyGetResourceDefinitionV2Query} from 'redux/apis/TraceTest.api';
import {ResourceType} from 'types/Resource.type';

const useDefinitionFile = () => {
  const [definition, setDefinition] = useState<string>('');
  const [getResourceDefinition] = useLazyGetResourceDefinitionQuery();
  const [getResourceDefinitionV2] = useLazyGetResourceDefinitionV2Query();

  const loadDefinition = useCallback(
    async (resourceType: ResourceType, resourceId: string, version?: number) => {
      const data = await (resourceType === ResourceType.Environment || resourceType === ResourceType.Transaction
        ? getResourceDefinitionV2({resourceId, resourceType}).unwrap()
        : getResourceDefinition({resourceId, version, resourceType}).unwrap());

      setDefinition(data);
    },
    [getResourceDefinition, getResourceDefinitionV2]
  );

  return {definition, loadDefinition};
};

export default useDefinitionFile;
