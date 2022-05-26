import {MoreOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const ResultCard = styled.div`
  display: grid;
  align-items: center;
  grid-template-columns: 300px 300px 100px 220px 40px 40px 40px 1fr;
  gap: 16px;
  padding: 16px 12px;
  border: 1px solid rgba(3, 24, 73, 0.1);
  border-radius: 2px;
  background: #fbfbff;
  cursor: pointer;
`;

export const TextContainer = styled.div`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const Text = styled(Typography.Text)`
  overflow-x: ellipsis;
  font-size: 12px;
`;

export const ActionButton = styled(MoreOutlined).attrs({
  style: {fontSize: 24, color: '#9AA3AB', cursor: 'pointer'},
})``;
