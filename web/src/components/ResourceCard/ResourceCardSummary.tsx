import {Tooltip} from 'antd';
import AnalyzerScore from 'components/AnalyzerScore';
import Summary from 'models/Summary.model';
import Date from 'utils/Date';
import * as S from './ResourceCard.styled';

interface IProps {
  summary: Summary;
}

const ResourceCardSummary = ({
  summary: {
    lastRun: {time, passes, fails, analyzerScore},
  },
}: IProps) => (
  <>
    <div>
      {!!analyzerScore && (
        <Tooltip title="Trace Analyzer score">
          <div>
            <AnalyzerScore width="28px" height="28px" score={analyzerScore} />
          </div>
        </Tooltip>
      )}
    </div>
    <div>
      <S.Text>Last run result:</S.Text>
      <S.Row>
        <Tooltip title="Passed assertions">
          <S.HeaderDetail>
            <S.HeaderDot $passed />
            {passes}
          </S.HeaderDetail>
        </Tooltip>
        <Tooltip title="Failed assertions">
          <S.HeaderDetail>
            <S.HeaderDot $passed={false} />
            {fails}
          </S.HeaderDetail>
        </Tooltip>
      </S.Row>
    </div>
    <div>
      <S.Text>Last run time:</S.Text>
      <Tooltip title={Date.format(time ?? '')}>
        <S.Text>{Date.getTimeAgo(time ?? '')}</S.Text>
      </Tooltip>
    </div>
  </>
);

export default ResourceCardSummary;
