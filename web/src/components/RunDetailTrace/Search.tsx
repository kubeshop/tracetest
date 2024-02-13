import {CloseCircleOutlined} from '@ant-design/icons';
import {Col, Row, Space} from 'antd';
import {debounce} from 'lodash';
import {useCallback, useMemo, useState} from 'react';

import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';
import TracetestAPI from 'redux/apis/Tracetest';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {matchSpans, selectSpan, setSearchText} from 'redux/slices/Trace.slice';
import TraceSelectors from 'selectors/Trace.selectors';
import * as S from './RunDetailTrace.styled';

const {useGetSearchedSpansMutation} = TracetestAPI.instance;

interface IProps {
  runId: number;
  testId: string;
}

const Search = ({runId, testId}: IProps) => {
  const [search, setSearch] = useState('');
  const dispatch = useAppDispatch();
  const matchedSpans = useAppSelector(TraceSelectors.selectMatchedSpans);
  const [getSearchedSpans] = useGetSearchedSpansMutation();

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
    <Row>
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
