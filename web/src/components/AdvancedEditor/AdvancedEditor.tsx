import {useMemo} from 'react';
import {noop} from 'lodash';
import CodeMirror from '@uiw/react-codemirror';
import {autocompletion} from '@codemirror/autocomplete';
import {tracetest} from 'utils/grammar';
import {linter} from '@codemirror/lint';
import useAutoComplete from './hooks/useAutoComplete';
import useLint from './hooks/useLint';

interface IProps {
  testId: string;
  runId: string;
  value?: string;
  onChange?(value: string): void;
}

const AdvancedEditor = ({testId, runId, onChange = noop, value = ''}: IProps) => {
  const completionFn = useAutoComplete({testId, runId});
  const lintFn = useLint({testId, runId});

  const extensionList = useMemo(
    () => [autocompletion({override: [completionFn]}), linter(lintFn), tracetest()],
    [completionFn, lintFn]
  );

  return (
    <CodeMirror
      data-cy="advanced-selector"
      value={value}
      maxHeight="100px"
      extensions={extensionList}
      onChange={onChange}
      spellCheck={false}
    />
  );
};

export default AdvancedEditor;
