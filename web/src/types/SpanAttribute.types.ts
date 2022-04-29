export interface ISpanAttributeValue {
  stringValue: string;
  intValue: number;
  booleanValue: boolean;
  doubleValue: number;
  kvlistValue: {values: ISpanAttribute[]};
}

export interface ISpanAttribute {
  key: string;
  value: ISpanAttributeValue;
}
