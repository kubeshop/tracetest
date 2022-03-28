import {Table} from 'antd';
import {ColumnsType} from 'antd/lib/table';
import {AssertionResult} from 'types';

interface IProps {
  assertionResults: AssertionResult[];
}

const AssertionsResultTable = ({assertionResults}: IProps) => {
  const data = assertionResults.map(el => {
    return {
      key: el.spanAssertionId,
      name: el.selector,
      selectedSpans: el.spanCount,
      property: el.propertyName,
      comparison: el.operator,
      value: el.comparisonValue,
      results: `${el.passedSpanCount}/${el.spanCount - el.passedSpanCount}`,
    };
  });

  const columns: ColumnsType<any> = [
    {
      title: 'Span selector',
      dataIndex: 'name',
      key: 'name',
      render: (value, record, index) => {
        const obj = {
          children: value,
          props: {rowSpan: 1},
        };
        if (data.filter(el => el.name === value).length === 1) {
          return obj;
        }
        if (data.findIndex(el => el.name === value) === index) {
          const count = data.filter(item => item.name === value).length;
          obj.props.rowSpan = count;
          return obj;
        }
        obj.props.rowSpan = 0;
        return obj;
      },
    },
    {
      title: '# Selected',
      dataIndex: 'selectedSpans',
      key: 'selectedSpans',
    },

    {
      title: 'Property',
      dataIndex: 'property',
      key: 'property',
    },
    {
      title: 'Comparison',
      dataIndex: 'comparison',
      key: 'comparison',
    },
    {
      title: 'Value',
      dataIndex: 'value',
      key: 'value',
    },
    {
      title: 'Pass/Fail',
      dataIndex: 'results',
      key: 'results',
    },
  ];

  return <Table size="small" pagination={{hideOnSinglePage: true}} dataSource={data} columns={columns} />;
};

export default AssertionsResultTable;
