/* eslint-disable class-methods-use-this */
class AttributesService {
  requestAttributesListFN = (attributeList: any[]): any[] =>
    attributeList?.filter(a => {
      return (
        [
          'http.method',
          'http.url',
          'http.target',
          'http.host',
          'http.scheme',
          'http.request_content_length',
          'http.request_content_length_uncompressed',
          'http.retry_count exists',
          'http.user_agent',
        ].includes(a.key) || a.key.includes('http.request')
      );
    });

  responseAttributesFN = (attributeList: any[]) =>
    attributeList?.filter(a => {
      return a.key.includes('http.response');
    });
}

export default new AttributesService();
