import {Typography} from 'antd';
import {withCustomization} from 'providers/Customization';
import * as S from './Typescript.styled';

const Typescript = () => (
  <S.Container>
    <S.TitleContainer>
      <S.Title>Typescript Integration</S.Title>
    </S.TitleContainer>
    <Typography.Paragraph>* this capability is limited to the commercial version of Tracetest.</Typography.Paragraph>
  </S.Container>
);

export default withCustomization(Typescript, 'typescript');
