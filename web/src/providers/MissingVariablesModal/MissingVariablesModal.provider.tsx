import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import MissingVariablesModal from 'components/MissingVariablesModal';
import {TEnvironmentValue} from 'types/Environment.types';
import {TMissingVariable} from 'types/Variables.types';

type TOnOPenProps = {
  missingVariables: TMissingVariable[];
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
  const [{missingVariables = [], onSubmit, name}, setProps] = useState<TOnOPenProps>({
    missingVariables: [],
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

  const value = useMemo<IContext>(() => ({onOpen}), [onOpen]);

  return (
    <Context.Provider value={value}>
      {children}
      <MissingVariablesModal
        missingVariables={missingVariables}
        onClose={() => setIsOpen(false)}
        onSubmit={handleSubmit}
        isOpen={isOpen}
        name={name}
      />
    </Context.Provider>
  );
};

export default MissingVariablesModalProvider;
