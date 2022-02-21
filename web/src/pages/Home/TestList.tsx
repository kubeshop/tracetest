import {Table} from 'antd';
import {useGetTestsQuery} from '../../services/TestService';

const TestList = () => {
  const {data: tests, isLoading} = useGetTestsQuery();
  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Url',
      dataIndex: 'url',
      key: 'url',
    },
  ];
  return (
    <Table
      dataSource={tests?.map(el => ({...el, url: el.serviceUnderTest.url}))}
      rowKey={test => test.id}
      loading={isLoading}
      columns={columns}
    />
  );
};

export default TestList;
