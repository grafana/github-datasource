export const prepareAnnotation = (json: any) => {
  // migration: move annotation object to target for old annotation
  json.target = json.annotation ?? json.target;
  return json;
};
