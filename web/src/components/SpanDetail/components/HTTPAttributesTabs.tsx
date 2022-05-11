import {Table, Tabs, Typography} from 'antd';
import React from 'react';
import AttributesService from '../../../services/Atributes.service';
import CustomTable from '../../CustomTable';
import {HTTPAttributesTabsComponent} from './Component.styled';

interface IProps {
  attributeList: any[];
}

export const HTTPAttributesTabs: React.FC<IProps> = ({attributeList}) => {
  const render = (value: any) => <span data-cy="assertion-check-property">{value}</span>;
  return (
    <>
      <Typography.Title level={5} style={{margin: 0}}>
        Attributes
      </Typography.Title>
      <HTTPAttributesTabsComponent
        data-cy="span-details-attributes"
        id="the_tabs"
        style={{marginBottom: 150, overflow: 'hidden'}}
      >
        <Tabs.TabPane style={{overflow: 'hidden'}} tab={<span>Request</span>} key="span-request-detail">
          <CustomTable
            size="small"
            pagination={{hideOnSinglePage: true}}
            dataSource={AttributesService.requestAttributesListFN(attributeList)}
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
            dataSource={AttributesService.responseAttributesFN(attributeList)}
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
      </HTTPAttributesTabsComponent>
    </>
  );
};
