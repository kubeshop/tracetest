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

    it('should handle valid value', () => {
      const value = 'https://tracetest.io/';
      expect(Validator.required(value)).toBeTruthy();
    });
  });
});
