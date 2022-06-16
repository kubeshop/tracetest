import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import FileViewerModal from 'components/FileViewerModal';
import {useLazyGetJUnitByRunIdQuery} from 'redux/apis/TraceTest.api';

interface IContext {
  loadJUnit(testId: string, runId: string): void;
  loadTestDefinitionYaml(testId: string, version: number): void;
}

export const Context = createContext<IContext>({
  loadJUnit: noop,
  loadTestDefinitionYaml: noop,
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

  const loadJUnit = useCallback(
    async (testId: string, runId: string) => {
      const data = await getJUnit({runId, testId}).unwrap();
      setIsFileViewerOpen(true);
      setFileViewerData({data, type: 'junit'});
    },
    [getJUnit]
  );

  const loadTestDefinitionYaml = useCallback((testId: string, version: number) => {
    setIsFileViewerOpen(true);
    setFileViewerData({data: `name: ${testId}-${version}`, type: 'definition'});
  }, []);

  const value: IContext = useMemo(() => ({loadJUnit, loadTestDefinitionYaml}), [loadJUnit, loadTestDefinitionYaml]);
  const fileProps = propsMap[fileViewerData.type];

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
