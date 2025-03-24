const express = require('express');
const app = express();
const port = process.env.PORT || 3000;



// Раздача статических файлов из папки public
app.use(express.static('public'));

// Пример API-эндпоинта, который возвращает данные для графиков.
// В реальном приложении вы можете получать данные с вашего Go микросервиса.
app.get('/api/calculations', (req, res) => {
    // Пример данных — массивы значений индикаторов.
    const data = {
        ema: req.ema,
        macdLine: req.macdLine,
        histogram: req.histogram,
        rsi: req.rsi,
        signalLine: req.signalLine,
        // Добавьте другие массивы, если нужно
    };
    res.json(data);
});

app.listen(port, () => {
    console.log(`Server is running on port ${port}`);
});
