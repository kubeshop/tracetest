import {Badge, Tooltip, Typography} from 'antd';
import styled from 'styled-components';
import {SemanticGroupNames, SemanticGroupNamesToColor} from 'constants/SemanticGroupNames.constants';

export const AssertionCheckRow = styled.div`
  display: grid;
  grid-template-columns: 1.5fr 1fr 0.8fr 1fr 1fr;
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

export const SelectorEntry = styled.div`
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
`;

export const Label = styled(Typography.Text).attrs({
  type: 'secondary',
})`
  font-size: ${({theme}) => theme.size.sm};
`;

export const Value = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.sm};
`;

export const LabelBadge = styled(Badge)<{$spanType?: SemanticGroupNames}>`
  > sup {
    background: ${({$spanType, theme}) => ($spanType ? SemanticGroupNamesToColor[$spanType] : theme.color.borderLight)};
    color: black;
    margin-right: 6px;
    border-radius: 2px;
    margin-bottom: 8px;
  }
`;

export const LabelTooltip = styled(Tooltip).attrs({
  placement: 'top',
  arrowPointAtCenter: true,
})``;

export const SelectedLabelBadge = styled(LabelBadge)`
  > sup {
    color: ${({theme}) => theme.color.primary};
    border: ${({theme}) => `1px solid ${theme.color.primary}`};
  }
`;
