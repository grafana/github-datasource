import type { Configuration } from 'webpack';
import { merge } from 'webpack-merge';
import CopyWebpackPlugin from 'copy-webpack-plugin';
import grafanaConfig, { Env } from './.config/webpack/webpack.config';

const config = async (env: Env): Promise<Configuration> => {
    const baseConfig = await grafanaConfig(env);
    return merge(baseConfig, {
        plugins: [new CopyWebpackPlugin({
            patterns: [
                {
                    from: '../skills/**/*',
                    to: './skills',
                    noErrorOnMissing: true
                },
                {
                    from: '../pkg/schema/dsconfig.json',
                    to: './schema/settings.schema.json',
                    noErrorOnMissing: true
                },
                {
                    from: '../pkg/schema/schema.gen.json',
                    to: './schema/v0alpha1.json',
                    noErrorOnMissing: true
                },
                {
                    from: '../pkg/schema/settings.gen.json',
                    to: './schema/v0alpha1/settings.json',
                    noErrorOnMissing: true
                },
                {
                    from: '../pkg/schema/settings.examples.gen.json',
                    to: './schema/v0alpha1/settings.examples.json',
                    noErrorOnMissing: true
                },
            ]
        })],
    });
};

export default config;
