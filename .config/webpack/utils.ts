import fs from 'fs';
import path from 'path';
import util from 'util';
import glob from 'glob';
import { SOURCE_DIR } from './constants';

const globAsync = util.promisify(glob);

export function getPackageJson() {
  return require(path.resolve(process.cwd(), 'package.json'));
}

export function getPluginId() {
  const { id } = require(path.resolve(process.cwd(), `${SOURCE_DIR}/plugin.json`));

  return id;
}

export function hasReadme() {
  return fs.existsSync(path.resolve(process.cwd(), SOURCE_DIR, 'README.md'));
}

// Support bundling nested plugins by finding all plugin.json files in src directory
// then checking for a sibling module.[jt]sx? file.
export async function getEntries(): Promise<Record<string, string>> {
  const pluginsJson = await globAsync('**/src/**/plugin.json', { absolute: true });

  const plugins = await Promise.all(pluginsJson.map((pluginJson) => {
      const folder = path.dirname(pluginJson);
      return globAsync(`${folder}/module.{ts,tsx,js,jsx}`, { absolute: true });
    })
  );

  return plugins.reduce((result, modules) => {
    return modules.reduce((result, module) => {
      const pluginPath = path.dirname(module);
      const pluginName = path.relative(process.cwd(), pluginPath).replace(/src\/?/i, '');
      const entryName = pluginName === '' ? 'module' : `${pluginName}/module`;

      result[entryName] = module;
      return result;
    }, result);
  }, {});
}
