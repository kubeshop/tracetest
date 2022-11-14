import {Tooltip} from 'antd';

import {TSummary} from 'types/Test.types';
import Date from 'utils/Date';
import * as S from './ResourceCard.styled';

interface IProps {
  summary?: TSummary;
}

const ResourceCardSummary = ({summary}: IProps) => (
  <>
    <div>
      <S.Text>Last run time:</S.Text>
      <Tooltip title={Date.format(summary?.lastRun?.time ?? '')}>
        <S.Text>{Date.getTimeAgo(summary?.lastRun?.time ?? '')}</S.Text>
      </Tooltip>
    </div>

    <div>
      <S.Text>Last run result:</S.Text>
      <S.Row>
        <Tooltip title="Passed assertions">
          <S.HeaderDetail>
            <S.HeaderDot $passed />
            {summary?.lastRun?.passes}
          </S.HeaderDetail>
        </Tooltip>
        <Tooltip title="Failed assertions">
          <S.HeaderDetail>
            <S.HeaderDot $passed={false} />
            {summary?.lastRun?.fails}
          </S.HeaderDetail>
        </Tooltip>
      </S.Row>
    </div>
  </>
);

export default ResourceCardSummary;
