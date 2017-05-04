package view

const TPL_INDEX_PAGE = `
<html>

<head>
	<title>
		index
	</title>
	
	<script src="http://apps.bdimg.com/libs/jquery/2.0.0/jquery.js"></script>
</head>

<body>
	<center>
		<br>
		<h2>添加或删除</h2>
		<form id="myForm">
			<p>型号: <input type="text" name="modelName"/></p>
			<p>参考价: <input type="text" name="standardPrice"/></p>
			<p>最低价: <input type="text" name="minPrice"/></p>
			<p>最高价: <input type="text" name="maxPrice"/></p>
			<input type="button" value="添加型号" id="btnAdd"/>
			<input type="button" value="删除型号" id="btnDel"/>
		</form>
		<br>
		<br>
		<h2>更新</h2>
		<select id="myUpdateSelect" onchange="displayPrice(this.options[this.selectedIndex].value)">
			<option selected>选择型号</OPTION> 
			{{range $_,$v := .Selects }}
				<option value ="{{$v}}">{{$v}}</option>
			{{end}}
		</select>
		<form id="updateForm" >
			<input type="hidden" type="text" id="modelName" name="modelName"/>
			<p>参考价: <input type="text" id="standardPrice" name="standardPrice"/></p>
			<p>最低价: <input type="text" id="minPrice" name="minPrice"/></p>
			<p>最高价: <input type="text" id="maxPrice" name="maxPrice"/></p>
			<input type="button" value="更新" id="btnUpdate"/>
		</form>
		<br>
		<br>
		<h2>查询</h2>
		<select id="mySelect" onchange="window.open(this.options[this.selectedIndex].value)">
			<option selected>选择型号</OPTION> 
			{{range $_,$v := .Selects }}
				<option value ="/jd?model={{$v}}">{{$v}}</option>
			{{end}}
		</select>
	</center>

	<script>
		var prices = new Map([
			{{range $index,$v := .Prices}}{{if $index}},{{end}}["{{$v.Name}}",{standardPrice:{{$v.StandardPrice}},minPrice:{{$v.MinPrice}},maxPrice:{{$v.MaxPrice}}}]{{end}}
		]);
		
		$('#btnAdd').click(function() {
			var AjaxURL= "/addModel";
			$.ajax({
				type: "POST",
				dataType: "html",
				url: AjaxURL,
				data: $('#myForm').serialize(),
				success: function (data) {
					alert(data);
				},
				error: function(data) {
					alert("error:"+data);
				}
			});
		});
			
		$('#btnDel').click(function() {
			var AjaxURL= "/delModel";
			$.ajax({
				type: "POST",
				dataType: "html",
				url: AjaxURL,
				data: $('#myForm').serialize(),
				success: function (data) {
						alert(data);
				},
				error: function(data) {
						alert("error:"+data);
				}
			});
		});
		
		$('#btnUpdate').click(function() {
			var AjaxURL= "/updatePrice";
			$.ajax({
				type: "POST",
				dataType: "html",
				url: AjaxURL,
				data: $('#updateForm').serialize(),
				success: function (data) {
					alert(data);
				},
				error: function(data) {
					alert("error:"+data);
				}
			});
		});
		
		function displayPrice(model){ 
			var modelName = document.getElementById("modelName"); 
        		var inputSP = document.getElementById("standardPrice"); 
        		var inputMinP = document.getElementById("minPrice"); 
			var inputMaxP = document.getElementById("maxPrice");
			var p = prices.get(model)
			//alert(p)
			modelName.value = model
			inputSP.value = p.standardPrice
			inputMinP.value = p.minPrice
			inputMaxP.value = p.maxPrice
       } 
	</script>
</body>

</html>`
