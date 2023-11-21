import {Overlay} from 'components/Inputs';
import * as S from 'components/RunDetailLayout/RunDetailLayout.styled';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {Form} from 'antd';

interface IProps {
  triggerType: string;
  origin: string;
}

const HeaderLeft = ({triggerType, origin}: IProps) => {
  const {navigate} = useDashboard();

  return (
    <S.Section $justifyContent="flex-start">
      <a onClick={() => navigate(origin)}>
        <S.BackIcon />
      </a>
      <S.InfoContainer>
        <S.Row $height={24}>
          <Form.Item name="name" noStyle>
            <Overlay />
          </Form.Item>
        </S.Row>
        <S.Text>{triggerType.toUpperCase()}</S.Text>
      </S.InfoContainer>
    </S.Section>
  );
};

export default HeaderLeft;
