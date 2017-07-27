var url = 'http://search.jd.com/Search?enc=utf-8&keyword='
var page = require('webpage').create();
var system = require('system');

if (system.args.length === 1) {
  console.log('Usage: jd_spider.js <some modelName>');
  phantom.exit();
}

url = url + system.args[1];

page.injectJs("jquery-1.6.1.min.js");

page.onConsoleMessage = function(msg) {
	console.log(msg);
};

function getData(status){
	if ( status === "success" ) {
		page.evaluate(function(text) {
			var result = '';
			var shopJD="京东自营";
			var spilt = ",";
	   		$(".gl-i-wrap").each(function(){
				var shop = $(this).find(".p-shop").find("a");
				var shopName = shop.text();
				var shopHerf = shop.attr("href");
				if(shopName.length === 0){
					shopName = shopJD;
				}
				var pName = $(this).find(".p-name").find("em").text();
				var pNameHref =  $(this).find(".p-name").children("a").attr("href");
				var price = $(this).find("strong").text();
				var iconString = "";
				var icons = $(this).find(".p-icons").find("i").each(function(){iconString+=$(this).text()+spilt});
				result += shopName + ":" + pName + ":"+price + ":"+iconString+":" +shopHerf+":"+pNameHref+"\n";
			});

			console.log(result);
	  	});
	}
  	phantom.exit();
}

page.open(url, getData);


/*
page.open(url,function(status) {

    //当打开成功后。输入检索内容并点击搜索
    //注意这两个按钮的的ID还是需要人为去看一下的
    if ( status === "success" ) {
        page.evaluate(function(text) {
            $("#kw1").val(text);
            $("#su1").click();
        }, "hello world");

    }

});
*/