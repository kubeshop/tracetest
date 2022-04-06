import {Button, Table, Typography} from 'antd';
import {useMemo} from 'react';
import {Assertion, SpanAssertionResult} from 'types';
import {getOperator} from 'utils';
import CustomTable from '../CustomTable';
import * as S from './AssertionsTable.styled';

interface IProps {
  assertionResults: SpanAssertionResult[];
  assertion: Assertion;
  sort: number;
}

type TParsedAssertion = {
  key: string;
  property: string;
  comparison: string;
  value: string;
  actualValue: string;
  hasPassed: boolean;
};

const AssertionsResultTable = ({assertionResults, assertion: {selectors = []}, sort}: IProps) => {
  const parsedAssertionList = useMemo<Array<TParsedAssertion>>(
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
        <Typography.Title level={3} style={{margin: 0}}>
          Assertion #{sort}
          {selectors.map(({value, propertyName}) => (
            <S.AssertionsTableBadge count={value} key={propertyName} />
          ))}
        </Typography.Title>
        <Button type="link">Edit</Button>
      </S.AssertionsTableHeader>
      <CustomTable
        size="small"
        pagination={{hideOnSinglePage: true}}
        dataSource={parsedAssertionList}
        bordered
        tableLayout="fixed"
      >
        <Table.Column title="Property" dataIndex="property" key="property" ellipsis width="60%" />
        <Table.Column title="Comparison" dataIndex="comparison" key="comparison" render={value => getOperator(value)} />
        <Table.Column title="Value" dataIndex="value" key="value" />
        <Table.Column
          title="Actual"
          dataIndex="actualValue"
          key="actualValue"
          render={(value, record: TParsedAssertion) => (
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
