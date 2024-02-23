import {CheckCircleFilled, CloseCircleFilled, MinusCircleFilled} from '@ant-design/icons';
import {Drawer, Typography} from 'antd';
import styled from 'styled-components';
import {SemanticGroupNames, SemanticGroupNamesToColor} from 'constants/SemanticGroupNames.constants';

export const AssertionsContainer = styled.div`
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 16px;
  border-top: 1px solid ${({theme}) => theme.color.borderLight};
`;

export const AssertionContainer = styled.div`
  span {
    overflow-wrap: anywhere;
  }
`;

export const DrawerContainer = styled(Drawer)<{$type: SemanticGroupNames}>`
  position: absolute;
  overflow: hidden;

  .ant-drawer-body {
    display: flex;
    flex-direction: column;
    border-top: ${({$type}) => `4px solid ${SemanticGroupNamesToColor[$type]}`};
  }
`;

export const DrawerRow = styled.div`
  flex: 1;
`;

export const GridContainer = styled.div`
  display: grid;
  grid-template-columns: 4.5% 1fr;
  align-items: center;
`;

export const CheckItemContainer = styled.div<{$isSuccessful: boolean}>`
  padding: 10px 12px;

  border: 1px solid ${({theme}) => theme.color.borderLight};

  border-top: ${({$isSuccessful, theme}) => `4px solid ${$isSuccessful ? theme.color.success : theme.color.error}`};
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
    text-overflow: ellipsis;
    max-width: 250px;
    text-wrap: nowrap;
    overflow: hidden;
    margin-bottom: 0;
  }
`;

export const IconError = styled(MinusCircleFilled)`
  color: ${({theme}) => theme.color.error};
`;

export const IconSuccess = styled(CheckCircleFilled)`
  color: ${({theme}) => theme.color.success};
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

export const DetailsWrapper = styled.div`
  align-items: center;
  cursor: pointer;
  justify-content: space-between;
  display: flex;
`;

export const ClearSearchIcon = styled(CloseCircleFilled)`
  position: absolute;
  right: 8px;
  top: 8px;
  color: ${({theme}) => theme.color.textLight};
  cursor: pointer;
`;

export const SearchContainer = styled(Row)`
  margin-bottom: 16px;
`;
