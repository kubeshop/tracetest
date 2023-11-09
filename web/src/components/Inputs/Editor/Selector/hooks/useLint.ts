import {useCallback} from 'react';
import {uniqBy} from 'lodash';
import {syntaxTree} from '@codemirror/language';
import {Diagnostic, LintSource} from '@codemirror/lint';
import {useAppStore} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';

interface IProps {
  testId: string;
  runId: number;
}

const useLint = ({runId, testId}: IProps): LintSource => {
  const {getState} = useAppStore();

  const getValidAttributeList = useCallback(() => {
    const state = getState();
    const attributeList = uniqBy(AssertionSelectors.selectAllAttributeList(state, testId, runId), 'key');

    return attributeList.map(({key}) => key);
  }, [getState, runId, testId]);

  return useCallback(
    async view => {
      let diagnostics: Diagnostic[] = [];
      const validAttributeList = getValidAttributeList();

      syntaxTree(view.state)
        .cursor()
        .iterate(node => {
          if (node.name === 'Identifier') {
            const attributeName = view.state.doc.sliceString(node.from, node.to);

            if (!validAttributeList.includes(attributeName)) {
              diagnostics.push({
                from: node.from,
                to: node.to,
                severity: 'error',
                message: "Attribute doesn't exist",
              });
            }
          }
        });

      return diagnostics;
    },
    [getValidAttributeList]
  );
};

export default useLint;
