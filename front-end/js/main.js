let myImage = document.querySelector('img');
let myButton = document.querySelector('button');
let myHeading = document.querySelector('h1');
const para = document.querySelector('p');
const sourcefiles = document.querySelector('.sourcefiles')


import BMF from 'browser-md5-file'; //这个使用npm
const el = document.getElementById('upload');
const bmf = new BMF();
 
el.addEventListener('change', handle, false);//原生调用handle方法，监听选择的文件
 
function handle(e) {
  const file = e.target.files[0];
  bmf.md5(
    file,
    (err, md5) => {
      console.log('err:', err);
      console.log('md5 string:', md5); // 97027eb624f85892c69c4bcec8ab0f11
    },
    progress => {
      console.log('progress number:', progress);
    },
  );
}
// 上传图片
// var imgPosY = 0;
// sourcefiles.addEventListener('click', draw);

// function draw() {
//     // 获取选择到的文件
//     let files = sourcefiles.files;

//     for (var file of files) {
//         // console.log(file.webkitRelativePath);   // 显示图片的相对路径

//         // 读取文件内容
//         let reader = new FileReader();
//         reader.readAsDataURL(file);
//         // 读完完成onload再执行function
//         reader.onload = function () {
//             let img = new Image();
//             // 图片内容从result获取
//             img.src = this.result;
//             // 图片读取内容完成后，执行function来添加到画布上
//             img.onload = function() {
//                 let myCanvas = document.getElementById("myCanvas");
//                 let cxt = myCanvas.getContext('2d');
//                 console.log(img.height);
//                 // (img, 图片左上角在画布的横坐标，~的纵坐标)
//                 cxt.drawImage(img, 0, imgPosY);			
//                 imgPosY += img.height;		// 采用接着上一张图的纵向排布
//             }
//         }
//     }
// }


  

