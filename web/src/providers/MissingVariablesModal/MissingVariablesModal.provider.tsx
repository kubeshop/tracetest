import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {TTestVariables} from 'types/Variables.types';
import VariablesService from 'services/Variables.service';
import MissingVariablesModal from 'components/MissingVariablesModal';
import {TEnvironmentValue} from 'types/Environment.types';

type TOnOPenProps = {
  testsVariables: TTestVariables[];
  name: string;
  onSubmit(draft: TEnvironmentValue[]): void;
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
  const [{testsVariables = [], onSubmit, name}, setProps] = useState<TOnOPenProps>({
    testsVariables: [],
    onSubmit: noop,
    name: '',
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

  const variables = useMemo(() => VariablesService.getVariableEntries(testsVariables), [testsVariables]);
  const value = useMemo<IContext>(() => ({onOpen}), [onOpen]);

  return (
    <Context.Provider value={value}>
      {children}
      <MissingVariablesModal
        variables={variables}
        onClose={() => setIsOpen(false)}
        onSubmit={handleSubmit}
        isOpen={isOpen}
        name={name}
      />
    </Context.Provider>
  );
};

export default MissingVariablesModalProvider;
