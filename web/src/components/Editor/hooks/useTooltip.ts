import {useCallback} from 'react';
import {EditorView} from '@codemirror/view';
import EditorService from 'services/Editor.service';
import EnvironmentSelectors from 'selectors/Environment.selectors';
import {useAppStore} from 'redux/hooks';
import SpanSelectors from 'selectors/Span.selectors';
import AssertionSelectors from 'selectors/Assertion.selectors';
import {uniqBy} from 'lodash';

interface IProps {
  testId?: string;
  runId?: string;
}

const useTooltip = ({testId = '', runId = ''}: IProps = {}) => {
  const {getState} = useAppStore();

  const getAttributeList = useCallback(() => {
    const state = getState();
    const spanIdList = SpanSelectors.selectMatchedSpans(state);
    const attributeList = AssertionSelectors.selectAttributeList(state, testId, runId, spanIdList);

    return uniqBy(attributeList, 'key');
  }, [getState, runId, testId]);

  const getSelectedEnvironmentEntryList = useCallback(() => {
    const state = getState();

    return EnvironmentSelectors.selectSelectedEnvironmentValues(state);
  }, [getState]);

  return useCallback(
    async (view: EditorView, pos: number, side: -1 | 1) => {
      const attributeList = testId && runId ? getAttributeList() : [];
      const envEntryList = getSelectedEnvironmentEntryList();

      return EditorService.getTooltip({view, pos, side, attributeList, envEntryList});
    },
    [getAttributeList, getSelectedEnvironmentEntryList, runId, testId]
  );
};

export default useTooltip;
