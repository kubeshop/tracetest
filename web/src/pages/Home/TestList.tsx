import {Table} from 'antd';
import {useNavigate} from 'react-router-dom';
import {useGetTestsQuery} from 'services/TestService';
import CustomTable from '../../components/CustomTable';
import {Test} from '../../types';
import NoResults from './NoResults';

const TestList = () => {
  const navigate = useNavigate();
  const {data: testList = [], isLoading} = useGetTestsQuery();

  return (
    <CustomTable
      dataSource={testList?.map(el => ({...el, url: el.serviceUnderTest.request.url})).reverse()}
      rowKey="testId"
      locale={{emptyText: <NoResults />}}
      loading={isLoading}
      onRow={record => {
        return {
          onClick: () => navigate(`/test/${(record as Test).testId}`),
        };
      }}
    >
      <Table.Column title="Name" dataIndex="name" key="name" width="25%" />
      <Table.Column title="Endpoint" dataIndex="url" key="url" />
    </CustomTable>
  );
};

export default TestList;
