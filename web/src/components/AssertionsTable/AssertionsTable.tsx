import {Table} from 'antd';
import {ColumnsType} from 'antd/lib/table';

const AssertionsResultTable = () => {
  const dataSource = [
    {
      key: '1',
      name: 'HTTP',
      selectedSpans: 4,
      property: 'http.status_code',
      comparison: 'equals',
      value: '200',
      results: '3 / 1',
    },
  ];

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
        if (dataSource.filter(el => el.name === value).length === 1) {
          return obj;
        }
        if (dataSource.findIndex(el => el.name === value) === index) {
          const count = dataSource.filter(item => item.name === value).length;
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

  return <Table size="small" pagination={{hideOnSinglePage: true}} dataSource={dataSource} columns={columns} />;
};

export default AssertionsResultTable;
