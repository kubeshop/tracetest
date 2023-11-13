import styled from 'styled-components';
import {Form, Typography} from 'antd';
import {InfoCircleOutlined} from '@ant-design/icons';

export const TestsContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

export const InputContainer = styled.div`
  display: grid;
  gap: 24px;
  grid-template-columns: 1fr 1fr;
`;

export const DetailContainer = styled.div`
  overflow-wrap: break-word;
  width: 270px;
`;

export const TestList = styled.ul`
  list-style-type: none;
  margin: 0;
  padding-left: 0;
`;

export const ToolTipTitle = styled(Typography.Text)`
  && {
    margin: 0;
    font-weight: 600;
  }
`;

export const InfoIcon = styled(InfoCircleOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: pointer;
  margin: 4px;
`;

export const FromItem = styled(Form.Item)<{$hasValue: boolean}>`
  margin: 0;

  input {
    border: 1px solid ${({theme, $hasValue}) => !$hasValue && theme.color.error};
  }
`;
