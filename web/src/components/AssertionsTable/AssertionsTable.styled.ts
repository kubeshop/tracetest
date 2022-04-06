import {Badge} from 'antd';
import styled from 'styled-components';

export const AssertionsTableContainer = styled.div`
  margin-bottom: 32px;
`;

export const AssertionsTableHeader = styled.div`
  padding: 8px 0;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  width: 100%;
`;

export const AssertionsTableBadge = styled(Badge)`
  > sup {
    background-color: #f0f0f0;
    color: black;
    margin-left: 6px;
  }
`;
