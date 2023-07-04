import {Typography} from 'antd';
import styled from 'styled-components';

export const DocsBannerContainer = styled.div`
  display: flex;
  gap: 8px;
  align-items: center;
  padding: 12px 18px;
  border-radius: 2px;
  width: max-content;
  background: ${({theme}) => theme.color.backgroundInteractive};
`;

export const Text = styled(Typography.Text)`
  a {
    font-weight: 600;
  }
`;
