const moment = require('moment');
import ApexCharts from 'apexcharts'

document.addEventListener("DOMContentLoaded", function () {
    function loadAndRenderCharts(apiQuery, suffix) {
        const currentUrl = new URL(window.location.href);
        currentUrl.port = '8080';
        const apiUrl = currentUrl.origin + '/api/calculations' + apiQuery;

        fetch(apiUrl)
            .then(response => response.json())
            .then(data => {
                if (!suffix) {
                    const emaCtx = document.getElementById('emaChart').getContext('2d');
                    new Chart(emaCtx, {
                        type: 'line',
                        data: {
                            labels: data.ema.map((_, i) => i + 1),
                            datasets: [{
                                label: 'EMA',
                                data: data.ema,
                                borderColor: 'rgba(75, 192, 192, 1)',
                                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                                fill: false
                            }]
                        },
                        options: {
                            responsive: true,
                            scales: {
                                x: { title: { display: true, text: 'Период' } },
                                y: { title: { display: true, text: 'Значение EMA' } }
                            }
                        }
                    });

                    const macdCtx = document.getElementById('macdChart').getContext('2d');
                    new Chart(macdCtx, {
                        type: 'line',
                        data: {
                            labels: data.macdLine.map((_, i) => i + 1),
                            datasets: [{
                                label: 'MACD',
                                data: data.macdLine,
                                borderColor: 'rgba(255, 99, 132, 1)',
                                backgroundColor: 'rgba(255, 99, 132, 0.2)',
                                fill: false
                            }]
                        },
                        options: {
                            responsive: true,
                            scales: {
                                x: { title: { display: true, text: 'Период' } },
                                y: { title: { display: true, text: 'Значение MACD' } }
                            }
                        }
                    });

                    const histogramCtx = document.getElementById('histogram').getContext('2d');
                    new Chart(histogramCtx, {
                        type: 'line',
                        data: {
                            labels: data.histogram.map((_, i) => i + 1),
                            datasets: [{
                                label: 'Histogram',
                                data: data.histogram,
                                borderColor: 'rgba(75, 192, 192, 1)',
                                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                                fill: false
                            }]
                        },
                        options: {
                            responsive: true,
                            scales: {
                                x: { title: { display: true, text: 'Период' } },
                                y: { title: { display: true, text: 'Значение Histogram' } }
                            }
                        }
                    });

                    const rsiCtx = document.getElementById('rsi').getContext('2d');
                    new Chart(rsiCtx, {
                        type: 'line',
                        data: {
                            labels: data.rsi.map((_, i) => i + 1),
                            datasets: [{
                                label: 'RSI',
                                data: data.rsi,
                                borderColor: 'rgba(75, 192, 192, 1)',
                                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                                fill: false
                            }]
                        },
                        options: {
                            responsive: true,
                            scales: {
                                x: { title: { display: true, text: 'Период' } },
                                y: { title: { display: true, text: 'Значение RSI' } }
                            }
                        }
                    });
                }

                // Создаем графики с ApexCharts (пример)
                var baseOptions = {
                    chart: {
                        id: 'chartyear',
                        type: 'area',
                        height: 270,
                        background: '#F6F8FA',
                        toolbar: { show: false, autoSelected: 'pan' },
                        events: {}
                    },
                    stroke: { width: 0, curve: 'monotoneCubic' },
                    dataLabels: { enabled: false },
                    fill: { opacity: 1, type: 'solid' },
                    yaxis: { show: true, forceNiceScale: true }
                };

                function createChartOptions(series, colors) {
                    var options = Object.assign({}, baseOptions, { series: series, colors: colors });
                    return options;
                }

                var rsiOptions = createChartOptions(
                    [{ name: 'RSI', data: data.rsi }],
                    ['#FF7F00']
                );
                var histogramOptions = createChartOptions(
                    [{ name: 'Histogram', data: data.histogram }],
                    ['#00e396']
                );
                var macdOptions = createChartOptions(
                    [{ name: 'MACD', data: data.macdLine }],
                    ['#FEB019']
                );
                var emaOptions = createChartOptions(
                    [{ name: 'EMA', data: data.ema }],
                    ['#40bf40']
                );

                new ApexCharts(document.querySelector("#rsiNewChart" + suffix), rsiOptions).render();
                new ApexCharts(document.querySelector("#histogramNewChart" + suffix), histogramOptions).render();
                new ApexCharts(document.querySelector("#macdNewChart" + suffix), macdOptions).render();
                new ApexCharts(document.querySelector("#emaNewChart" + suffix), emaOptions).render();
            })
            .catch(error => console.error('Ошибка загрузки данных:', error));
    }

    loadAndRenderCharts('?x=2', '');
    loadAndRenderCharts('?x=4', '2');
    loadAndRenderCharts('?x=8', '3');
});
