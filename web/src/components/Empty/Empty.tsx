import {Col, Typography} from 'antd';

import icon from 'assets/empty.svg';
import * as S from './Empty.styled';

interface IProps {
  message?: React.ReactNode;
  title?: string;
  action?: React.ReactNode;
}

const Empty = ({message = '', title = '', action}: IProps) => (
  <S.Container align="middle">
    <Col lg={{span: 10, offset: 7}}>
      <S.Content>
        <div>
          <S.Icon alt="empty" src={icon} />
        </div>
        <Typography.Title>{title}</Typography.Title>
        <Typography.Text>{message}</Typography.Text>
      </S.Content>
      <S.ActionContainer>{action}</S.ActionContainer>
    </Col>
  </S.Container>
);

export default Empty;
