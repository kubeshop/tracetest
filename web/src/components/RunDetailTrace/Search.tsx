import {Button, Col, Form, Row} from 'antd';

import AdvancedEditor from 'components/AdvancedEditor';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useLazyGetSelectedSpansQuery} from 'redux/apis/TraceTest.api';
import {useAppDispatch} from 'redux/hooks';
import {matchSpans, selectSpan, setSearchText} from 'redux/slices/Trace.slice';
import SelectorService from 'services/Selector.service';
import SpanService from 'services/Span.service';

interface IValues {
  query: string;
}

interface IProps {
  runId: string;
  testId: string;
}

const Search = ({runId, testId}: IProps) => {
  const dispatch = useAppDispatch();
  const {
    run: {trace: {spans = []} = {}},
  } = useTestRun();
  const [getSelectedSpans] = useLazyGetSelectedSpansQuery();

  const [form] = Form.useForm<IValues>();

  const onFinish = async (values: IValues) => {
    const {query} = values;
    const isValidSelector = SelectorService.getIsValidSelector(query);

    let spanIds = [];
    if (isValidSelector) {
      spanIds = await getSelectedSpans({query, runId, testId}).unwrap();
    } else {
      dispatch(setSearchText({searchText: query}));
      spanIds = SpanService.searchSpanList(spans, values.query);
    }

    dispatch(matchSpans({spanIds}));
    dispatch(selectSpan({spanId: spanIds[0]}));
  };

  const onClear = () => {
    form.resetFields();
    dispatch(setSearchText({searchText: ''}));
    dispatch(matchSpans({spanIds: []}));
  };

  return (
    <Form form={form} name="search" onFinish={onFinish}>
      <Row>
        <Col flex="auto">
          <Form.Item
            name="query"
            rules={[
              {
                required: true,
                message: 'Try with a search term or a query using our Selector Language',
              },
            ]}
          >
            <AdvancedEditor placeholder="Search in trace" runId={runId} testId={testId} />
          </Form.Item>
        </Col>

        <Col flex="150px">
          <Button type="primary" htmlType="submit">
            Search
          </Button>
          <Button onClick={onClear}>Clear</Button>
        </Col>
      </Row>
    </Form>
  );
};

export default Search;
