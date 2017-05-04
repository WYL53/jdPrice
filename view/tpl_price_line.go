package view

const TPL_PRICE_LINE_PAGE = `
<html>

<head>
	<script src="/static/Chart.js"></script>
</head>

<body>
	<canvas id="myChart" width="1000" height="400"></canvas>
<script>
var ctx = document.getElementById("myChart");
var myChart = new Chart(ctx, {
    type: 'line',
    data: {
        labels : [ {{range $index,$v := .Times}}{{if $index}},{{end}}"{{$v}}"{{end}} ],
		datasets : [{
				label: "价格走势",
				fillColor : "rgba(220,220,220,0.5)",
				strokeColor : "rgba(220,220,220,1)",
				pointColor : "rgba(220,220,220,1)",
				pointStrokeColor : "#fff",
				pointHoverRadius: 5,
				pointRadius: 1,
				data : [{{range $index,$v := .Prices}}{{if $index}},{{end}}{{$v}}{{end}}]
			}]
    		},
   	options: { responsive: false }
});
</script>
</body>
</html>
`
