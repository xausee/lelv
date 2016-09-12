function uploadAvatar() {
    var diag = new Dialog();
    // diag.CancelEvent = function() {
    //     alert("点击取消或关闭按钮时执行方法");
    //     diag.close();
    // };

    diag.Title = "修改头像";
    diag.Drag = false;
    diag.Width = 750;
    diag.Height = 500;
    diag.URL = 'Avatar';
    diag.show();
}

// 提交form表单
function postForm() {
    post_profile()
        // putb64()
}

// 提交form表单
function post_profile() {
    data = new FormData();
    var base64string = document.getElementById("avatar").src;
    if (base64string.indexOf("base64,") > 0) {
        data.append("Base64String", base64string.split("base64,")[1]);
    }
    console.log("kokjoj")
    //data.append("NickName", document.getElementById("NickName").value);
    data.append("Introduction", document.getElementById("Introduction").value);

    $.ajax({
        data: data,
        type: "POST",
        url: "/User/PostProfile",
        cache: false,
        processData: false, // 告诉jQuery不要去处理发送的数据
        contentType: false, // 告诉jQuery不要去设置Content-Type请求头
        success: function(data) {
            console.log(data)
            window.location.href = "/User/Home";
        }
    });
}

function putb64() {
    var pic = document.getElementById("avatar").src.split("base64,")[1];
    var url = "http://up.qiniu.com/yixing/-1";
    var xhr = new XMLHttpRequest();

    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4) {
            document.getElementById("feadback").innerHTML = xhr.responseText;
        }
    }
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/octet-stream");
    xhr.setRequestHeader("Authorization", "UpToken " + $("#uptoken").val());
    xhr.send(pic);
}

// 上传图片到七牛云存储
function UploadToQiNiu(file) {
    var timestamp = new Date().getTime();
    var name = timestamp + "_" + filename;
    data = new FormData();
    data.append("file", file);
    data.append("key", name);
    // 七牛token
    data.append("token", $("#uptoken").val());
    $.ajax({
        data: data,
        type: "POST",
        url: "http://upload.qiniu.com/",
        cache: false,
        contentType: false,
        processData: false,
        success: function(data) {
            //data是返回的hash,key之类的值，key是定义的文件名
            //http://7xsp9p.com1.z0.glb.clouddn.com/: 七牛云domain
            var url = "http://7xsp9p.com1.z0.glb.clouddn.com/" + data["key"];
        },
        error: function(res) {
            console.log(res)
        }
    });
}