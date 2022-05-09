import {Table, Typography} from 'antd';
import {useStore} from 'react-flow-renderer';
import {difference, sortBy} from 'lodash';
import {FC, useCallback, useMemo} from 'react';
import CustomTable from '../CustomTable';
import * as S from './TraceAssertionsTable.styled';
import TraceAssertionTableAnalyticsService from '../../services/Analytics/TraceAssertionTableAnalytics.service';
import {IAssertionResult} from '../../types/Assertion.types';
import OperatorService from '../../services/Operator.service';

const {onSpanAssertionClick} = TraceAssertionTableAnalyticsService;

interface IProps {
  assertionResult: IAssertionResult;
  onSpanSelected(spanId: string): void;
}

type TParsedAssertion = {
  key: string;
  spanLabels: string[];
  property: string;
  comparison: string;
  value: string;
  actualValue: string;
  hasPassed: boolean;
  spanId: string;
};

const TraceAssertionsResultTable: FC<IProps> = ({
  assertionResult: {
    assertion: {selectors = []},
    spanListAssertionResult,
  },
  onSpanSelected,
}) => {
  const selectorValueList = useMemo(() => selectors.map(({value}) => value), [selectors]);
  const parsedAssertionList = useMemo(() => {
    const spanAssertionList = spanListAssertionResult.reduce<Array<TParsedAssertion>>((list, {resultList, span}) => {
      const subResultList = resultList.map<TParsedAssertion>(
        ({propertyName, comparisonValue, operator, actualValue, hasPassed, spanId}) => {
          const spanLabelList = span.signature.map(({value}) => value).concat([`#${spanId.slice(-4)}`]) || [];

          return {
            spanLabels: difference(spanLabelList, selectorValueList),
            spanId,
            key: `${propertyName}-${spanId}`,
            property: propertyName,
            comparison: operator,
            value: comparisonValue,
            actualValue,
            hasPassed,
          };
        }
      );

      return list.concat(subResultList);
    }, []);

    return sortBy(spanAssertionList, ({spanLabels}) => spanLabels.join(''));
  }, [selectorValueList, spanListAssertionResult]);

  const spanCount = spanListAssertionResult.length;
  const store = useStore();

  const getIsSelected = useCallback(
    (spanId: string): boolean => {
      const {selectedElements} = store.getState();
      const found = selectedElements ? selectedElements.find(({id}) => id === spanId) : undefined;

      return Boolean(found);
    },
    [store]
  );

  return (
    <S.Container data-cy="test-results-assertion-table">
      <S.Header>
        <Typography.Title level={5} style={{margin: 0}}>
          {selectors.map(({value}) => value).join(' ')}
        </Typography.Title>
        <Typography.Title level={5} style={{margin: 0}}>
          {`${spanCount} ${spanCount > 1 ? 'spans' : 'span'}`}
        </Typography.Title>
      </S.Header>
      <CustomTable
        size="small"
        pagination={false}
        dataSource={parsedAssertionList}
        onRow={record => ({
          onClick: () => {
            const spanId = (record as TParsedAssertion).spanId;

            onSpanAssertionClick(spanId);
            onSpanSelected(spanId);
          },
        })}
      >
        <Table.Column
          title="Span Labels"
          dataIndex="spanLabels"
          key="spanLabels"
          ellipsis
          width="40%"
          render={(value: string[], record: TParsedAssertion) =>
            (getIsSelected(record.spanId) ? [<S.SelectedLabelBadge count="selected" key="selected" />] : []).concat(
              value
                // eslint-disable-next-line react/no-array-index-key
                .map((label, index) => <S.LabelBadge count={label} key={`${label}-${index}`} />)
            )
          }
        />
        <Table.Column title="Property" dataIndex="property" key="property" ellipsis width="25%" />
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
          render={(value, record: TParsedAssertion) => (
            <Typography.Text strong type={record.hasPassed ? 'success' : 'danger'}>
              {value}
            </Typography.Text>
          )}
        />
      </CustomTable>
    </S.Container>
  );
};

export default TraceAssertionsResultTable;
