import {Tooltip} from 'antd';
import Date from 'utils/Date';
import Summary from 'models/Summary.model';
import * as S from './ResourceCard.styled';

interface IProps {
  summary: Summary;
}

const ResourceCardSummary = ({
  summary: {
    lastRun: {time, passes, fails},
  },
}: IProps) => (
  <>
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
