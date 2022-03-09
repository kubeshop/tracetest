import {Button} from 'antd';
import styled from 'styled-components';

export const Header = styled.div`
  height: 64px;
  padding: 16px 32px;
  border-bottom: 1px solid rgb(213, 215, 224);
`;

export const Content = styled.main`
  display: flex;
  padding: 16px 0;
`;

export const SideMenu = styled.div`
  display: flex;
  flex-direction: column;
  flex: 0.2;
  padding: 0 16px;
`;

export const TestsContainer = styled.div`
  flex: 0.8;
  padding: 0 8px;
`;

export const CreateTestButton = styled(Button)`
  border-radius: 16px;
  align-self: flex-start;
`;
