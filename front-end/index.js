function openIFrame(){
    let head = document.getElementsByTagName('head')[0];
	var script1 = document.createElement('script');
    script1.setAttribute('type','text/javascript');
    script1.setAttribute('src','https://cdn.bootcdn.net/ajax/libs/jquery/3.6.0/jquery.min.js');
    head.appendChild(script1);
    var script2 = document.createElement('script');
    script2.setAttribute('type','text/javascript');
    script2.setAttribute('src','https://cdn.bootcdn.net/ajax/libs/fancybox/3.5.7/jquery.fancybox.js');
    head.appendChild(script2);
    var link = document.createElement('link');
    link.setAttribute('rel','stylesheet');
    link.setAttribute('href','https://cdn.bootcdn.net/ajax/libs/fancybox/3.5.7/jquery.fancybox.css');
    head.appendChild(link);
	setTimeout(function(){  
       $.fancybox.open({
           src: 'test2_multipart_success.html',  //弹框中需要展示的广告页面
           type: 'iframe',
           width: 2000,
           height: 2000,
           autoSize: true,
        });
      
        //可以通过下面的方式设置iframe的一些css样式          
        let fancyboxFrame = document.getElementsByClassName('fancybox-iframe')[0];
        fancyboxFrame.style.paddingTop = '40px'
        fancyboxFrame.style.paddingLeft = '20px'
    },100);
           
}