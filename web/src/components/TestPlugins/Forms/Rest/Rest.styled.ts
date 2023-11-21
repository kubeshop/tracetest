import {Typography} from 'antd';
import styled from 'styled-components';

export const Row = styled.div`
  display: flex;
`;

export const Label = styled(Typography.Text).attrs({as: 'div'})`
  margin-bottom: 8px;
`;

export const SettingsContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 14px;
`;