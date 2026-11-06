import type { Configuration } from 'webpack';
import { merge } from 'webpack-merge';
import CopyWebpackPlugin from 'copy-webpack-plugin';
import grafanaConfig, { Env } from './.config/webpack/webpack.config';

const config = async (env: Env): Promise<Configuration> => {
    const baseConfig = await grafanaConfig(env);
    return merge(baseConfig, {
        plugins: [new CopyWebpackPlugin({
            patterns: [
                { from: '../pkg/schema/dsconfig.json', to: './schema/dsconfig.json' },
                { from: '../pkg/schema/schema.gen.json', to: './schema/v0alpha1.json' },
                { from: '../pkg/schema/settings.gen.json', to: './schema/v0alpha1/settings.json' },
                { from: '../pkg/schema/settings.examples.gen.json', to: './schema/v0alpha1/settings.examples.json' },
            ],
        })],
    });
};

export default config;
