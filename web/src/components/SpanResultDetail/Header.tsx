import Span from 'models/Span.model';
import {ArrowLeftOutlined} from '@ant-design/icons';
import {Button, Divider} from 'antd';
import * as S from './SpanResultDetail.styled';
import Detail from './Detail';

interface IProps {
  assertionsFailed: number;
  assertionsPassed: number;
  onClose(): void;
  span: Span;
}

const Header = ({span, assertionsFailed, assertionsPassed, onClose}: IProps) => (
  <>
    <S.HeaderContainer>
      <S.Row $hasGap>
        <Button icon={<ArrowLeftOutlined />} onClick={onClose} type="link" />
        <S.HeaderTitle level={2}>Span Result Detail</S.HeaderTitle>
      </S.Row>
    </S.HeaderContainer>
    <Divider />
    <Detail span={span} assertionsFailed={assertionsFailed} assertionsPassed={assertionsPassed} />
    <Divider />
  </>
);

export default Header;
