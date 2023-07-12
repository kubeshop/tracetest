import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import MissingVariablesModal from 'components/MissingVariablesModal';
import VariablesService from 'services/Variables.service';
import MissingVariables from 'models/MissingVariables.model';
import {TEnvironmentValue} from 'models/Environment.model';
import Test from 'models/Test.model';

type TOnOPenProps = {
  missingVariables: MissingVariables;
  name: string;
  onSubmit(draft: TEnvironmentValue[]): void;
  onCancel?(): void;
  testList: Test[];
};

interface IContext {
  onOpen(props: TOnOPenProps): void;
}

export const Context = createContext<IContext>({
  onOpen: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useMissingVariablesModal = () => useContext(Context);

const MissingVariablesModalProvider = ({children}: IProps) => {
  const [{missingVariables = [], testList = [], onSubmit, onCancel = noop, name}, setProps] = useState<TOnOPenProps>({
    missingVariables: [],
    onSubmit: noop,
    onCancel: noop,
    name: '',
    testList: [],
  });
  const [isOpen, setIsOpen] = useState(false);

  const onOpen = useCallback((newProps: TOnOPenProps) => {
    setProps(newProps);
    setIsOpen(true);
  }, []);

  const handleSubmit = useCallback(
    (draft: TEnvironmentValue[]) => {
      onSubmit(draft);
      setIsOpen(false);
    },
    [onSubmit]
  );

  const value = useMemo<IContext>(() => ({onOpen}), [onOpen]);

  const testVariables = useMemo(
    () => VariablesService.getVariableEntries(missingVariables, testList),
    [missingVariables, testList]
  );

  return (
    <Context.Provider value={value}>
      {children}
      <MissingVariablesModal
        testVariables={testVariables}
        onClose={() => {
          setIsOpen(false);
          onCancel();
        }}
        onSubmit={handleSubmit}
        isOpen={isOpen}
        name={name}
      />
    </Context.Provider>
  );
};

export default MissingVariablesModalProvider;
