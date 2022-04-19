import {Button} from 'antd';
import {ColumnsType} from 'antd/lib/table';
import Title from 'antd/lib/typography/Title';
import CustomTable from '../../components/CustomTable';

const Assertions = () => {
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
    {
      key: '2',
      name: 'HTTP POST /user/verify',
      selectedSpans: 1,
      property: 'http.response.body',
      comparison: 'contains',
      value: 'ok',
      results: '1 / 0',
    },
    {
      key: '3',
      name: 'HTTP POST /user/verify',
      selectedSpans: 1,
      property: 'traceDurationMs',
      comparison: 'less-than',
      value: '1000 ms',
      results: '1 / 0',
    },
    {
      key: '4',
      name: 'HTTP POST /user/verify',
      selectedSpans: 1,
      property: 'http.status_text',
      comparison: 'equals',
      value: 'OK',
      results: '1 / 0',
    },

    {
      key: '5',
      name: 'MONGO stock-db find statement contains cartId',
      selectedSpans: 1,
      property: 'array(db.response).length',
      comparison: 'equals',
      value: '2',
      results: '1 / 0',
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
  ];

  return (
    <>
      <div style={{display: 'flex', justifyContent: 'space-between', marginBottom: 16}}>
        <Title level={4}>Assertions</Title>
        <Button>New Assertion</Button>
      </div>
      <CustomTable dataSource={dataSource} columns={columns} />
    </>
  );
};

export default Assertions;
