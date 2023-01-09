import { AnnotationEventFieldSource } from '@grafana/data';

const isNewAnnotation = (json: any) => {
  return !('annotation' in json);
};

export const prepareAnnotation = (json: any) => {
  if (isNewAnnotation(json)) {
    return json;
  }

  const { annotation, ...annotationRest } = json;
  const { field, timeField, ...targetRest } = annotation;

  const migratedAnnotation = {
    ...annotationRest,
    target: targetRest,
    mappings: {
      text: field
        ? {
            source: AnnotationEventFieldSource.Field,
            value: field,
          }
        : undefined,
      time: timeField
        ? {
            source: AnnotationEventFieldSource.Field,
            value: timeField,
          }
        : undefined,
      ...annotationRest.mappings,
    },
  };

  return migratedAnnotation;
};
