import {useMemo} from 'react';
import SpanService from '../../services/Span.service';
import {ISpan, ISpanFlatAttribute} from '../../types/Span.types';

const renderTitle = (title: string, index: number) => <span key={`KEY_${title}_${index}`}>{title}</span>;

const renderItem = ({key}: ISpanFlatAttribute) => ({
  value: key,
  label: (
    <div
      key={key}
      style={{
        display: 'flex',
        justifyContent: 'space-between',
      }}
    >
      {key}
    </div>
  ),
});

const useAttributeList = (selectedSpan: ISpan, affectedSpanList: ISpan[]) => {
  return useMemo(() => {
    const {intersectedList, differenceList} = SpanService.getSelectedSpanListAttributes(selectedSpan, affectedSpanList);
    return [
      {
        label: renderTitle('Across all Spans', 0),
        options: intersectedList.map(el => renderItem(el)),
      },
      {
        label: renderTitle('For selected span', 1),
        options: differenceList.map(el => renderItem(el)),
      },
    ];
  }, [affectedSpanList, selectedSpan]);
};

export default useAttributeList;
