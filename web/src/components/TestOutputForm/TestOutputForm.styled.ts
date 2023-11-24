import {ReadOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  background-color: ${({theme}) => theme.color.white};
  height: 100%;
  overflow-y: auto;
  padding: 24px;
  position: relative;
`;

export const Title = styled(Typography.Title).attrs({level: 2})`
  && {
    margin-bottom: 16px;
  }
`;

export const InputContainer = styled.div`
  display: grid;
  gap: 24px;
  grid-template-columns: 100%;
`;

export const Footer = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 8px;
`;

export const FormSection = styled.div`
  margin-bottom: 24px;
`;

export const FormSectionRow = styled.div`
  margin-bottom: 8px;
`;

export const FormSectionRow1 = styled.div`
  align-items: center;
  display: flex;
  gap: 12px;
`;

export const FormSectionHeaderSelector = styled.div`
  align-items: center;
  display: flex;
  justify-content: space-between;
`;

export const FormSectionTitle = styled(Typography.Title).attrs({level: 3})<{$noMargin?: boolean}>`
  && {
    margin-bottom: ${({$noMargin}) => ($noMargin ? '0' : '4px')};
  }
`;

export const FormSectionText = styled(Typography.Text)`
  color: ${({theme}) => theme.color.textSecondary};
`;

export const ReadIcon = styled(ReadOutlined)`
  margin-top: 4px;
`;

export const SelectorTitleContainer = styled.div`
  align-items: center;
  display: flex;
  gap: 12px;
`;

export const SelectorLabel = styled(Typography.Text)`
  margin: 0;
`;
