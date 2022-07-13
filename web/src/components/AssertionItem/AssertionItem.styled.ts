import {ApartmentOutlined, CheckCircleFilled, MinusCircleFilled} from '@ant-design/icons';
import {Card, Collapse, Tag, Typography} from 'antd';
import styled from 'styled-components';

import {SemanticGroupNames, SemanticGroupNamesToColor} from 'constants/SemanticGroupNames.constants';

export const ActionTag = styled(Tag)`
  background-color: ${({theme}) => theme.color.primary};
  border-color: transparent;
  color: ${({theme}) => theme.color.white};
`;

export const AssertionCollapse = styled(Collapse)<{$isSelected: boolean}>`
  background-color: ${({theme}) => theme.color.white};
  border: ${({$isSelected, theme}) =>
    $isSelected ? `1px solid ${theme.color.text}` : `1px solid ${theme.color.border}`};
  border-bottom: ${({$isSelected, theme}) => ($isSelected ? `1px solid ${theme.color.text}` : 0)};

  > .ant-collapse-item > .ant-collapse-header {
    align-items: center;
  }

  .ant-collapse-extra {
    border-right: ${({theme}) => `1px solid ${theme.color.border}`};
    padding-right: 8px;
  }
`;

export const Column = styled.div`
  display: flex;
  flex-direction: column;
`;

export const GridContainer = styled.div`
  display: grid;
  column-gap: 24px;
  grid-template-columns: 30% 1fr;
  padding: 10px 0 10px 30px;
`;

export const HeaderDetail = styled(Typography.Text)`
  color: ${({theme}) => theme.color.text};
  font-size: ${({theme}) => theme.size.sm};
  margin-right: 8px;
`;

export const HeaderDot = styled.span<{$passed: boolean}>`
  background-color: ${({$passed, theme}) => ($passed ? theme.color.success : theme.color.error)};
  height: 10px;
  width: 10px;
  display: inline-block;
  margin-right: 4px;
  line-height: 0;
  vertical-align: -0.1em;
  border-radius: 50%;
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

export const HeaderSpansIcon = styled(ApartmentOutlined)`
  margin-right: 4px;
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

export const Row = styled.div<{$align?: string}>`
  align-items: ${({$align}) => $align || 'center'};
  display: flex;
`;

export const SecondaryText = styled(Typography.Text)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.sm};
`;

export const SpanCard = styled(Card)<{$isSelected: boolean; $type: SemanticGroupNames}>`
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

export const SpanHeaderContainer = styled.div`
  align-items: center;
  display: flex;
  gap: 8px;
`;
