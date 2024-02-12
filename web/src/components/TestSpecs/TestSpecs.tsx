import TestSpec from 'components/TestSpec';
import AutoSizer, {Size} from 'react-virtualized-auto-sizer';
import {FixedSizeList as List} from 'react-window';
import AssertionResults, {TAssertionResultEntry} from 'models/AssertionResults.model';
import Empty from './Empty';

interface IProps {
  assertionResults?: AssertionResults;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry, name: string): void;
  onOpen(selector: string): void;
  onRevert(originalSelector: string): void;
}

const TestSpecs = ({assertionResults, onDelete, onEdit, onOpen, onRevert}: IProps) => {
  if (!assertionResults?.resultList?.length) {
    return <Empty />;
  }

  return (
    <AutoSizer>
      {({height, width}: Size) => (
        <List
          height={height}
          itemCount={assertionResults.resultList.length}
          itemData={assertionResults.resultList}
          itemSize={10}
          width={width}
        >
          {({index, data}) => {
            const specResult = data[index];

            return specResult.resultList.length ? (
              <TestSpec
                key={specResult.id}
                onDelete={onDelete}
                onEdit={onEdit}
                onOpen={onOpen}
                onRevert={onRevert}
                testSpec={specResult}
              />
            ) : null;
          }}
        </List>
      )}
    </AutoSizer>
  );
};

export default TestSpecs;
