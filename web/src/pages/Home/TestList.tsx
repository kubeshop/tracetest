import {Table} from 'antd';
import {useNavigate} from 'react-router-dom';
import {useGetTestsQuery} from '../../services/TestService';

const TestList = () => {
  const navigate = useNavigate();
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
      dataSource={tests?.map(el => ({...el, url: el.serviceUnderTest.url})).reverse()}
      rowKey={test => test.id}
      loading={isLoading}
      columns={columns}
      onRow={(record, rowIndex) => {
        return {
          onClick: () => navigate(`/test/${record.id}`),
        };
      }}
    />
  );
};

export default TestList;
