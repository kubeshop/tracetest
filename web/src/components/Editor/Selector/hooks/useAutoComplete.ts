import {useCallback} from 'react';
import {uniqBy} from 'lodash';
import {CompletionContext} from '@codemirror/autocomplete';
import {syntaxTree} from '@codemirror/language';
import {
  completeIdentifierAfter,
  completePseudoSelectorAfter,
  completeValueAfter,
  comparatorList,
  pseudoSelectorList,
  Tokens,
} from 'constants/Editor.constants';
import {useAppStore} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import {escapeString} from 'utils/Common';

interface IProps {
  testId: string;
  runId: number;
}

const useAutoComplete = ({testId, runId}: IProps) => {
  const {getState} = useAppStore();

  const getAttributeList = useCallback(() => {
    const state = getState();
    const defaultList = AssertionSelectors.selectAllAttributeList(state, testId, runId);

    return defaultList;
  }, [getState, runId, testId]);

  return useCallback(
    async (context: CompletionContext) => {
      const {state, pos} = context;
      const word = context.matchBefore(/\w*/) ?? {
        from: 0,
        to: 0,
      };

      const tree = syntaxTree(state);
      const nodeBefore = tree.resolveInner(pos, -1);
      const attributeList = getAttributeList();

      if (
        nodeBefore.prevSibling?.name === Tokens.Identifier ||
        nodeBefore.nextSibling?.name === Tokens.ClosingBracket
      ) {
        return {
          from: word.from,
          options: comparatorList,
        };
      }

      if (completeIdentifierAfter.includes(nodeBefore.name)) {
        const uniqueList = uniqBy(attributeList, 'key');
        const identifierText = state.doc.sliceString(nodeBefore.from, nodeBefore.to);
        const isIdentifier = nodeBefore.name === Tokens.Identifier;
        const list = isIdentifier ? uniqueList.filter(({key}) => key.toLowerCase().includes(identifierText.toLowerCase())) : uniqueList;

        return {
          from: isIdentifier ? nodeBefore.from : word.from,
          to: isIdentifier ? nodeBefore.to : word.to,
          options: list.map(({key, value}) => ({
            label: key,
            type: 'variableName',
            apply: `${key} = ${JSON.stringify(value)}`,
          })),
        };
      }

      if (completePseudoSelectorAfter.includes(nodeBefore.name)) {
        return {
          from: word.from,
          options: pseudoSelectorList,
        };
      }

      const parentNode = nodeBefore.parent;
      const identifierNode = parentNode?.prevSibling?.prevSibling;
      if (parentNode && completeValueAfter.includes(parentNode.name) && identifierNode) {
        const attributeName = state.doc.sliceString(identifierNode.from, identifierNode.to);
        const valueList = attributeList.filter(({key}) => key === attributeName);

        return {
          from: parentNode.from + 1,
          to: parentNode.to - 1,
          options: valueList.map(({value}) => ({label: value, apply: escapeString(value), type: 'string'})),
        };
      }

      const planeWord = context.matchBefore(/\w*/) ?? {from: 0, to: 0};

      if (planeWord.from === planeWord.to && !context.explicit) return null;

      return {
        from: planeWord.from,
        options: [{label: 'span', apply: 'span[', type: 'variableName'}],
      };
    },
    [getAttributeList]
  );
};

export default useAutoComplete;
