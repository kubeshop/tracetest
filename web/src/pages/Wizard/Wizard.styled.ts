import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  padding: 24px;

  .ant-collapse {
    background-color: ${({theme}) => theme.color.white};
    margin-bottom: 24px;
  }

  && {
    .ant-collapse-header {
      align-items: center;
    }
  }

  .ant-collapse-content > .ant-collapse-content-box {
    padding: 0;
  }
`;

export const Header = styled.div`
  margin-bottom: 16px;
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 6px;
  }
`;

export const Text = styled(Typography.Text)``;
