import {SemanticGroupNameNodeMap} from '../constants/SemanticGroupNames.constants';
import {ISpan} from '../types/Span.types';

const SpanService = () => ({
  getSpanNodeInfo(span: ISpan) {
    const signatureObject = span.signature.reduce<Record<string, string>>(
      (signature, {propertyName, value}) => ({
        ...signature,
        [propertyName]: value,
      }),
      {}
    );

    const {primary, type} = SemanticGroupNameNodeMap[span.type];

    const attributeKey = primary.find(key => signatureObject[key]) || '';

    return {
      primary: signatureObject[attributeKey] || '',
      heading: signatureObject[type],
    };
  },
});

export default SpanService();
