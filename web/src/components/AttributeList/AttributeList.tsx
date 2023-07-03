import {IPropsAttributeRow} from 'components/SpanDetail/SpanDetail';
import {OtelReference} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import TestRunOutput from 'models/TestRunOutput.model';
import {TSpanFlatAttribute} from 'types/Span.types';
import {TTestSpecSummary} from 'types/TestRun.types';
import * as S from './AttributeList.styled';
import EmptyAttributeList from './EmptyAttributeList';

interface IProps {
  attributeList: TSpanFlatAttribute[];
  searchText?: string;
  semanticConventions: OtelReference;
  testSpecs?: TTestSpecSummary;
  testOutputs?: TestRunOutput[];
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  onCreateOutput(attribute: TSpanFlatAttribute): void;
  AttributeRowComponent: React.ComponentType<IPropsAttributeRow>;
}

const AttributeList = ({
  attributeList,
  searchText,
  semanticConventions,
  testSpecs,
  testOutputs,
  onCreateTestSpec,
  onCreateOutput,
  AttributeRowComponent,
}: IProps) => {
  return attributeList.length ? (
    <S.AttributeList data-cy="attribute-list">
      {attributeList.map(attribute => (
        <AttributeRowComponent
          key={attribute.key}
          attribute={attribute}
          searchText={searchText}
          semanticConventions={semanticConventions}
          testSpecs={testSpecs}
          testOutputs={testOutputs}
          onCreateTestSpec={onCreateTestSpec}
          onCreateOutput={onCreateOutput}
        />
      ))}
    </S.AttributeList>
  ) : (
    <EmptyAttributeList />
  );
};

export default AttributeList;
