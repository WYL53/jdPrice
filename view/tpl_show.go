package view

const TPL_SHOW_PAGE = `<html>
<body>

<table id="tbl" border="1">
  	<tr>
		<th>    </th>
	    <th>店名</th>
	    <th>商品名</th>
	    <th>参考价</th>
		<th>当前价</th>
		<th>当前价-参考价</th>
		<th>评论数</th>
		<th>优惠信息</th>
  	</tr>
	{{range $i, $v :=  .Goods }}
		{{if eq $v.PriceDiff 0 }}
			<tr>
		{{end}}
		{{if gt $v.PriceDiff 0 }}
			<tr bgcolor="#FF0000">
		{{end}}
		{{if lt $v.PriceDiff 0 }}
			<tr bgcolor="#00FF00">
		{{end}}
			<td></td>
		    <td>{{$v.ShopName}}</td>
		    <td><a target="_blank" href="{{$v.GoodHref}}">{{$v.Name}}</a></td>
			<td>{{$.StandardPrice}}</td>
			<td><a target="_blank" href="/price?model={{$.Model}}&id={{$i}}">{{$v.Price}}</a></td>
			<td>{{$v.PriceDiff}}</td>
		    <td>{{$v.Sales}}</td>
			<td>{{$v.Etc}}</td>
		</tr>
	{{end}}
</table>

<script type="text/javascript">
	var obj = document.getElementById("tbl");
	var rowNum = obj.rows.length;
	for(var i=0;i<rowNum;i++){
		obj.rows[i+1].cells[0].innerHTML = "<td>" + (i+1).toString() +"</td>"
	}

</script>
</body>
</html>`
