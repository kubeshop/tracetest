import {Col, Typography} from 'antd';

import icon from 'assets/empty.svg';
import * as S from './Empty.styled';

interface IProps {
  message?: string;
}

const Empty = ({message = ''}: IProps) => (
  <S.Container align="middle">
    <Col span={12} offset={6}>
      <S.Content>
        <div>
          <S.Icon alt="empty" src={icon} />
        </div>
        <Typography.Text>{message}</Typography.Text>
      </S.Content>
    </Col>
  </S.Container>
);

export default Empty;
