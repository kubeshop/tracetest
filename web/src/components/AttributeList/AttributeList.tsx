import AttributeRow from 'components/AttributeRow';
import {OtelReference} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import TestOutput from 'models/TestOutput.model';
import * as S from './AttributeList.styled';
import EmptyAttributeList from './EmptyAttributeList';

interface IProps {
  assertions?: TResultAssertions;
  attributeList: TSpanFlatAttribute[];
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  searchText?: string;
  semanticConventions: OtelReference;
  onCreateOutput(attribute: TSpanFlatAttribute): void;
  outputs: TestOutput[];
}

const AttributeList = ({
  assertions,
  attributeList,
  onCreateTestSpec,
  onCreateOutput,
  searchText,
  semanticConventions,
  outputs,
}: IProps) => {
  return attributeList.length ? (
    <S.AttributeList data-cy="attribute-list">
      {attributeList.map(attribute => (
        <AttributeRow
          searchText={searchText}
          assertions={assertions}
          attribute={attribute}
          key={attribute.key}
          onCreateTestSpec={onCreateTestSpec}
          onCreateOutput={onCreateOutput}
          semanticConventions={semanticConventions}
          outputs={outputs}
        />
      ))}
    </S.AttributeList>
  ) : (
    <EmptyAttributeList />
  );
};

export default AttributeList;
