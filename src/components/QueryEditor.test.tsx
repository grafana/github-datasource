import { getGoogleSheetRangeInfoFromURL } from './QueryEditor';

describe('QueryEditor', () => {
  it('should extract id from URL', () => {
    const url = 'https://docs.google.com/spreadsheets/d/1m2idieRUHdzWTu3_cpYs1lUfP_jwfgL8NBaLtqLmia8/edit#gid=790763898&range=B19:F20';
    const info = getGoogleSheetRangeInfoFromURL(url);
    expect(info.spreadsheet).toBe('1m2idieRUHdzWTu3_cpYs1lUfP_jwfgL8NBaLtqLmia8');
    expect(info.range).toBe('B19:F20');
  });
});
