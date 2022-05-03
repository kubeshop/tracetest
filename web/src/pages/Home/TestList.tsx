import {Table} from 'antd';
import {useRef} from 'react';
import {useNavigate} from 'react-router-dom';
import {useGetTestListQuery} from 'redux/apis/Test.api';
import CustomTable from '../../components/CustomTable';
import HomeAnalyticsService from '../../services/Analytics/HomeAnalytics.service';
import NoResults from './NoResults';
import {TTest} from '../../types/Test.types';

const {onTestClick} = HomeAnalyticsService;

const TestList = () => {
  const navigate = useNavigate();
  const eventRef = useRef<{previousPageX: number; currentPageX: number}>({previousPageX: 0, currentPageX: 0});
  const {data: testList = [], isLoading} = useGetTestListQuery();

  const handleMouseUp = (event: any) => {
    if (event.type === 'mousedown') {
      eventRef.current.previousPageX = event.pageX;
    } else if (event.type === 'mouseup') {
      eventRef.current.currentPageX = event.pageX;
    } else if (event.type === 'click') {
      if (Math.abs(eventRef.current.currentPageX - eventRef.current.previousPageX) > 0) {
        event?.stopPropagation();
      }
    }
  };

  return (
    <CustomTable
      scroll={{y: 'calc(100vh - 300px)'}}
      dataSource={testList?.map(el => ({...el, url: el?.serviceUnderTest?.request?.url})).reverse()}
      rowKey="testId"
      locale={{emptyText: <NoResults />}}
      loading={isLoading}
      onRow={record => {
        return {
          onClick: () => {
            const testId = (record as TTest).testId;

            onTestClick(testId!);
            navigate(`/test/${testId}`);
          },
        };
      }}
    >
      <Table.Column title="Name" dataIndex="name" key="name" width="25%" />
      <Table.Column
        title="Endpoint"
        dataIndex="url"
        key="url"
        render={value => {
          return (
            <span
              style={{paddingLeft: 16, paddingRight: 16}}
              onMouseDown={handleMouseUp}
              onMouseUp={handleMouseUp}
              onClick={handleMouseUp}
            >
              {value}
            </span>
          );
        }}
      />
    </CustomTable>
  );
};

export default TestList;
