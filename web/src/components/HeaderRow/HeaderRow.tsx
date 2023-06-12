import AttributeActions from 'components/AttributeActions';
import {TSpanFlatAttribute} from 'types/Span.types';
import {THeader} from 'types/Test.types';
import Highlighted from '../Highlighted';
import * as S from './HeaderRow.styled';

interface IProps {
  header: THeader;
  onCreateTestOutput(attribute: TSpanFlatAttribute): void;
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
}

const HeaderRow = ({header: {key = '', value = ''}, onCreateTestOutput, onCreateTestSpec}: IProps) => {
  return (
    <S.HeaderContainer>
      <S.Header>
        <S.HeaderKey>{key}</S.HeaderKey>
        <S.HeaderValue>
          <Highlighted text={value} highlight="" />
        </S.HeaderValue>
      </S.Header>
      <AttributeActions
        attribute={{key: `tracetest.response.headers | json_path '$[?(@.Key=="${key}")].Value'`, value}}
        onCreateTestOutput={onCreateTestOutput}
        onCreateTestSpec={onCreateTestSpec}
      />
    </S.HeaderContainer>
  );
};

export default HeaderRow;
