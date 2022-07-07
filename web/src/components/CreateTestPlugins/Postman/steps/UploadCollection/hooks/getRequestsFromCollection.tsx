import {Collection, Item, ItemGroup, Request} from 'postman-collection';

export interface RequestDefinitionExtended extends Request {
  id: string;
  name: string;
}

function flattenCollectionItemsIntoRequestDefinition(structure: Item | ItemGroup<Item>): RequestDefinitionExtended[] {
  const itemRequest: RequestDefinitionExtended[] =
    'request' in structure ? [{...structure.request, name: structure.name, id: structure.id} as Request] : [];
  const groupRequest: RequestDefinitionExtended[] =
    'items' in structure
      ? structure.items
          .all()
          .map(a => flattenCollectionItemsIntoRequestDefinition(a))
          .flat()
      : [];
  return [...itemRequest, ...groupRequest];
}

export function getRequestsFromCollection(collection: Collection): RequestDefinitionExtended[] {
  try {
    return collection.items.all().flatMap(test => flattenCollectionItemsIntoRequestDefinition(test));
  } catch (e) {
    return [];
  }
}
