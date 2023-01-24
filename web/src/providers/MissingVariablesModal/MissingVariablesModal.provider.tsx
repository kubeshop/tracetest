import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import MissingVariablesModal from 'components/MissingVariablesModal';
import {TEnvironmentValue} from 'types/Environment.types';
import {TMissingVariable} from 'types/Variables.types';
import VariablesService from 'services/Variables.service';
import {TTest} from 'types/Test.types';

type TOnOPenProps = {
  missingVariables: TMissingVariable[];
  name: string;
  onSubmit(draft: TEnvironmentValue[]): void;
  testList: TTest[];
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
  const [{missingVariables = [], testList = [], onSubmit, name}, setProps] = useState<TOnOPenProps>({
    missingVariables: [],
    onSubmit: noop,
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
        onClose={() => setIsOpen(false)}
        onSubmit={handleSubmit}
        isOpen={isOpen}
        name={name}
      />
    </Context.Provider>
  );
};

export default MissingVariablesModalProvider;
