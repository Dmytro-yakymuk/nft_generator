const path = require("path");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const StylelintPlugin = require("stylelint-webpack-plugin");
const zlib = require("zlib");
const CompressionPlugin = require("compression-webpack-plugin");

const PRODUCTION_PLUGINS = [
    new HtmlWebpackPlugin({
        title: "Generator",
        template: "./public/index.html",
    }),
    new CleanWebpackPlugin(),
    // the content hash in the names lets the browser keep the files and see a new build at once.
    new MiniCssExtractPlugin({
        filename: "[name].[contenthash].css",
        chunkFilename: "[id].[contenthash].css",
    }),
    new StylelintPlugin({ fix: true }),
    // compressed copies next to the files; the server sends one when the browser accepts its encoding.
    new CompressionPlugin({
        filename: "[path][base].br",
        algorithm: "brotliCompress",
        test: /\.(js|css|html|svg|ttf|otf)$/,
        compressionOptions: {
            params: {
                [zlib.constants.BROTLI_PARAM_QUALITY]: 11,
            },
        },
        threshold: 1024,
        minRatio: 0.8,
        deleteOriginalAssets: false,
    }),
    new CompressionPlugin({
        filename: "[path][base].gz",
        algorithm: "gzip",
        test: /\.(js|css|html|svg|ttf|otf)$/,
        compressionOptions: { level: 9 },
        threshold: 1024,
        minRatio: 0.8,
        deleteOriginalAssets: false,
    }),
];

const DEVELOPMENT_PLUGINS = [
    new HtmlWebpackPlugin({
        title: "Generator",
        template: "./public/index.html",
    }),
];

const OPTIMIZATION = {
    splitChunks: {
        chunks: "all",
    },
    minimize: true,
    // "..." keeps the default JavaScript minimizer (Terser) next to the CSS one.
    minimizer: [
        "...",
        new CssMinimizerPlugin(),
    ],
};

const IMAGE_COMPRESSION = [
    {
        loader: "image-webpack-loader",
        options: {
            mozjpeg: {
                progressive: false,
            },
            optipng: {
                enabled: true,
            },
            pngquant: {
                quality: [0.8, 0.8],
                speed: 2,
            },
        },
    },
];

module.exports = (env, argv) => {
    const isProduction = argv.mode === 'production';

    return {
        experiments: {
            asset: true,
        },
        entry: "./src/index.tsx",
        target: "web",
        watch: !isProduction,
        watchOptions: {
            aggregateTimeout: 200,
            poll: 1000,
            ignored: /node_modules/,
        },
        output: {
            path: path.resolve(__dirname, "dist/"),
            filename: "[name].[hash].js",
            publicPath: "/static/dist/",
        },
        plugins: isProduction ? PRODUCTION_PLUGINS : DEVELOPMENT_PLUGINS,
        devServer: {
            port: 3000,
            open: true,
            historyApiFallback: true,
        },
        resolve: {
            alias: {
                "@": path.resolve(__dirname, "./src/"),
                "@components": path.resolve(__dirname, "./src/app/components/"),
                "@internal": path.resolve(__dirname, "./src/app/internal/"),
                "@static": path.resolve(__dirname, "./src/app/static/"),
                "@store": path.resolve(__dirname, "./src/app/store"),
                "@views": path.resolve(__dirname, "./src/app/views/")
            },
            extensions: [".ts", ".tsx", ".js", ".jsx"],
            modules: ["node_modules"],
        },
        module: {
            rules: [
                {
                    test: /\.m?(tsx|ts)$/i,
                    exclude: /(node_modules)/,
                    use: [
                        {
                            loader: "ts-loader",
                        },
                    ],
                },
                {
                    test: /\.(scss)$/,
                    exclude: /(node_modules)/,
                    use: isProduction
                        ?
                        [
                            MiniCssExtractPlugin.loader,
                            "css-loader",
                            "sass-loader",
                        ]
                        :
                        [
                            "style-loader",
                            "css-loader",
                            "sass-loader",
                        ]
                },
                {
                    test: /\.(css)$/,
                    use: [
                        "style-loader",
                        "css-loader",
                    ],
                },
                {
                    test: /\.(woff|woff2|eot|ttf|otf)$/i,
                    exclude: /(node_modules)/,
                    type: "asset/resource",
                    generator: {
                        filename: "fonts/[name][hash:5][ext]",
                    },
                },
                {
                    test: /\.(jpe|jpg|png|svg|gif|webp)(\?.*$|$)/,
                    exclude: /(node_modules)/,
                    type: "asset/resource",
                    generator: {
                        filename: "images/[name][hash:5][ext]",
                    },
                    use: isProduction ? IMAGE_COMPRESSION : [],
                },
            ],
        },
        optimization: isProduction ? OPTIMIZATION : {},
    };
}
