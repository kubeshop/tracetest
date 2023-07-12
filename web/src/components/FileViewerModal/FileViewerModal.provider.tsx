import {capitalize, noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import FileViewerModal from 'components/FileViewerModal';
import {ResourceType} from 'types/Resource.type';
import useDefinitionFile from 'hooks/useDefinitionFile';
import useJUnitResult from 'hooks/useJUnitResult';

interface IContext {
  onJUnit(testId: string, runId: string): void;
  onDefinition(resourceType: ResourceType, resourceId: string, version?: number): void;
}

export const Context = createContext<IContext>({
  onJUnit: noop,
  onDefinition: noop,
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
  const {definition, loadDefinition} = useDefinitionFile();
  const {jUnit, loadJUnit} = useJUnitResult();
  const [fileProps, setProps] = useState({
    title: '',
    language: '',
    subtitle: '',
    fileName: '',
  });

  const onJUnit = useCallback(
    async (testId: string, runId: string) => {
      loadJUnit(testId, runId);
      setIsFileViewerOpen(true);
      setProps(propsMap.junit);
    },
    [loadJUnit]
  );

  const onDefinition = useCallback(
    async (resourceType: ResourceType, resourceId: string, version?: number) => {
      setIsFileViewerOpen(true);
      loadDefinition(resourceType, resourceId, version);
      setProps({
        title: `${capitalize(resourceType)} Definition`,
        language: 'yaml',
        subtitle: 'Preview your YAML file',
        fileName: `${resourceType}-${resourceId}-${version || 0}-definition.yaml`,
      });
    },
    [loadDefinition]
  );

  const value: IContext = useMemo(() => ({onJUnit, onDefinition}), [onJUnit, onDefinition]);

  return (
    <>
      <Context.Provider value={value}>{children}</Context.Provider>
      <FileViewerModal
        isOpen={isFileViewerOpen}
        {...fileProps}
        data={definition || jUnit}
        onClose={() => setIsFileViewerOpen(false)}
      />
    </>
  );
};

export default FileViewerModalProvider;
