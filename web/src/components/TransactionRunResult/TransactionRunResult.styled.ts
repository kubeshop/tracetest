import {CheckCircleFilled, CloseCircleFilled} from '@ant-design/icons';
import {Typography, Tag} from 'antd';
import {Link} from 'react-router-dom';
import styled from 'styled-components';

export const Container = styled.div`
  align-items: center;
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  border-radius: 2px;
  background: ${({theme}) => theme.color.background};
  display: grid;
  grid-template-columns: auto 1fr auto auto;
  gap: 16px;
  padding: 7px 16px;
  margin-bottom: 8px;
  height: 58px;
`;

export const Info = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
`;

export const Stack = styled.div`
  display: flex;
  justify-content: space-between;
  flex-direction: column;
`;

export const Text = styled(Typography.Text)<{opacity?: number; fontWeight?: string}>`
  && {
    font-size: ${({theme}) => theme.size.sm};
    color: ${({opacity}) => `rgba(3, 24, 73, ${opacity || 1})`};
    margin: 0 !important;
    font-weight: ${({fontWeight}) => fontWeight || 'normal'};
  }
`;

export const ExecutionStepName = styled(Typography.Text)`
  && {
    font-size: ${({theme}) => theme.size.sm};
  }
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.lg};
    margin-bottom: 24px;
  }
`;

export const IconSuccess = styled(CheckCircleFilled)`
  color: ${({theme}) => theme.color.success};
`;

export const IconFail = styled(CloseCircleFilled)`
  color: ${({theme}) => theme.color.error};
`;

export const TagContainer = styled.div`
  display: grid;
  grid-template-columns: fit-content(8ch) fit-content(8ch) fit-content(8ch);

  > span:nth-child(1) {
    border-radius: 2px 0px 0px 2px;
  }

  > span:last-child {
    border-radius: 0px 2px 2px 0px;
  }
`;

export const TextTag = styled(Tag)<{$isLight?: boolean}>`
  && {
    margin: 0;
    border-radius: 0px;
    border: none;
    background: ${({theme, $isLight}) => ($isLight ? 'rgba(3, 24, 73, 0.05)' : theme.color.borderLight)};
  }
`;

export const ExecutionStepStatus = styled.div`
  color: ${({theme}) => theme.color.textLight};
  font-weight: 700;
`;

export const ExecutionStepRunLink = styled(Link)`
  && {
    color: ${({theme}) => theme.color.textLight};
  }
`;

export const VariableContainer = styled(Container)`
  grid-template-columns: 1fr 1fr;
`;

export const ResultContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 24px;
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

export const AssertionResultContainer = styled.div`
  display: flex;
  align-items: center;
`;
