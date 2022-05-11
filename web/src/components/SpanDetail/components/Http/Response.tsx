import {Table} from 'antd';
import React from 'react';
import {ISpanFlatAttribute} from '../../../../types/Span.types';
import CustomTable from '../../../CustomTable';

interface IRequestProps {
  attributeList: ISpanFlatAttribute[];
}

const filterAttributeList = (attributeList: ISpanFlatAttribute[]) =>
  attributeList?.filter(a => a.key.includes('http.response'));

const Request: React.FC<IRequestProps> = ({attributeList}) => {
  return (
    <CustomTable
      size="small"
      pagination={{hideOnSinglePage: true}}
      dataSource={filterAttributeList(attributeList)}
    >
      <Table.Column title="Key" dataIndex="key" key="key" ellipsis width="50%" />
      <Table.Column title="Value" dataIndex="value" key="value" ellipsis width="50%" />
    </CustomTable>
  );
};

export default Request;
