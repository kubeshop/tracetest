import {CheckCircleFilled, InfoCircleFilled, MinusCircleFilled, MoreOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const ActionButton = styled(MoreOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: pointer;
  font-size: ${({theme}) => theme.size.lg};
`;

export const Container = styled.div<{$isWhite?: boolean}>`
  align-items: center;
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  border-radius: 2px;
  background: ${({$isWhite, theme}) => ($isWhite ? theme.color.white : theme.color.background)};
  display: flex;
  gap: 16px;
  padding: 8px 16px;
`;

export const HeaderDetail = styled(Typography.Text)`
  display: flex;
  align-items: center;
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.sm};
  margin-right: 8px;
`;

export const HeaderDot = styled.span<{$passed: boolean}>`
  background-color: ${({$passed, theme}) => ($passed ? theme.color.success : theme.color.error)};
  border-radius: 50%;
  display: inline-block;
  height: 10px;
  line-height: 0;
  margin-right: 4px;
  vertical-align: -0.1em;
  width: 10px;
`;

export const IconContainer = styled.div`
  width: 14px;
`;

export const IconSuccess = styled(CheckCircleFilled)`
  color: ${({theme}) => theme.color.success};
`;

export const IconFail = styled(MinusCircleFilled)`
  color: ${({theme}) => theme.color.error};
`;

export const IconInfo = styled(InfoCircleFilled)`
  color: ${({theme}) => theme.color.textLight};
`;

export const Info = styled.div`
  flex: 1;
`;

export const Text = styled(Typography.Text).attrs({
  as: 'p',
})<{$hasLink?: boolean}>`
  color: ${({$hasLink, theme}) => ($hasLink ? theme.color.primary : theme.color.textSecondary)};
  font-size: ${({theme}) => theme.size.sm};
  margin: 0;
`;

export const Title = styled(Typography.Title).attrs({level: 3})`
  && {
    margin: 0;
  }
`;

export const Row = styled.div`
  display: flex;
`;
