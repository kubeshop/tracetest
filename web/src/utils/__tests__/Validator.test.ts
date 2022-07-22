import Validator from '../Validator';

describe('Validator', () => {
  describe('required', () => {
    it('should handle empty value', () => {
      const value = '';
      expect(Validator.required(value)).toBeFalsy();
    });

    it('should handle valid value', () => {
      const value = 'value';
      expect(Validator.required(value)).toBeTruthy();
    });
  });

  describe('url', () => {
    it('should handle non string value', () => {
      const value = 123;
      expect(Validator.url(value)).toBeFalsy();
    });

    it('should be invalid', () => {
      expect(Validator.url('http://demo.@')).toBe(false);
      expect(Validator.url('htt://demo')).toBe(false);
    });

    it('should be valid', () => {
      expect(Validator.required('https://tracetest.io/')).toBeTruthy();
      expect(Validator.url('http://demo')).toBe(true);
      expect(Validator.url('http://demo.com')).toBe(true);
      expect(Validator.url('https://www.demo.com')).toBe(true);
    });
  });
});
