import {MoreOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const TestCard = styled.div`
  display: grid;
  align-items: center;
  grid-template-columns: 20px 1fr 60px 2fr 220px 100px 20px;
  gap: 24px;
  padding: 24px;
  box-shadow: 0px 4px 8px rgba(153, 155, 168, 0.1);
  background: #fff;
  cursor: pointer;
`;

export const TextContainer = styled.div`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const ButtonContainer = styled.div`
  display: flex;
  justify-content: flex-end;
`;

export const NameText = styled(Typography.Text)`
  font-weight: 700;
  overflow-x: ellipsis;
`;

export const Text = styled(Typography.Text)``;

export const ActionButton = styled(MoreOutlined).attrs({
  style: {fontSize: 24, color: '#9AA3AB', cursor: 'pointer'},
})``;
