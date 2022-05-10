import {Table, Tabs, Typography} from 'antd';
import React from 'react';
import styled from 'styled-components';
import CustomTable from '../../CustomTable';

const STabs = styled(Tabs)`
  .ant-tabs-content-holder {
    overflow: hidden;
    padding-bottom: 100px;
  }
`;

interface IProps {
  attributeList: any[];
}

export const HTTPAttributesTabs = ({attributeList}: IProps): JSX.Element => {
  const render = (value: any) => <span data-cy="assertion-check-property">{value}</span>;
  return (
    <>
      <Typography.Title level={5} style={{margin: 0}}>
        Attributes
      </Typography.Title>
      <STabs style={{marginBottom: 150, overflow: 'hidden'}}>
        <Tabs.TabPane style={{overflow: 'hidden'}} tab={<span>Request</span>} key="span-request-detail">
          <CustomTable
            size="small"
            pagination={{hideOnSinglePage: true}}
            dataSource={attributeList?.filter(a => {
              return (
                [
                  'http.method',
                  'http.url',
                  'http.target',
                  'http.host',
                  'http.scheme',
                  'http.request_content_length',
                  'http.request_content_length_uncompressed',
                  'http.retry_count exists',
                  'http.user_agent',
                ].includes(a.key) || a.key.includes('http.request')
              );
            })}
          >
            <Table.Column
              title="Key"
              dataIndex="key"
              key="key"
              ellipsis
              width="50%"
              render={value => <span data-cy="assertion-check-property">{value}</span>}
            />
            <Table.Column title="Value" dataIndex="value" key="value" ellipsis width="50%" render={render} />
          </CustomTable>
        </Tabs.TabPane>
        <Tabs.TabPane tab={<span>Response</span>} key="span-response-results">
          <CustomTable
            size="small"
            pagination={{hideOnSinglePage: true}}
            dataSource={attributeList?.filter(a => {
              return a.key.includes('http.response');
            })}
          >
            <Table.Column title="Key" dataIndex="key" key="key" ellipsis width="50%" render={render} />
            <Table.Column title="Value" dataIndex="value" key="value" ellipsis width="50%" render={render} />
          </CustomTable>
        </Tabs.TabPane>
        <Tabs.TabPane tab={<span>Detailed attribute list</span>} key="detailed-attribute-list">
          <CustomTable size="small" pagination={{hideOnSinglePage: true}} dataSource={attributeList}>
            <Table.Column title="Key" dataIndex="key" key="key" ellipsis width="50%" render={render} />
            <Table.Column title="Value" dataIndex="value" key="value" ellipsis width="50%" render={render} />
          </CustomTable>
        </Tabs.TabPane>
      </STabs>
    </>
  );
};
