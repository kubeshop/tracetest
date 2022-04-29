import {Table} from 'antd';
import {FC} from 'react';
import { ISpanFlatAttribute } from '../../types/Span.types';
import CustomTable from '../CustomTable';

type TSpanAttributesTableProps = {
  spanAttributesList: ISpanFlatAttribute[];
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
