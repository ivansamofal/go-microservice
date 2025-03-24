const path = require('path');

module.exports = {
    entry: './public/js/mychart.js', // ваша основная точка входа
    output: {
        filename: 'bundle.js', // собранный файл
        path: path.resolve(__dirname, 'public/js/dist'),
    },
    mode: 'development', // или 'production'
    module: {
        rules: [
            {
                test: /\.js$/,            // обрабатываем все файлы .js
                exclude: /node_modules/,  // исключаем node_modules
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env']
                    }
                }
            }
        ]
    },
    mode: 'development' // или 'production'
};
