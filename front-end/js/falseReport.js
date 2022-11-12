;(function(){
	var local = {
		// 初始化
		init: function(){
			// window.add
			window.addEventListener('error',(message, source, lineno, colno, error)=>{
				this.report(message, source, lineno, colno, error);
			});

			window.LN = window.LN || {};
			window.LN.report = window.LN.report || {};
			window.LN.report.notice = (message, source, lineno, colno, error)=>{
				this.report(message, source, lineno, colno, error);
			}
		},
		report: function(message, source, lineno, colno, error){
			// message：错误信息（字符串）。可用于HTML onerror=""处理程序中的event。
		    // source：发生错误的脚本URL（字符串）
		    // lineno：发生错误的行号（数字）
		    // colno：发生错误的列号（数字）
		    // error：Error对象（对象）;
		    var _Img = document.createElement('img');
		    var _url = 'http://123.56.16.33/falseReport?'+
		    	'message='+encodeURI(message)+'&'
				'source='+encodeURI(source)+'&'
				'lineno='+encodeURI(lineno)+'&'
				'colno='+encodeURI(colno)+'&'
				'error='+encodeURI(error)+'';
		}
	};
	local.init();
})();