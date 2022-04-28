export type TSpanAttributeValue = {
  stringValue: string;
  intValue: number;
  booleanValue: boolean;
  doubleValue: number;
  kvlistValue: {values: TSpanAttribute[]};
};

export type TSpanAttribute = {
  key: string;
  value: TSpanAttributeValue;
};
