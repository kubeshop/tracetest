import {Typography, Divider as AntdDivider} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  min-height: 100px;
  padding-top: 12px;
  padding-left: 16px;
  padding-right: 16px;
`;
export const TitleContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

export const Header = styled.div`
  display: flex;
  justify-content: space-between;
  padding: 16px;
  height: 48px;
  background: linear-gradient(180deg, #2f1e61 -11.46%, #bc334a 134.37%);
  > div {
    width: 100%;
  }
`;
export const Divider = styled(AntdDivider)`
  margin-top: 12px;
  margin-bottom: 0;
`;

export const TitleText = styled(Typography.Text).attrs({level: 3})`
  && {
    color: white;
    opacity: 0.6;
  }
`;
export const Title = styled(Typography.Title).attrs({level: 3})`
  && {
    color: white;
    margin-bottom: 0;
  }
`;
