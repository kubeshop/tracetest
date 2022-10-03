import {TSpanFlatAttribute} from '../types/Span.types';

const OutputService = () => ({
  getValueFromAttributeList(attributeList: TSpanFlatAttribute[], attribute: string, regex = '') {
    const value = attributeList.find(({key}) => attribute === key)?.value || '<No value>';

    return regex ? value.match(regex)?.[0] || '' : value;
  },
});

export default OutputService();
