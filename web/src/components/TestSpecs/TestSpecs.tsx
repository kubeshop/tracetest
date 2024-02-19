import TestSpec from 'components/TestSpec';
import AssertionResults, {TAssertionResultEntry} from 'models/AssertionResults.model';
import {useCallback, useRef} from 'react';
import AutoSizer, {Size} from 'react-virtualized-auto-sizer';
import {VariableSizeList as List} from 'react-window';
import {useAppSelector} from 'redux/hooks';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import Empty from './Empty';
import * as S from './TestSpecs.styled';

interface IProps {
  assertionResults?: AssertionResults;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry, name: string): void;
  onOpen(selector: string): void;
  onRevert(originalSelector: string): void;
}

const TestSpecs = ({assertionResults, onDelete, onEdit, onOpen, onRevert}: IProps) => {
  const hiddenElementRef = useRef<HTMLDivElement>(null);
  const specs = useAppSelector(state => TestSpecsSelectors.selectSpecs(state));

  const getItemSize = useCallback(
    index => {
      const item = assertionResults?.resultList?.[index];
      const selector = item?.selector ?? '';
      const {name = ''} = specs.find(spec => spec.selector === selector) ?? {};
      const label = name || selector || 'All Spans';

      if (hiddenElementRef.current) {
        hiddenElementRef.current.innerText = label;
        return hiddenElementRef.current.offsetHeight;
      }

      return 0;
    },
    [assertionResults?.resultList, specs]
  );

  if (!assertionResults?.resultList?.length) {
    return <Empty />;
  }

  return (
    <>
      <S.ListContainer>
        <AutoSizer>
          {({height, width}: Size) => (
            <List
              height={height}
              itemCount={assertionResults.resultList.length}
              itemData={assertionResults.resultList}
              itemSize={getItemSize}
              width={width}
            >
              {({index, data, style}) => {
                const specResult = data[index];

                return specResult.resultList.length ? (
                  <div style={style}>
                    <TestSpec
                      key={specResult.id}
                      onDelete={onDelete}
                      onEdit={onEdit}
                      onOpen={onOpen}
                      onRevert={onRevert}
                      testSpec={specResult}
                    />
                  </div>
                ) : null;
              }}
            </List>
          )}
        </AutoSizer>
      </S.ListContainer>

      <S.HiddenElementContainer>
        <S.HiddenElement ref={hiddenElementRef} />
      </S.HiddenElementContainer>
    </>
  );
};

export default TestSpecs;
