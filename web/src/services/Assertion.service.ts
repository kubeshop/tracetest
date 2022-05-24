import {TRawAssertionResult} from '../types/Assertion.types';

const AssertionService = () => ({
  getSpanCount(resultList: TRawAssertionResult[]): number {
    const spanIdList = resultList.reduce<string[]>((list, {spanResults}) => {
      const tmpList: string[] = [];

      spanResults?.forEach(({spanId = ''}) => {
        if (!tmpList.includes(spanId)) tmpList.push(spanId);
      });

      return list.concat(tmpList);
    }, []);

    return spanIdList.length;
  },
});

export default AssertionService();
