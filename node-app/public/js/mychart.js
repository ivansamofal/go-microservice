// chart.js
// const moment = require('moment');
// Обязательно дождитесь загрузки DOM, чтобы элемент с id "chart" уже был доступен
document.addEventListener("DOMContentLoaded", function() {
    // Здесь dataSet – предполагаемый массив с данными для графика.
    // Если dataSet передается динамически, вы можете определить его заранее или получить через AJAX.
    var dataSet = [
        [10, 20, 30, 40, 50, 60, 70, 80],
        [15, 25, 35, 45, 55, 65, 75, 85],
        [12, 22, 32, 42, 52, 62, 72, 82]
    ];

    var options = {
        series: [{
            name: 'PRODUCT A',
            data: dataSet[0]
        }, {
            name: 'PRODUCT B',
            data: dataSet[1]
        }, {
            name: 'PRODUCT C',
            data: dataSet[2]
        }],
        chart: {
            type: 'area',
            stacked: false,
            height: 350,
            zoom: {
                enabled: false
            }
        },
        dataLabels: {
            enabled: false
        },
        markers: {
            size: 0
        },
        fill: {
            type: 'gradient',
            gradient: {
                shadeIntensity: 1,
                inverseColors: false,
                opacityFrom: 0.45,
                opacityTo: 0.05,
                stops: [20, 100, 100, 100]
            }
        },
        yaxis: {
            labels: {
                style: {
                    colors: '#8e8da4'
                },
                offsetX: 0,
                formatter: function(val) {
                    return (val / 1000000).toFixed(2);
                }
            },
            axisBorder: {
                show: false
            },
            axisTicks: {
                show: false
            }
        },
        xaxis: {
            type: 'datetime',
            tickAmount: 8,
            min: new Date("01/01/2014").getTime(),
            max: new Date("01/20/2014").getTime(),
            labels: {
                rotate: -15,
                rotateAlways: true,
                formatter: function(val, timestamp) {
                    return moment(new Date(timestamp)).format("DD MMM YYYY");
                }
            }
        },
        title: {
            text: 'Irregular Data in Time Series',
            align: 'left',
            offsetX: 14
        },
        tooltip: {
            shared: true
        },
        legend: {
            position: 'top',
            horizontalAlign: 'right',
            offsetX: -10
        }
    };

    var chart = new ApexCharts(document.querySelector("#chart"), options);
    chart.render();
    console.log(chart);
    console.log('loaded');
});
