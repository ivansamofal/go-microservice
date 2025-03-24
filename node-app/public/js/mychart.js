const moment = require('moment');
import ApexCharts from 'apexcharts'

document.addEventListener("DOMContentLoaded", function () {
    const currentUrl = new URL(window.location.href);
    currentUrl.port = '8080';
    const apiUrl = currentUrl.origin + '/api/calculations';
    fetch(apiUrl)
        .then(response => response.json())
        .then(data => {
            // Создаем график EMA
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
                        x: {
                            title: {
                                display: true,
                                text: 'Период'
                            }
                        },
                        y: {
                            title: {
                                display: true,
                                text: 'Значение EMA'
                            }
                        }
                    }
                }
            });

            // Создаем график MACD
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
                        x: {
                            title: {
                                display: true,
                                text: 'Период'
                            }
                        },
                        y: {
                            title: {
                                display: true,
                                text: 'Значение MACD'
                            }
                        }
                    }
                }
            });

            const histogram = document.getElementById('histogram').getContext('2d');
            new Chart(histogram, {
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
                        x: {
                            title: {
                                display: true,
                                text: 'Период'
                            }
                        },
                        y: {
                            title: {
                                display: true,
                                text: 'Значение Histogram'
                            }
                        }
                    }
                }
            });

            const rsi = document.getElementById('rsi').getContext('2d');
            new Chart(rsi, {
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
                        x: {
                            title: {
                                display: true,
                                text: 'Период'
                            }
                        },
                        y: {
                            title: {
                                display: true,
                                text: 'Значение RSI'
                            }
                        }
                    }
                }
            });

            const signalLine = document.getElementById('signalLine').getContext('2d');
            new Chart(signalLine, {
                type: 'line',
                data: {
                    labels: data.rsi.map((_, i) => i + 1),
                    datasets: [{
                        label: 'Signal Line',
                        data: data.signalLine,
                        borderColor: 'rgba(75, 192, 192, 1)',
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        fill: false
                    }]
                },
                options: {
                    responsive: true,
                    scales: {
                        x: {
                            title: {
                                display: true,
                                text: 'Период'
                            }
                        },
                        y: {
                            title: {
                                display: true,
                                text: 'Значение Signal Line'
                            }
                        }
                    }
                }
            });

            // Предполагается, что data – объект с данными для графиков: data.rsi, data.histogram, data.macdLine, data.ema

// Общие настройки, которые повторяются для всех графиков
            var baseOptions = {
                chart: {
                    id: 'chartyear',
                    type: 'area',
                    height: 270,
                    background: '#F6F8FA',
                    toolbar: {
                        show: false,
                        autoSelected: 'pan'
                    },
                    events: {}
                },
                stroke: {
                    width: 0,
                    curve: 'monotoneCubic'
                },
                dataLabels: {
                    enabled: false
                },
                fill: {
                    opacity: 1,
                    type: 'solid'
                },
                yaxis: {
                    show: true,
                    forceNiceScale: true
                }
            };

// Функция для создания настроек графика на основе базовых настроек
            function createChartOptions(series, colors) {
                // Создаем копию базовых опций, чтобы не изменять исходный объект
                var options = Object.assign({}, baseOptions, {
                    series: series,
                    colors: colors
                });
                return options;
            }

// Настройки для каждого графика с разными данными и цветами
            var rsiOptions = createChartOptions(
                [{
                    name: 'RSI',
                    data: data.rsi
                }],
                ['#FF7F00'] // например, оранжевый для RSI
            );

            var histogramOptions = createChartOptions(
                [{
                    name: 'Histogram',
                    data: data.histogram
                }],
                ['#00e396'] // зеленый для Histogram
            );

            var macdOptions = createChartOptions(
                [{
                    name: 'MACD',
                    data: data.macdLine
                }],
                ['#FEB019'] // желтый/оранжевый для MACD
            );

            var emaOptions = createChartOptions(
                [{
                    name: 'EMA',
                    data: data.ema
                }],
                ['#40bf40'] // зеленый для EMA
            );

            new ApexCharts(document.querySelector("#rsiNewChart"), rsiOptions).render();
            new ApexCharts(document.querySelector("#histogramNewChart"), histogramOptions).render();
            new ApexCharts(document.querySelector("#macdNewChart"), macdOptions).render();
            new ApexCharts(document.querySelector("#emaNewChart"), emaOptions).render();

        })
        .catch(error => console.error('Ошибка загрузки данных:', error));;


/* ================================ FIN JS  INDICATORS ====================================== */
        });