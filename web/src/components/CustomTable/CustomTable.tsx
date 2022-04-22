import {Table} from 'antd';
import styled from 'styled-components';

const CustomTable = styled(Table).attrs({
  bordered: true,
  tableLayout: 'fixed',
})`
  .ant-table-thead > tr > th {
    font-weight: 600;
  }

  .ant-table.ant-table-bordered > .ant-table-container > .ant-table-body > table > tbody > tr > td,
  .ant-table.ant-table-bordered > .ant-table-container > .ant-table-header > table > thead > tr > th {
    border-right: none;
  }

  .ant-table.ant-table-bordered > .ant-table-body {
    border-right: 1px solid #f0f0f0;
  }

  tr {
    cursor: pointer;
  }
`;

export default CustomTable;
