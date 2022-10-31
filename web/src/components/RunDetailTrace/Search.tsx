import {Col, Row} from 'antd';
import {debounce} from 'lodash';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useCallback, useMemo, useState} from 'react';
import {useLazyGetSelectedSpansQuery} from 'redux/apis/TraceTest.api';
import {useAppDispatch} from 'redux/hooks';
import {matchSpans, selectSpan, setSearchText} from 'redux/slices/Trace.slice';
import SpanService from 'services/Span.service';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';
import EditorService from 'services/Editor.service';
import * as S from './RunDetailTrace.styled';

interface IProps {
  runId: string;
  testId: string;
}

const Search = ({runId, testId}: IProps) => {
  const [search, setSearch] = useState('');
  const dispatch = useAppDispatch();
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
      dispatch(selectSpan({spanId: spanIds[0]}));
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
        {Boolean(search) && <S.ClearSearchIcon onClick={onClear} />}
      </Col>
    </Row>
  );
};

export default Search;
