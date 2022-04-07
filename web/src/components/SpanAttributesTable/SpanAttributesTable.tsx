import {Table} from 'antd';
import {FC} from 'react';
import {TSpanAttributesList} from 'types';
import CustomTable from '../CustomTable';

type TSpanAttributesTableProps = {
  spanAttributesList: TSpanAttributesList;
};

const SpanAttributesTable: FC<TSpanAttributesTableProps> = ({spanAttributesList}) => {
  return (
    <CustomTable
      size="small"
      pagination={{hideOnSinglePage: true}}
      dataSource={spanAttributesList}
      bordered
      tableLayout="fixed"
      showHeader={false}
    >
      <Table.Column
        dataIndex="key"
        key="key"
        width="20%"
        render={value => ({
          props: {
            style: {
              background: '#FAFAFA',
            },
          },
          children: value,
        })}
      />
      <Table.Column dataIndex="value" key="value" />
    </CustomTable>
  );
};

export default SpanAttributesTable;
