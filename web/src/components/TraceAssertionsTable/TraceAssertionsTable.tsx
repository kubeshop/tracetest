import {Table, Typography} from 'antd';
import {FC, useMemo} from 'react';
import {AssertionResult} from 'types';
import {getOperator} from 'utils';
import CustomTable from '../CustomTable';
import * as S from './TraceAssertionsTable.styled';

interface IProps {
  assertionResult: AssertionResult;
}

type TParsedAssertion = {
  key: string;
  spanLabels: string[];
  property: string;
  comparison: string;
  value: string;
  actualValue: string;
  hasPassed: boolean;
};

const TraceAssertionsResultTable: FC<IProps> = ({
  assertionResult: {
    assertion: {selectors = []},
    spanListAssertionResult,
    spanCount,
  },
}) => {
  const selectorsList = useMemo(() => selectors?.map(({value}) => value), [selectors]);
  const parsedAssertionList = useMemo(
    () =>
      spanListAssertionResult.reduce<Array<TParsedAssertion>>((list, currentResultList) => {
        const subResultList = currentResultList.map<TParsedAssertion>(
          ({propertyName, comparisonValue, operator, actualValue, hasPassed}) => ({
            spanLabels: selectorsList,
            key: propertyName,
            property: propertyName,
            comparison: operator,
            value: comparisonValue,
            actualValue,
            hasPassed,
          })
        );

        return list.concat(subResultList);
      }, []),
    [selectorsList, spanListAssertionResult]
  );

  return (
    <S.AssertionsTableContainer>
      <S.AssertionsTableHeader>
        <Typography.Title level={4} style={{margin: 0}}>
          {selectors.map(({value}) => value).join(' ')}
        </Typography.Title>
        <Typography.Title level={4} style={{margin: 0}}>
          {spanCount} spans
        </Typography.Title>
      </S.AssertionsTableHeader>
      <CustomTable
        size="small"
        pagination={{hideOnSinglePage: true}}
        dataSource={parsedAssertionList}
        bordered
        tableLayout="fixed"
      >
        <Table.Column
          title="Span Labels"
          dataIndex="spanLabels"
          key="spanLabels"
          ellipsis
          width="40%"
          render={(value: string[], record, index) => {
            const spanLabelText = value.join('');
            const obj = {
              children: value.map(label => <S.AssertionsTableBadge count={label} key={label} />),
              props: {rowSpan: 1},
            };

            if (parsedAssertionList.filter(({spanLabels}) => spanLabels.join('') === spanLabelText).length === 1) {
              return obj;
            }

            if (parsedAssertionList.findIndex(({spanLabels}) => spanLabels.join('') === spanLabelText) === index) {
              const count = parsedAssertionList.filter(({spanLabels}) => spanLabels.join('') === spanLabelText).length;
              obj.props.rowSpan = count;

              return obj;
            }

            obj.props.rowSpan = 0;
            return obj;
          }}
        />
        <Table.Column title="Property" dataIndex="property" key="property" ellipsis width="25%" />
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

export default TraceAssertionsResultTable;
