import {Row} from 'antd';
import styled from 'styled-components';

export const ActionContainer = styled.div`
  align-items: center;
  display: flex;
  justify-content: center;
`;

export const Container = styled(Row)`
  height: calc(100% - 130px);
`;

export const Content = styled.div`
  margin-bottom: 24px;
  text-align: center;
`;

export const Icon = styled.img`
  margin-bottom: 25px;
  height: auto;
  width: 140px;
`;
