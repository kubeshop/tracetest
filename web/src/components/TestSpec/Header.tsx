import {Tooltip} from 'antd';
import {singularOrPlural} from 'utils/Common';
import {SupportedEditors} from 'constants/Editor.constants';
import {Editor} from 'components/Inputs';
import * as S from './TestSpec.styled';

interface IProps {
  affectedSpans: number;
  assertionsFailed: number;
  assertionsPassed: number;
  hasError: boolean;
  selector: string;
  title: string;
}

const Header = ({affectedSpans, assertionsFailed, assertionsPassed, hasError, selector, title}: IProps) => {
  return (
    <S.Column>
      {title ? (
        <S.Title>{title}</S.Title>
      ) : (
        <Editor
          type={SupportedEditors.Selector}
          editable={false}
          basicSetup={{lineNumbers: false}}
          placeholder="Selecting All Spans"
          value={selector}
        />
      )}

      <div>
        {!!assertionsPassed && (
          <S.HeaderDetail>
            <S.HeaderDot $passed />
            {assertionsPassed}
          </S.HeaderDetail>
        )}
        {!!assertionsFailed && (
          <S.HeaderDetail>
            <S.HeaderDot $passed={false} />
            {assertionsFailed}
          </S.HeaderDetail>
        )}
        <S.HeaderDetail>
          <S.HeaderSpansIcon />
          {`${affectedSpans} ${singularOrPlural('span', affectedSpans)}`}
        </S.HeaderDetail>
        {hasError && (
          <span>
            <Tooltip title="This spec has errors">
              <S.WarningIcon />
            </Tooltip>
          </span>
        )}
      </div>
    </S.Column>
  );
};

export default Header;
