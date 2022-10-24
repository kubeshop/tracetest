import {CheckCircleFilled, MinusCircleFilled} from '@ant-design/icons';
import {Card, Drawer, Typography} from 'antd';
import styled from 'styled-components';

import {SemanticGroupNames, SemanticGroupNamesToColor} from 'constants/SemanticGroupNames.constants';

export const AssertionContainer = styled.div`
  span {
    overflow-wrap: anywhere;
  }
`;

export const CardContainer = styled(Card)<{$isSelected: boolean; $type: SemanticGroupNames}>`
  border: ${({$isSelected, theme}) =>
    $isSelected ? `1px solid ${theme.color.interactive}` : `1px solid ${theme.color.borderLight}`};

  :not(:last-child) {
    margin-bottom: 16px;
  }

  .ant-card-head {
    border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
    border-top: ${({$type}) => `4px solid ${SemanticGroupNamesToColor[$type]}`};
    background-color: ${({theme}) => theme.color.white};
  }

  > .ant-card-body {
    padding: 0px 12px;
  }
`;

export const DrawerContainer = styled(Drawer)`
  position: absolute;
  overflow: hidden;
`;

export const GridContainer = styled.div`
  display: grid;
  column-gap: 24px;
  grid-template-columns: 5% 1fr;
  align-items: center;
`;

export const CheckItemContainer = styled.div`
  padding: 10px 0 10px 30px;

  &:hover {
    background: ${({theme}) => theme.color.background};
  }
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
  margin-right: 8px;
`;
export const IconSuccess = styled(CheckCircleFilled)`
  color: ${({theme}) => theme.color.success};
  margin-right: 8px;
`;

export const Row = styled.div<{$align?: string; $hasGap?: boolean}>`
  align-items: ${({$align}) => $align || 'center'};
  display: flex;
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
