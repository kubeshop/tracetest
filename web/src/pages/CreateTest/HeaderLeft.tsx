import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import * as S from 'components/RunDetailLayout/RunDetailLayout.styled';
import {Overlay} from 'components/Inputs';
import {Form} from 'antd';

interface IProps {
  triggerType: string;
  origin: string;
}

const HeaderLeft = ({triggerType, origin}: IProps) => {
  const {navigate} = useDashboard();

  return (
    <S.Section $justifyContent="flex-start">
      <a data-cy="create-test-header-back-button" onClick={() => navigate(origin)}>
        <S.BackIcon />
      </a>
      <S.InfoContainer>
        <Form.Item name="name" noStyle>
          <Overlay />
        </Form.Item>
        <S.Row>
          <S.Text>{triggerType.toUpperCase()}</S.Text>
        </S.Row>
      </S.InfoContainer>
    </S.Section>
  );
};

export default HeaderLeft;
