import styled from 'styled-components';
import Layout from 'antd/lib/layout';

export const Content = styled(Layout.Content)`
  background: ${({theme}) => theme.color.bg};
  display: flex;
  flex-direction: column;
  flex-grow: 1;
`;
