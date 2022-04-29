import {Button, Table, Typography} from 'antd';
import {useMemo, useState} from 'react';
import AssertionTableAnalyticsService from '../../services/Analytics/AssertionTableAnalytics.service';
import {IAssertion, ISpanAssertionResult} from '../../types/Assertion.types';
import OperatorService from '../../services/Operator.service';
import {ISpan} from '../../types/Span.types';
import {ITrace} from '../../types/Trace.types';
import CreateAssertionModal from '../CreateAssertionModal';
import CustomTable from '../CustomTable';
import * as S from './AssertionsTable.styled';

interface IAssertionsResultTableProps {
  assertionResults: ISpanAssertionResult[];
  assertion: IAssertion;
  sort: number;
  span: ISpan;
  testId: string;
  trace: ITrace;
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
  trace,
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);

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
    <S.AssertionsTableContainer>
      <S.AssertionsTableHeader>
        <Typography.Title level={5} style={{margin: 0}}>
          Assertion #{sort}
          {selectors.map(({value, propertyName}) => (
            <S.AssertionsTableBadge count={value} key={propertyName} />
          ))}
        </Typography.Title>
        <Button
          type="link"
          onClick={() => {
            AssertionTableAnalyticsService.onEditAssertionButtonClick(assertion.assertionId);
            setIsModalOpen(true);
          }}
        >
          Edit
        </Button>
      </S.AssertionsTableHeader>
      <CustomTable size="small" pagination={{hideOnSinglePage: true}} dataSource={parsedAssertionList}>
        <Table.Column title="Property" dataIndex="property" key="property" ellipsis width="50%" />
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
      <CreateAssertionModal
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        assertion={assertion}
        span={span}
        testId={testId}
        trace={trace}
      />
    </S.AssertionsTableContainer>
  );
};

export default AssertionsResultTable;
