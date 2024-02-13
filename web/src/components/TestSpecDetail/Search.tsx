import {Col} from 'antd';
import {debounce} from 'lodash';
import {useCallback, useMemo, useState} from 'react';

import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';
import TracetestAPI from 'redux/apis/Tracetest';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useTest} from 'providers/Test/Test.provider';
import {useAppDispatch} from 'redux/hooks';
import {matchSpans, selectSpan, setSearchText} from 'redux/slices/Trace.slice';
import * as S from './TestSpecDetail.styled';

const {useGetSearchedSpansMutation} = TracetestAPI.instance;

const Search = () => {
  const [search, setSearch] = useState('');
  const dispatch = useAppDispatch();
  const [getSearchedSpans] = useGetSearchedSpansMutation();
  const {
    run: {id: runId},
  } = useTestRun();
  const {
    test: {id: testId},
  } = useTest();

  const handleSearch = useCallback(
    async (query: string) => {
      if (!query) {
        dispatch(matchSpans({spanIds: []}));
        dispatch(selectSpan({spanId: ''}));
        return;
      }

      const {spanIds} = await getSearchedSpans({query, runId, testId}).unwrap();
      dispatch(setSearchText({searchText: query}));
      dispatch(matchSpans({spanIds}));

      if (spanIds.length) {
        dispatch(selectSpan({spanId: spanIds[0]}));
      }
    },
    [dispatch, getSearchedSpans, runId, testId]
  );

  const onSearch = useMemo(() => debounce(handleSearch, 500), [handleSearch]);
  const onClear = useCallback(() => {
    onSearch('');
    setSearch('');
  }, [onSearch]);

  return (
    <S.SearchContainer>
      <Col flex="auto">
        <Editor
          type={SupportedEditors.Selector}
          placeholder={'Try `span[tracetest.span.type="general" name="Tracetest trigger"]` or just "Tracetest trigger"'}
          onChange={query => {
            onSearch(query);
            setSearch(query);
          }}
          value={search}
        />
        {!!search && <S.ClearSearchIcon onClick={onClear} />}
      </Col>
    </S.SearchContainer>
  );
};

export default Search;
