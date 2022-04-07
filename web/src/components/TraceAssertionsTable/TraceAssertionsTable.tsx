import {Table, Typography} from 'antd';
import {difference, sortBy} from 'lodash';
import {FC, useMemo} from 'react';
import {AssertionResult} from 'types';
import {getOperator} from 'utils';
import {getSpanSignature} from '../../services/SpanService';
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
  },
}) => {
  const selectorValueList = useMemo(() => selectors.map(({value}) => value), [selectors]);
  const parsedAssertionList = useMemo(() => {
    const spanAssertionList = spanListAssertionResult.reduce<Array<TParsedAssertion>>((list, {span, resultList}) => {
      const subResultList = resultList.map<TParsedAssertion>(
        ({propertyName, comparisonValue, operator, actualValue, hasPassed}) => ({
          spanLabels: difference(getSpanSignature(span), selectorValueList),
          key: propertyName,
          property: propertyName,
          comparison: operator,
          value: comparisonValue,
          actualValue,
          hasPassed,
        })
      );

      return list.concat(subResultList);
    }, []);

    return sortBy(spanAssertionList, ({spanLabels}) => spanLabels.join(''));
  }, [selectorValueList, spanListAssertionResult]);

  const spanCount = spanListAssertionResult.length;

  return (
    <S.AssertionsTableContainer>
      <S.AssertionsTableHeader>
        <Typography.Title level={4} style={{margin: 0}}>
          {selectors.map(({value}) => value).join(' ')}
        </Typography.Title>
        <Typography.Title level={4} style={{margin: 0}}>
          {`${spanCount} ${spanCount > 1 ? 'spans' : 'span'}`}
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
          render={(value: string[]) =>
            value.map(label => <S.AssertionsTableBadge count={label} key={label} />)
          }
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
