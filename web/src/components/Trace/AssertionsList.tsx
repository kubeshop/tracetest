import React, {useState} from 'react';
import {Button} from 'antd';

import jemsPath from 'jmespath';

import {filterBySpanId} from 'utils';
import {IAttribute, ISpan} from 'types';
import {useGetTestAssertionsQuery} from 'services/TestService';
import AssertionsResultTable from 'components/AssertionsTable/AssertionsTable';

import CreateAssertionModal from './CreateAssertionModal';

interface IProps {
  testId: string;
  targetSpan: ISpan;
  trace: any;
}

const AssertionList = ({testId, targetSpan, trace}: IProps) => {
  const [openCreateAssertion, setOpenCreateAssertion] = useState(false);
  const {data: testAssertions} = useGetTestAssertionsQuery(testId);
  const attrs: IAttribute[] = jemsPath.search(trace, filterBySpanId(targetSpan.spanId));

  const attributesTree = attrs?.reduce((acc: any, item: any) => {
    const resource = acc[item.type] || {};
    resource.title = item.type;
    resource.key = item.type;
    resource.children = resource.children || [];
    resource.children.push({title: `${item.key} = ${item.value}`, key: item.key});
    acc[item.type] = resource;
    return acc;
  }, {});

  return (
    <div>
      <Button style={{marginBottom: 8}} onClick={() => setOpenCreateAssertion(true)}>
        New Assertion
      </Button>
      <AssertionsResultTable />
      <CreateAssertionModal
        key={`KEY_${targetSpan.spanId}`}
        testId={testId}
        trace={trace}
        span={targetSpan}
        open={openCreateAssertion}
        onClose={() => setOpenCreateAssertion(false)}
      />
    </div>
  );
};

export default AssertionList;
