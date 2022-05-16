import {Button, Table, Typography} from 'antd';
import {useMemo} from 'react';
import AssertionTableAnalyticsService from '../../services/Analytics/AssertionTableAnalytics.service';
import {IAssertion, ISpanAssertionResult} from '../../types/Assertion.types';
import OperatorService from '../../services/Operator.service';
import {ISpan} from '../../types/Span.types';
import CustomTable from '../CustomTable';
import * as S from './AssertionsTable.styled';
import {useCreateAssertionModal} from '../CreateAssertionModal/CreateAssertionModalProvider';

interface IAssertionsResultTableProps {
  assertionResults: ISpanAssertionResult[];
  assertion: IAssertion;
  sort: number;
  span: ISpan;
  testId: string;
  resultId: string;
}

interface IParsedAssertion {
  key: string;
  property: string;
  comparison: string;
  value: string;
  actualValue: string;
  hasPassed: boolean;
}

const AssertionsResultTable: React.FC<IAssertionsResultTableProps> = ({
  assertionResults,
  assertion: {selectors = []},
  assertion,
  sort,
  span,
  testId,
  resultId,
}) => {
  const {open} = useCreateAssertionModal();

  const parsedAssertionList = useMemo<Array<IParsedAssertion>>(
    () =>
      assertionResults.map(({propertyName, comparisonValue, operator, actualValue, hasPassed}) => ({
        key: propertyName,
        property: propertyName,
        comparison: operator,
        value: comparisonValue,
        actualValue,
        hasPassed,
      })),
    [assertionResults]
  );

  return (
    <S.AssertionsTableContainer data-cy="assertion-table">
      <S.AssertionsTableHeader>
        <Typography.Title level={5} style={{margin: 0}}>
          Assertion #{sort}
          {selectors.map(({value, propertyName}) => (
            <S.AssertionsTableBadge count={value} key={propertyName} />
          ))}
        </Typography.Title>
        <Button
          type="link"
          data-cy="edit-assertion-button"
          onClick={() => {
            AssertionTableAnalyticsService.onEditAssertionButtonClick(assertion.assertionId);
            open({
              span,
              assertion,
              resultId,
              testId,
            });
          }}
        >
          Edit
        </Button>
      </S.AssertionsTableHeader>
      <CustomTable size="small" pagination={{hideOnSinglePage: true}} dataSource={parsedAssertionList}>
        <Table.Column
          title="Property"
          dataIndex="property"
          key="property"
          ellipsis
          width="50%"
          render={value => <span data-cy="assertion-check-property">{value}</span>}
        />
        <Table.Column
          title="Comparison"
          dataIndex="comparison"
          key="comparison"
          render={value => OperatorService.getOperatorName(value)}
        />
        <Table.Column title="Value" dataIndex="value" key="value" />
        <Table.Column
          title="Actual"
          dataIndex="actualValue"
          key="actualValue"
          render={(value, record: IParsedAssertion) => (
            <Typography.Text strong type={record.hasPassed ? 'success' : 'danger'}>
              {value}
            </Typography.Text>
          )}
        />
      </CustomTable>
    </S.AssertionsTableContainer>
  );
};

export default AssertionsResultTable;
