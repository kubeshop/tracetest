import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {ISpanResult} from 'types/TestSpecs.types';
import * as S from './SpanResultDetail.styled';
import Content from './Content';

interface IProps {
  isOpen: boolean;
  onClose(): void;
  spanResult?: ISpanResult;
}

const SpanResultDetail = ({isOpen, onClose, spanResult}: IProps) => {
  const {
    run: {trace},
  } = useTestRun();

  if (!spanResult) return null;

  const span = trace?.flat[spanResult.spanId] || {};

  return (
    <S.DrawerContainer
      closable={false}
      getContainer={false}
      mask={false}
      onClose={onClose}
      placement="right"
      visible={isOpen}
      width="100%"
      height="100%"
      $type={span.type}
    >
      {!!spanResult && <Content onClose={onClose} span={span} checkResults={spanResult.checkResults} />}
    </S.DrawerContainer>
  );
};

export default SpanResultDetail;
