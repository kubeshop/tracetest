import {useCallback} from 'react';
import {uniqBy} from 'lodash';
import {CompletionContext} from '@codemirror/autocomplete';
import {syntaxTree} from '@codemirror/language';
import {completeSourceAfter, parserList, Tokens} from 'constants/Editor.constants';
import {useAppStore} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import SpanSelectors from 'selectors/Span.selectors';
import Env from 'utils/Env';

const {PokeshopHttp = ''} = Env.get('demoEndpoints');

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

interface IProps {
  testId: string;
  runId: string;
}

const useAutoComplete = ({testId, runId}: IProps) => {
  const {getState} = useAppStore();

  const getAttributeList = useCallback(() => {
    const state = getState();
    const spanIdList = SpanSelectors.selectMatchedSpans(state);
    const attributeList = AssertionSelectors.selectAttributeList(state, testId, runId, spanIdList);

    return uniqBy(attributeList, 'key');
  }, [getState, runId, testId]);

  return useCallback(
    async (context: CompletionContext) => {
      const {state, pos} = context;
      const tree = syntaxTree(state);
      const nodeBefore = tree.resolveInner(pos, -1);
      const attributeList = getAttributeList();

      if (nodeBefore.name === Tokens.Pipe) {
        return {
          from: nodeBefore.to,
          options: parserList,
        };
      }

      if (completeSourceAfter.includes(nodeBefore.name)) {
        const {from, to} = nodeBefore;
        const [sourceText] = state.doc.sliceString(from, to).split(':');

        return {
          from: nodeBefore.to,
          options:
            sourceText === 'env'
              ? environmentList.map(({key}) => ({
                  label: key,
                  type: 'variableName',
                  apply: key,
                }))
              : attributeList.map(({key, value}) => ({
                  label: key,
                  type: 'variableName',
                  apply: `${key} = ${JSON.stringify(value)}`,
                })),
        };
      }

      if (nodeBefore.prevSibling?.name === Tokens.Source) {
        const {from, to} = nodeBefore.prevSibling.firstChild || {from: 0, to: 0};
        const [sourceText] = state.doc.sliceString(from, to).split(':');

        return {
          from: nodeBefore.from,
          to: nodeBefore.to,
          options:
            sourceText === 'env'
              ? environmentList.map(({key}) => ({
                  label: key,
                  type: 'variableName',
                  apply: key,
                }))
              : attributeList.map(({key, value}) => ({
                  label: key,
                  type: 'variableName',
                  apply: `${key} = ${JSON.stringify(value)}`,
                })),
        };
      }

      return null;
    },
    [getAttributeList]
  );
};

export default useAutoComplete;
