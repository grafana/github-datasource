export const prepareAnnotation = (json: any) => {
  console.log('json', json);
  json.target = json.target ?? {};
  return json;
};
