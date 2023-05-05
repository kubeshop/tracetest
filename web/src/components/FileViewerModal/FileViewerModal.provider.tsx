import {capitalize, noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import FileViewerModal from 'components/FileViewerModal';
import {
  useLazyGetJUnitByRunIdQuery,
  useLazyGetResourceDefinitionQuery,
  useLazyGetResourceDefinitionV2Query,
} from 'redux/apis/TraceTest.api';
import {ResourceType} from 'types/Resource.type';

interface IContext {
  loadJUnit(testId: string, runId: string): void;
  loadDefinition(resourceType: ResourceType, resourceId: string, version?: number): void;
}

export const Context = createContext<IContext>({
  loadJUnit: noop,
  loadDefinition: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useFileViewerModal = () => useContext(Context);

const propsMap = {
  definition: {
    title: 'Test Definition',
    language: 'yaml',
    subtitle: 'Preview your YAML file',
    fileName: 'test-definition.yaml',
  },
  junit: {
    title: 'JUnit Results',
    language: 'xml',
    subtitle: 'Preview your JUnit results',
    fileName: 'junit.xml',
  },
};

const FileViewerModalProvider = ({children}: IProps) => {
  const [isFileViewerOpen, setIsFileViewerOpen] = useState(false);
  const [fileViewerData, setFileViewerData] = useState<{data: string; type: 'definition' | 'junit'}>({
    data: '',
    type: 'definition',
  });
  const [getJUnit] = useLazyGetJUnitByRunIdQuery();
  const [getResourceDefinition] = useLazyGetResourceDefinitionQuery();
  const [getResourceDefinitionV2] = useLazyGetResourceDefinitionV2Query();
  const [fileProps, setProps] = useState({
    title: '',
    language: '',
    subtitle: '',
    fileName: '',
  });

  const loadJUnit = useCallback(
    async (testId: string, runId: string) => {
      const data = await getJUnit({runId, testId}).unwrap();
      setIsFileViewerOpen(true);
      setFileViewerData({data, type: 'junit'});
      setProps(propsMap.junit);
    },
    [getJUnit]
  );

  const loadDefinition = useCallback(
    async (resourceType: ResourceType, resourceId: string, version?: number) => {
      const data = await (resourceType === ResourceType.Environment || resourceType === ResourceType.Transaction
        ? getResourceDefinitionV2({resourceId, resourceType}).unwrap()
        : getResourceDefinition({resourceId, version, resourceType}).unwrap());
      setIsFileViewerOpen(true);
      setFileViewerData({data, type: 'definition'});
      setProps({
        title: `${capitalize(resourceType)} Definition`,
        language: 'yaml',
        subtitle: 'Preview your YAML file',
        fileName: `${resourceType}-${resourceId}-${version || 0}-definition.yaml`,
      });
    },
    [getResourceDefinition, getResourceDefinitionV2]
  );

  const value: IContext = useMemo(() => ({loadJUnit, loadDefinition}), [loadJUnit, loadDefinition]);

  return (
    <>
      <Context.Provider value={value}>{children}</Context.Provider>
      <FileViewerModal
        isOpen={isFileViewerOpen}
        {...fileProps}
        data={fileViewerData.data}
        onClose={() => setIsFileViewerOpen(false)}
      />
    </>
  );
};

export default FileViewerModalProvider;
