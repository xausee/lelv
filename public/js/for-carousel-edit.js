window.onload = function () {
    var options = {
        imageBox: '.imageBox',
        thumbBox: '.thumbBox',
        spinner: '.spinner',
        imgSrc: '../public/pop/images/icon_alert.gif'
    }
    var cropper = new cropbox(options);
    document.querySelector('#file').addEventListener('change', function () {
        var reader = new FileReader();
        reader.onload = function (e) {
            options.imgSrc = e.target.result;
            cropper = new cropbox(options);
        }
        reader.readAsDataURL(this.files[0]);
        this.files = [];
    })

    // 存在安全性问题，不能执行
    // document.querySelector('#select').addEventListener('click', function () {
    //     var reader = new FileReader();
    //     // 修改为从结果列表中选择博客的方式
    //     var v = $("input[name='searchResults']:checked").val();
    //     options.imgSrc = $("#img_" + v).attr("src");


    //     cropper = new cropbox(options);
    //     var image = new Image();
    //     image.crossOrigin = 'anonymous'
    //     image.src = options.imgSrc;
    //     var base64 = getBase64Image(image);
    //     reader.readAsDataURL(base64);
    //     this.files = [];
    // })

    document.querySelector('#btnCrop').addEventListener('click', function () {
        var img = cropper.getDataURL();
        document.querySelector('.cropbox-cropped').innerHTML = '<img src="' + img + '" id="cropped_avatar">';
        // document.querySelector('.cropped').innerHTML += '<img src="' + img + '">';
    })
    document.querySelector('#btnZoomIn').addEventListener('click', function () {
        cropper.zoomIn();
    })
    document.querySelector('#btnZoomOut').addEventListener('click', function () {
        cropper.zoomOut();
    })
};

// 存在安全性问题，不能执行
// function getBase64Image(img) {
//     var canvas = document.createElement("canvas");
//     canvas.width = img.width;
//     canvas.height = img.height;
//     var ctx = canvas.getContext("2d");
//     ctx.drawImage(img, 0, 0, img.width, img.height);
//     var ext = img.src.substring(img.src.lastIndexOf(".") + 1).toLowerCase();
//     var dataURL = canvas.toDataURL("image/" + ext);
//     return dataURL;
// }
// var img = "https://img.alicdn.com/bao/uploaded/TB1qimQIpXXXXXbXFXXSutbFXXX.jpg";
// var image = new Image();
// image.src = img;
// image.onload = function () {
//     var base64 = getBase64Image(image);
//     console.log(base64);
// }

function save() {
    var n = $("input[name='SearchResults']:checked").val();

    var id = $("#BlogID_" + n).val();
    var title = $("#Title_" + n).html();
    var base64string = $("#cropped_avatar")[0].src;
    if (base64string.contains("base64,")) {
        base64string = base64string.split("base64,")[1];
    }

    data = new FormData();
    data.append("ID", id);
    data.append("Title", title);
    data.append("CoverBase64String", base64string);

    $.ajax({
        data: data,
        type: "POST",
        url: "/Admin/PostCarouselBlog",
        cache: false,
        processData: false, // 告诉jQuery不要去处理发送的数据
        contentType: false, // 告诉jQuery不要去设置Content-Type请求头
        success: function (data) {
            var html = '<div class="alert alert-success text-center" role="alert"><br><br>\
            编辑成功&nbsp;&nbsp;<a href=../Blog/View?id=' + data + '>查看</a>&nbsp;&nbsp;\
            <a href="">继续编辑</a><br><br><br></div>'
            alert("保存成功");
        }
    });
}

function search() {
    data = new FormData();
    var key = $("#key").val();
    data.append("key", $('#key').val());
    $.ajax({
        data: data,
        type: "POST",
        url: "/Admin/SearchForCarousel",
        cache: false,
        processData: false, // 告诉jQuery不要去处理发送的数据
        contentType: false, // 告诉jQuery不要去设置Content-Type请求头
        success: function (data) {
            $('#search_results').html(data);
        },
        error: function (data) {
            $('#search_results').html(data);
        }
    });
}

$(function () {
    $("input[name=options]").click(function () {
        // var v = $("input[name='options']:checked").val();
        // var img = new Image();
        // img.src = $("#img_" + v).attr("src");
        // var w = img.width;
        // var h = img.height;

        // $("#imageBox").css("background-image", "url(" + img.src + ")");
        // $("#imageBox").css("background-size", h + "px " + w + "px");
        // $("#imageBox").css("background-position", "395.355px 195.355px");

        // style="background-image: url(http://localhost:9000/public/pop/images/icon_alert.gif); background-size: 7.290000000000001px 7.290000000000001px; background-position: 395.355px 195.355px; background-repeat: no-repeat"
    });
});