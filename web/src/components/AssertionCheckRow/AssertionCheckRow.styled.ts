import {Badge, Typography} from 'antd';
import styled from 'styled-components';

export const AssertionCheckRow = styled.div`
  display: grid;
  grid-template-columns: 1.5fr 1fr .8fr 1fr 1fr;
  gap: 14px;
  cursor: pointer;
`;

export const Entry = styled.div`
  display: flex;
  flex-direction: column;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const Label = styled(Typography.Text).attrs({
  type: 'secondary',
})`
  font-size: 12px;
`;

export const Value = styled(Typography.Text)`
  font-size: 12px;
`;

export const LabelBadge = styled(Badge)`
  > sup {
    background-color: #f0f0f0;
    color: black;
    margin-right: 6px;
    border-radius: 2px;
  }
`;

export const SelectedLabelBadge = styled(LabelBadge)`
  > sup {
    color: #61175e;
    border: 1px solid #61175e;
  }
`;
