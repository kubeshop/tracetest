import {CloseCircleOutlined} from '@ant-design/icons';
import {Col, Row, Space} from 'antd';
import {debounce} from 'lodash';
import {useCallback, useMemo, useState} from 'react';

import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useLazyGetSelectedSpansQuery} from 'redux/apis/TraceTest.api';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {matchSpans, selectSpan, setSearchText} from 'redux/slices/Trace.slice';
import TraceSelectors from 'selectors/Trace.selectors';
import SpanService from 'services/Span.service';
import EditorService from 'services/Editor.service';
import * as S from './RunDetailTrace.styled';

interface IProps {
  runId: string;
  testId: string;
}

const Search = ({runId, testId}: IProps) => {
  const [search, setSearch] = useState('');
  const dispatch = useAppDispatch();
  const matchedSpans = useAppSelector(TraceSelectors.selectMatchedSpans);
  const {
    run: {trace: {spans = []} = {}},
  } = useTestRun();
  const [getSelectedSpans] = useLazyGetSelectedSpansQuery();

  const handleSearch = useCallback(
    async (query: string) => {
      const isValidSelector = EditorService.getIsQueryValid(SupportedEditors.Selector, query || '');
      if (!query) {
        dispatch(matchSpans({spanIds: []}));
        dispatch(selectSpan({spanId: ''}));
        return;
      }

      let spanIds = [];
      if (isValidSelector) {
        const selectedSpansData = await getSelectedSpans({query, runId, testId}).unwrap();
        spanIds = selectedSpansData.spanIds;
      } else {
        dispatch(setSearchText({searchText: query}));
        spanIds = SpanService.searchSpanList(spans, query);
      }

      dispatch(matchSpans({spanIds}));
      if (spanIds.length) {
        dispatch(selectSpan({spanId: spanIds[0]}));
      }
    },
    [dispatch, getSelectedSpans, runId, spans, testId]
  );

  const onSearch = useMemo(() => debounce(handleSearch, 500), [handleSearch]);
  const onClear = useCallback(() => {
    onSearch('');
    setSearch('');
  }, [onSearch]);

  return (
    <Row>
      <Col flex="auto">
        <Editor
          type={SupportedEditors.Selector}
          placeholder="Search in trace"
          onChange={query => {
            onSearch(query);
            setSearch(query);
          }}
          value={search}
        />
        {!!search && <S.ClearSearchIcon onClick={onClear} />}
        {!!search && !matchedSpans.length && (
          <S.NoMatchesContainer>
            <Space>
              <CloseCircleOutlined />
              <S.NoMatchesText>No matches found</S.NoMatchesText>
            </Space>
          </S.NoMatchesContainer>
        )}
      </Col>
    </Row>
  );
};

export default Search;
