import {downloadFile, escapeString, getObjectIncludesText} from '../Common';

describe('Common', () => {
  describe('miscelanious', () => {
    it('should call downloadFile', () => {
      const fileName = 'fileName';
      const link = downloadFile('fg', fileName);
      expect(link.getAttribute('download')).toBe(fileName);
    });
    it('should call getObjectIncludesText', () => {
      expect(getObjectIncludesText({hello: 33}, '33')).toBe(true);
    });
    it('should call escapeString', () => {
      expect(escapeString("'fg'")).toBe("\\'fg\\'");
    });
  });
});
