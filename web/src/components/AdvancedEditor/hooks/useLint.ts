import {useCallback} from 'react';
import {syntaxTree} from '@codemirror/language';
import {Diagnostic, LintSource} from '@codemirror/lint';
import {useAppStore} from '../../../redux/hooks';
import AssertionSelectors from '../../../selectors/Assertion.selectors';

interface IProps {
  testId: string;
  runId: string;
}

const useLint = ({runId, testId}: IProps): LintSource => {
  const {getState} = useAppStore();

  const getAttributeList = useCallback(() => {
    const state = getState();
    const defaultList = AssertionSelectors.selectAllAttributeList(state, testId, runId);

    return defaultList;
  }, [getState, runId, testId]);

  return useCallback(
    async view => {
      let diagnostics: Diagnostic[] = [];
      const attributeList = getAttributeList();
      const validAttributeList = attributeList.map(({key}) => key);

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
    [getAttributeList]
  );
};

export default useLint;
