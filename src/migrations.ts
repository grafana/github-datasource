export const prepareAnnotation = (json: any) => {
  json.target = json.annotation;
  return json;
};
