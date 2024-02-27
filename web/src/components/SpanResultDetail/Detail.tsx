import Span from 'models/Span.model';
import {ClockCircleOutlined, SettingOutlined, ToolOutlined} from '@ant-design/icons';
import {Tooltip} from 'antd';
import * as STestSpec from 'components/TestSpec/TestSpec.styled';
import {SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import * as SSpanNode from 'components/Visualization/components/DAG/BaseSpanNode/BaseSpanNode.styled';
import * as S from './SpanResultDetail.styled';

interface IProps {
  assertionsFailed: number;
  assertionsPassed: number;
  span: Span;
}

const Detail = ({assertionsFailed, assertionsPassed, span: {duration, name, service, kind, system, type}}: IProps) => {
  return (
    <>
      <S.DetailsWrapper>
        <S.SpanHeaderContainer>
          <SSpanNode.BadgeType count={SemanticGroupNamesToText[type]} $type={type} />
          <Tooltip title={name}>
            <S.HeaderTitle level={3}>{name}</S.HeaderTitle>
          </Tooltip>
        </S.SpanHeaderContainer>
        <div>
          {Boolean(assertionsPassed) && (
            <STestSpec.HeaderDetail>
              <STestSpec.HeaderDot $passed />
              {assertionsPassed}
            </STestSpec.HeaderDetail>
          )}
          {Boolean(assertionsFailed) && (
            <STestSpec.HeaderDetail>
              <STestSpec.HeaderDot $passed={false} />
              {assertionsFailed}
            </STestSpec.HeaderDetail>
          )}
        </div>
      </S.DetailsWrapper>

      <S.SpanHeaderContainer>
        <S.HeaderItem>
          <SettingOutlined />
          <S.HeaderItemText>{`${service} ${kind}`}</S.HeaderItemText>
        </S.HeaderItem>
        {Boolean(system) && (
          <S.HeaderItem>
            <ToolOutlined />
            <S.HeaderItemText>{system}</S.HeaderItemText>
          </S.HeaderItem>
        )}

        <S.HeaderItem>
          <ClockCircleOutlined />
          <S.HeaderItemText>{duration}</S.HeaderItemText>
        </S.HeaderItem>
      </S.SpanHeaderContainer>
    </>
  );
};

export default Detail;
