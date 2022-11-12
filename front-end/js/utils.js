utils={
	toHtml: function (val) {
	    var entityMap = {
	        "&": "&amp;",
	        "<": "&lt;",
	        ">": "&gt;",
	        '"': '&quot;',
	        "'": '&#39;',
	        "/": '&#x2F;'
	    };
	    return String(val).replace(/[&<>"'\/]/g, function (s) {
	        return entityMap[s];
	    });
	},
	toUrl: function (url) {
	    if(url.match(/^http/i)){
	        return encodeURI(url)
	    }
	    return '#'
	},
}