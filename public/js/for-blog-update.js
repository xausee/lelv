// 博客封面图片
var coverPic = ""
var setCover
// 提交form表单
function post_blog() {
    content = $('#summernote').summernote('code')
    data = new FormData();
    
    data.append("id", $('#BlogID').val());
    data.append("title", $('#BlogTitle').val());
    data.append("tags", $('#TagsInput').val());
    data.append("type", $('input:radio:checked').val());
    data.append("cover", coverPic);
    // 最多截取50个字符
    data.append("briefText", $('.note-editable.panel-body').text().substring(0, 50));
    data.append("content", content);
    $.ajax({
        data: data,
        type: "POST",
        url: "/Blog/PostEdit",
        cache: false,
        processData: false, // 告诉jQuery不要去处理发送的数据
        contentType: false, // 告诉jQuery不要去设置Content-Type请求头
        success: function (data) {
            var form = document.getElementById("BlogForm");
            form.style.display = "none";
            var html = '<div class="alert alert-success text-center" role="alert"><br><br>\
            发布成功&nbsp;&nbsp;<a href=../Blog/View?id=' + data + '>查看</a>&nbsp;&nbsp;\
            <a href="">再写一篇</a><br><br><br></div>'
            var info = document.getElementById("info");
            info.innerHTML = html;
            setCover = false;
            // setTimeout(function() {
            //     info.innerHTML = "";
            // }, 1000 * 10);
        }
    });
}

// 上传图片到七牛云存储
function UploadToQiNiu(file) {
    var info = '<button data-original-title="正在上传" title="正在上传" type="button" class="note-btn btn btn-default btn-sm info"> <img src="../public/img/loading-sm.gif" height="12px"></button>'
    $(".note-btn-group.btn-group.note-insert").append(info);
    var filename = false;
    try {
        filename = file['name'];
    } catch (e) {
        filename = false;
    }
    if (!filename) {
        $(".note-alarm").remove();
    }
    //以上防止在图片在编辑器内拖拽引发第二次上传导致的提示错误
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
        success: function (data) {
            $(".note-btn.btn.btn-default.btn-sm.info").remove();
            //data是返回的hash,key之类的值，key是定义的文件名
            //http://7xsp9p.com1.z0.glb.clouddn.com/: 七牛云domain
            var url = "http://7xsp9p.com1.z0.glb.clouddn.com/" + data["key"];
            if (!setCover) {
                coverPic = url;
                setCover = true;
            }
            // 传到七牛成功后，插入图片到编辑器 
            $('#summernote').summernote('insertImage', url, function ($image) {
                $image.css('width', 1900); //$image.width());
                $image.addClass("img-responsive img-rounded")
                $image.attr('data-filename', 'retriever');
            });
        },
        error: function (res) {
            console.log(res)
            $(".note-btn.btn.btn-default.btn-sm.info").remove();
            var info = '<button data-original-title="上传失败" title="" type="button" class="note-btn btn btn-default btn-sm error">上传失败</button>'
            $(".note-btn-group.btn-group.note-insert").append(info);
            setTimeout(function () {
                $(".note-btn.btn.btn-default.btn-sm.error").remove();
            }, 3000);
        }
    });
}

// 编写文字博客时的工具栏
var textBlogToolbar = [
    ['style', ['style']],
    ['font', ['bold', 'underline', 'clear']],
    ['fontname', ['fontname']],
    ['color', ['color']],
    ['para', ['ul', 'ol', 'paragraph']],
    ['view', ['fullscreen', 'codeview', 'help']]
]

// 编写图文博客时的工具栏
var pictureBlogToolbar = [
    ['insert', ['link', 'picture']],
    ['view', ['fullscreen', 'help']]
]

// 编写混合图文博客时的工具栏
var hybridToolbar = [
    ['style', ['style']],
    ['font', ['bold', 'underline', 'clear']],
    ['fontname', ['fontname']],
    ['color', ['color']],
    ['para', ['ul', 'ol', 'paragraph']],
    ['table', ['table']],
    ['insert', ['link', 'picture']],
    ['view', ['fullscreen', 'codeview', 'help']]
]

// 缺省的工具栏
var defaultToolbar = [
    ['style', ['style']],
    ['font', ['bold', 'underline', 'clear']],
    ['fontname', ['fontname']],
    ['color', ['color']],
    ['para', ['ul', 'ol', 'paragraph']],
    ['table', ['table']],
    ['insert', ['link', 'picture', 'video']],
    ['view', ['fullscreen', 'codeview', 'help']]
]

function strcnt(s, o, n) {
    var reg = new RegExp(o, "g");
    return s.replace(reg, n);
};

function RecoverContent() {
    var c = document.getElementById("OldContent")
    var str = c.innerHTML
    str = strcnt(str, "&amp;", "&");
    //str = strcnt(str, "&nbsp;", " ");
    str = strcnt(str, "&lt;", "<");
    str = strcnt(str, "&gt;", ">");
    str = strcnt(str, "＇", "'");
    str = strcnt(str, "<br>", "\n");
    $('#summernote').summernote('code', str);
}

var toolbar = pictureBlogToolbar;
$(document).ready(function () {
    var type = $('#TagsInput').val();
    switch (type) {
        case "Picture":
            toolbar = pictureBlogToolbar;
            break;
        case "Text":
            toolbar = textBlogToolbar;
            break;
        case "Hibrid":
            toolbar = hybridToolbar;
            break;
        default:
            toolbar = defaultToolbar;
    }
    InitSummernote();
    RecoverContent();

    $("#Picture").click(function () {
        DestroySummernote();
        toolbar = pictureBlogToolbar;
        InitSummernote();
    });

    $("#Text").click(function () {
        DestroySummernote();
        toolbar = textBlogToolbar;
        InitSummernote();
    });

    $("#Hybrid").click(function () {
        DestroySummernote();
        toolbar = hybridToolbar;
        InitSummernote();
    });

    $("#Publish").click(function () {
        post_blog();
    });
});

function InitSummernote() {
    $('#summernote').summernote({
        toolbar: toolbar, // set toolbar
        lang: 'zh-CN', // default: 'en-US'
        height: 500, // set editor height
        minHeight: null, // set minimum height of editor
        maxHeight: null, // set maximum height of editor
        focus: true, // set focus to editable area after initializing summernote
        callbacks: {
            onImageUpload: function (files) {
                console.log(files[0]);
                UploadToQiNiu(files[0])
            }
        }
    });
}

function DestroySummernote() {
    $('#summernote').summernote('destroy');
}

function OnTagLabelsClick() {
    var labels = document.getElementById("TagLabels")
    labels.style.display = "none"

    var tags = document.getElementsByName("Tags")
    var tstr = "";
    for (i = 0; i < tags.length; i++) {
        tstr += tags[i].innerText + " ";
    }
    console.log(tstr)

    var input = document.getElementById("TagsInput")
    input.value = tstr;
    input.style.display = "block"
}

function OnTagInputMouseMoveOut() {
    var input = document.getElementById("TagsInput")
    input.style.display = "none"

    var tags = input.value.split(" ");
    html = "标签：&nbsp;&nbsp;";
    for (i = 0; i < tags.length; i++) {
        color = "default";
        switch (i % 6) {
            case 0:
                color = "default";
                break;
            case 1:
                color = "primary";
                break;
            case 2:
                color = "success";
                break;
            case 3:
                color = "info";
                break;
            case 4:
                color = "warning";
                break;
            case 5:
                color = "danger";
                break;
        }
        if (tags[i] != "")
            html += '<span class="label label-' + color + '" name="Tags">' + tags[i] + '</span>\n';
    }
    var labels = document.getElementById("TagLabels")
    labels.innerHTML = html;
    labels.style.display = "block"
}