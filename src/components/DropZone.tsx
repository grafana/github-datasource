import React, { useMemo } from 'react';
import { useDropzone, DropzoneOptions, DropzoneRootProps } from 'react-dropzone';

const defaultBaseStyle = {
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  height: '150px',
  backgroundColor: '#212327',
  padding: 24,
  borderWidth: 2,
  borderRadius: 2,
  borderColor: '#808080',
  borderStyle: 'dashed',
  outline: 'none',
  transition: 'border .24s ease-in-out',
};

const defaultActiveStyle = {
  borderColor: '#2196f3',
};

const defaultAcceptStyle = {
  borderColor: '#EF843C',
};

const defaultRejectStyle = {
  borderColor: '#E24D42',
};

type Style = { [key: string]: any };

export interface Props extends DropzoneOptions {
  children: React.ReactNode;
  baseStyle?: Style;
  activeStyle?: Style;
  acceptStyle?: Style;
  rejectStyle?: Style;
}

export function DropZone({ children, baseStyle = {}, activeStyle = {}, acceptStyle = {}, rejectStyle = {}, ...rest }: Props) {
  const { getRootProps, getInputProps, isDragActive, isDragAccept, isDragReject } = useDropzone(rest);

  const style = useMemo(
    () => ({
      ...defaultBaseStyle,
      ...baseStyle,
      ...(isDragActive ? { ...defaultActiveStyle, ...activeStyle } : {}),
      ...(isDragAccept ? { ...defaultAcceptStyle, ...acceptStyle } : {}),
      ...(isDragReject ? { ...defaultRejectStyle, ...rejectStyle } : {}),
    }),
    [isDragActive, isDragReject]
  );

  return (
    <div {...getRootProps({ style } as DropzoneRootProps)}>
      <input {...getInputProps()} />
      {children}
    </div>
  );
}
