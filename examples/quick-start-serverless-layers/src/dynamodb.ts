import { DynamoDBClient } from '@aws-sdk/client-dynamodb';
import { PutCommand, DynamoDBDocumentClient, GetCommand } from '@aws-sdk/lib-dynamodb';

const client = new DynamoDBClient({});
const docClient = DynamoDBDocumentClient.from(client);

const { TABLE_NAME = '' } = process.env;

const DynamoDb = {
  async put<T extends Record<string, any>>(item: T): Promise<T> {
    const command = new PutCommand({
      TableName: TABLE_NAME,
      Item: item,
    });

    const result = await docClient.send(command);

    return result.Attributes as T;
  },

  async get<T>(id: number): Promise<T | undefined> {
    const command = new GetCommand({
      TableName: TABLE_NAME,
      Key: {
        id,
      },
    });

    const result = await docClient.send(command);

    return result.Item as T | undefined;
  },
};

export default DynamoDb;
