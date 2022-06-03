import {Divider as AntdDivider, Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  padding-left: 16px;
  padding-right: 16px;
`;
export const TitleContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
`;

export const Header = styled.div`
  display: flex;
  justify-content: space-between;
  padding-left: 16px;
  padding-right: 16px;
  padding-bottom: 0;
  margin-top: 16px;
  height: 24px;
`;
export const Divider = styled(AntdDivider)`
  margin-top: 12px;
  margin-bottom: 12px;
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 0;
    line-height: 24px;
    font-size: 14px;
  }
`;
