import {Table} from 'antd';
import {FC} from 'react';
import { TSpanFlatAttribute } from '../../entities/Span/Span.types';
import CustomTable from '../CustomTable';

type TSpanAttributesTableProps = {
  spanAttributesList: TSpanFlatAttribute[];
};

const SpanAttributesTable: FC<TSpanAttributesTableProps> = ({spanAttributesList}) => {
  return (
    <CustomTable
      size="small"
      pagination={false}
      dataSource={spanAttributesList}
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
