import {Skeleton, SkeletonProps, Table} from 'antd';
import {ColumnsType} from 'antd/lib/table';

export type SkeletonTableColumnsType = {
  key: string;
};

type SkeletonTableProps = SkeletonProps & {
  columns?: ColumnsType<SkeletonTableColumnsType>;
  rowCount?: number;
};

export const SkeletonTable = ({
  loading = false,
  active = false,
  rowCount = 5,
  columns = [{key: 'loading'}] as SkeletonTableColumnsType[],
  children,
  className,
}: SkeletonTableProps): JSX.Element => {
  return loading ? (
    <div>
      <Table
        showHeader={false}
        rowKey="key"
        pagination={false}
        dataSource={[...Array(rowCount)].map((_, index) => ({
          key: `key${index}`,
        }))}
        columns={columns.map(column => {
          return {
            ...column,
            render: function renderPlaceholder() {
              return <Skeleton key={column.key} title active={active} paragraph={false} className={className} />;
            },
          };
        })}
      />
    </div>
  ) : (
    <div>{children}</div>
  );
};

export default SkeletonTable;
