const webpack = require('webpack');

module.exports = {
  context: __dirname + "/src",
  entry: './index.js',
  output: {
    path: __dirname + '/dist',
    filename: 'un-scrapy-client-bundle.js'
  },
  devServer: {
    contentBase: './dist',
    hot: true,
    inline: true
  },
  module: {
    loaders: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: 'babel-loader',
        query: {
          presets: ['es2015']
        }
      },
    ]
  },
  plugins: [
    new webpack.DefinePlugin({
      __DEV__: 'true',
      'process.env.NODE_ENV': `"${process.env.NODE_ENV || 'development'}"`
    }),
    new webpack.HotModuleReplacementPlugin()
  ]
}
