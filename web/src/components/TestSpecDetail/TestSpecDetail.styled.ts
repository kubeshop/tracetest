import {CheckCircleFilled, InfoCircleFilled, MinusCircleFilled} from '@ant-design/icons';
import {Card, Drawer, Typography} from 'antd';
import styled from 'styled-components';

import {SemanticGroupNames, SemanticGroupNamesToColor} from 'constants/SemanticGroupNames.constants';

export const AssertionsContainer = styled.div`
  cursor: pointer;
`;

export const AssertionContainer = styled.div`
  span {
    overflow-wrap: anywhere;
  }
`;

export const CardContainer = styled(Card)<{$isSelected: boolean; $type: SemanticGroupNames}>`
  border: ${({$isSelected, theme}) =>
    $isSelected ? `1px solid ${theme.color.interactive}` : `1px solid ${theme.color.borderLight}`};

  .ant-card-head {
    border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
    border-top: ${({$type}) => `4px solid ${SemanticGroupNamesToColor[$type]}`};
    background-color: ${({theme}) => theme.color.white};
    padding: 0;
  }

  > .ant-card-body {
    padding: 0;
  }

  .ant-card-head > .ant-card-head-wrapper > .ant-card-head-title {
    padding: 0;
  }
`;

export const DrawerContainer = styled(Drawer)`
  position: absolute;
  overflow: hidden;
`;

export const GridContainer = styled.div`
  display: grid;
  grid-template-columns: 4.5% 1fr;
  align-items: center;
`;

export const CheckItemContainer = styled.div`
  padding: 10px 12px 10px 42px;

  border-bottom: 1px solid ${({theme}) => theme.color.borderLight};
`;

export const HeaderContainer = styled.div`
  align-items: center;
  display: flex;
  justify-content: space-between;
`;

export const HeaderItem = styled.div`
  align-items: center;
  color: ${({theme}) => theme.color.text};
  display: flex;
  font-size: ${({theme}) => theme.size.sm};
`;

export const HeaderItemText = styled(Typography.Text)`
  color: inherit;
  margin-left: 5px;
`;

export const HeaderTitle = styled(Typography.Title)`
  && {
    margin-bottom: 0;
  }
`;

export const IconError = styled(MinusCircleFilled)`
  color: ${({theme}) => theme.color.error};
`;

export const IconSuccess = styled(CheckCircleFilled)`
  color: ${({theme}) => theme.color.success};
`;

export const IconWarning = styled(InfoCircleFilled)`
  color: ${({theme}) => theme.color.warningYellow};
`;

export const Row = styled.div<{$align?: string; $hasGap?: boolean; $justify?: string}>`
  align-items: ${({$align}) => $align || 'center'};
  display: flex;
  justify-content: ${({$justify}) => $justify || 'flex-start'};
  gap: ${({$hasGap}) => $hasGap && '8px'};
`;

export const SecondaryText = styled(Typography.Text)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.sm};
`;

export const SpanHeaderContainer = styled.div`
  align-items: center;
  cursor: pointer;
  display: flex;
  gap: 8px;
`;

export const Wrapper = styled.div`
  align-items: center;
  cursor: pointer;
  justify-content: space-between;
  display: flex;
  padding: 8px 12px;
`;
