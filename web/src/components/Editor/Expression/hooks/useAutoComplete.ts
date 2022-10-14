import {useCallback} from 'react';
import {uniqBy} from 'lodash';
import {CompletionContext} from '@codemirror/autocomplete';
import {useAppStore} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import SpanSelectors from 'selectors/Span.selectors';
import EditorService from 'services/Editor.service';
import {SupportedEditors} from 'constants/Editor.constants';

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
      const attributeList = getAttributeList();

      return EditorService.getAutocomplete(SupportedEditors.Expression, context, attributeList);
    },
    [getAttributeList]
  );
};

export default useAutoComplete;
