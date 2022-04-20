import {Table, Typography} from 'antd';
import {useStore} from 'react-flow-renderer';
import {difference, sortBy} from 'lodash';
import {FC, useCallback, useMemo} from 'react';
import {AssertionResult, ITrace} from 'types';
import {getOperator} from 'utils';
import {getSpanSignature} from '../../services/SpanService';
import CustomTable from '../CustomTable';
import * as S from './TraceAssertionsTable.styled';

interface IProps {
  assertionResult: AssertionResult;
  trace: ITrace;
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
  trace,
  onSpanSelected,
}) => {
  const selectorValueList = useMemo(() => selectors.map(({value}) => value), [selectors]);
  const parsedAssertionList = useMemo(() => {
    const spanAssertionList = spanListAssertionResult.reduce<Array<TParsedAssertion>>((list, {resultList}) => {
      const subResultList = resultList.map<TParsedAssertion>(
        ({propertyName, comparisonValue, operator, actualValue, hasPassed, spanId}) => {
          const spanLabelList = getSpanSignature(spanId, trace)
            .map(({value}) => value)
            .concat([`#${spanId.slice(-4)}`]);

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
  }, [selectorValueList, spanListAssertionResult, trace]);

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
    <S.Container>
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
          onClick: () => onSpanSelected((record as TParsedAssertion).spanId),
        })}
      >
        <Table.Column
          title="Span Labels"
          dataIndex="spanLabels"
          key="spanLabels"
          ellipsis
          width="40%"
          render={(value: string[], record: TParsedAssertion) =>
            value
              // eslint-disable-next-line react/no-array-index-key
              .map((label, index) => <S.LabelBadge count={label} key={`${label}-${index}`} />)
              .concat(getIsSelected(record.spanId) ? [<S.SelectedLabelBadge count="selected" key="selected" />] : [])
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
    </S.Container>
  );
};

export default TraceAssertionsResultTable;
