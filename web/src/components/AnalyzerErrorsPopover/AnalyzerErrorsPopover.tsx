import {Popover} from 'antd';
import {TAnalyzerError} from 'types/TestRun.types';
import * as S from './AnalyzerErrorsPopover.styled';
import Content from './Content';

interface IProps {
  errors: TAnalyzerError[];
}

const AnalyzerErrorsPopover = ({errors}: IProps) => (
  <S.Container>
    <Popover
      content={<Content errors={errors} />}
      placement="right"
      title={<S.Title level={3}>Analyzer errors</S.Title>}
    >
      <S.ErrorIcon />
    </Popover>
  </S.Container>
);

export default AnalyzerErrorsPopover;
