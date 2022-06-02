import Date from '../Date';

describe('Date', () => {
  describe('format', () => {
    it('should handle empty value', () => {
      const date = '';
      expect(Date.format(date)).toBe('');
    });

    it('should handle invalid value', () => {
      const date = 'invalid-date';
      expect(Date.format(date)).toBe('');
    });

    it('should handle valid value', () => {
      const date = '2022-06-02T17:28:48';
      expect(Date.format(date, 'yyyy/MM/dd')).toBe('2022/06/02');
    });
  });
});
