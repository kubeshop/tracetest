import {Table} from 'antd';
import React from 'react';
import {ISpanFlatAttribute} from '../../../../types/Span.types';
import CustomTable from '../../../CustomTable';

interface IRequestProps {
  attributeList: ISpanFlatAttribute[];
}

const HttpRequestAttributeList = [
  'http.method',
  'http.url',
  'http.target',
  'http.host',
  'http.scheme',
  'http.request_content_length',
  'http.request_content_length_uncompressed',
  'http.retry_count exists',
  'http.user_agent',
];

const filterAttributeList = (attributeList: ISpanFlatAttribute[]) =>
  attributeList?.filter(a => HttpRequestAttributeList.includes(a.key) || a.key.includes('http.request'));

const Request: React.FC<IRequestProps> = ({attributeList}) => {
  return (
    <CustomTable
      size="small"
      pagination={{hideOnSinglePage: true}}
      dataSource={filterAttributeList(attributeList)}
    >
      <Table.Column
        title="Key"
        dataIndex="key"
        key="key"
        ellipsis
        width="50%"
        render={value => <span data-cy="assertion-check-property">{value}</span>}
      />
      <Table.Column title="Value" dataIndex="value" key="value" ellipsis width="50%" />
    </CustomTable>
  );
};

export default Request;
