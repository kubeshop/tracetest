import {useCallback} from 'react';
import {CompletionContext} from '@codemirror/autocomplete';
import {syntaxTree} from '@codemirror/language';
import {completeSourceAfter, parserList, Tokens} from 'constants/Editor.constants';

const {demoEndpoints = '{}'} = window.ENV || {};
const {PokeshopHttp = ''}: Record<string, string> = JSON.parse(demoEndpoints);

const environmentList = [
  {
    key: 'HOST',
    value: PokeshopHttp,
  },
  {
    key: 'PORT',
    value: '8080',
  },
];

const useAutoComplete = () => {
  return useCallback(async (context: CompletionContext) => {
    const {state, pos} = context;
    const tree = syntaxTree(state);
    const nodeBefore = tree.resolveInner(pos, -1);

    if (nodeBefore.name === Tokens.Pipe) {
      return {
        from: nodeBefore.to,
        options: parserList,
      };
    }

    if (completeSourceAfter.includes(nodeBefore.name)) {
      return {
        from: nodeBefore.to,
        options: environmentList.map(({key}) => ({
          label: key,
          type: 'variableName',
          apply: key,
        })),
      };
    }

    if (nodeBefore.prevSibling?.name === Tokens.Source) {
      return {
        from: nodeBefore.from,
        to: nodeBefore.to,
        options: environmentList.map(({key}) => ({
          label: key,
          type: 'variableName',
          apply: key,
        })),
      };
    }

    return null;
  }, []);
};

export default useAutoComplete;
