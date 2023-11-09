import styled from 'styled-components';
import {Typography} from 'antd';
import {CheckCircleOutlined} from '@ant-design/icons';

const defaultHeight = '100vh - 106px - 60px - 40px';

export const FormContainer = styled.div`
  display: grid;
  grid-template-columns: auto 1fr;
  height: calc(${defaultHeight} - 48px);
  overflow: hidden;
`;

export const FactoryContainer = styled.div`
  display: flex;
  flex-direction: column;
  padding: 22px 0;
  justify-content: space-between;
  height: calc(${defaultHeight} - 25px);
  overflow-y: scroll;

  .ant-form-item {
    margin: 0;
  }
`;

export const TopContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0 22px;
`;

export const Circle = styled.div`
  min-height: 16px;
  min-width: 16px;
  max-height: 16px;
  max-width: 16px;
  border: ${({theme}) => `1px solid ${theme.color.primary}`};
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

export const Check = styled.div`
  height: 8px;
  width: 8px;
  background: ${({theme}) => theme.color.primary};
  border-radius: 50%;
  display: inline-block;
`;

export const DataStoreName = styled(Typography.Text)<{$isSelected: boolean}>`
  && {
    color: ${({theme, $isSelected}) => ($isSelected ? theme.color.primary : theme.color.text)};
    font-weight: ${({$isSelected}) => ($isSelected ? 700 : 400)};
  }
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.md};
    margin: 0 !important;
  }
`;

export const Description = styled(Typography.Text)`
  && {
    color: ${({theme}) => theme.color.textSecondary};
  }
`;

export const ButtonsContainer = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 8px;
  margin-top: 23px;
  padding: 16px 22px;
  border-top: 1px solid ${({theme}) => theme.color.borderLight};
  position: sticky;
  bottom: 0;
  background: white;
`;

export const SaveContainer = styled.div`
  display: flex;
  gap: 8px;
`;

export const InfoIcon = styled(CheckCircleOutlined)`
  color: ${({theme}) => theme.color.text};
  cursor: pointer;
  margin: 4px;
`;
