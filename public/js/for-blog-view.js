function strcnt(s, o, n) {
    var reg = new RegExp(o, "g");
    return s.replace(reg, n);
};

function RecoverContent() {
    var c = document.getElementById("content")
    var str = c.innerHTML
    str = strcnt(str, "&amp;", "&");
    //str = strcnt(str, "&nbsp;", " ");
    str = strcnt(str, "&lt;", "<");
    str = strcnt(str, "&gt;", ">");
    str = strcnt(str, "＇", "'");
    str = strcnt(str, "<br>", "\n");
    c.innerHTML = str;
}

$(document).ready(function() {
    RecoverContent();

    $("#Delete").click(function() {
        DeleteBlog();
    });

    // $("#comment_Body").focus(function () {
    //     OnCommentInput();
    // });
});

function collect() {
    if ($("#SigninedUserID").val() == "Guest") {
        var html = '<div class="alert alert-success text-center" style="margin-bottom:0;">\
                        <button type="button" class="close" data-dismiss="alert">×</button>\
                        账号未登陆，<a href=/User/SignIn>登陆</a>后收藏\
                    </div><br>'
        var info = document.getElementById("info");
        info.innerHTML = html;
        return;
    }

    data = new FormData();
    data.append("BlogID", $('#BlogID').val());
    data.append("Flag", $('#collect').html());
    $.ajax({
        data: data,
        type: "POST",
        url: "/User/Collect",
        cache: false,
        processData: false, // 告诉jQuery不要去处理发送的数据
        contentType: false, // 告诉jQuery不要去设置Content-Type请求头
        success: function(data) {
            $('#collect').html(data)
        },
        error: function(data) {
            $('#collect').html(data);
        }
    });
}

function OnCommentInput() {
    if ($("#SigninedUserID").val() == "Guest") {
        var html = '<div class="alert alert-success text-center" style="margin-bottom:0;">\
                        <button type="button" class="close" data-dismiss="alert">×</button>\
                        账号未登陆，<a href=/User/SignIn>发表评论</a>\
                    </div><br>'
        var info = document.getElementById("info");
        info.innerHTML = html;
        return;
    }
}

var response

function setInfo() {
    console.log(response)
    $('#content').html('<div class="alert alert-success text-center" role="alert">' + response + '</div>');
}

function Info() {
    $('.close').click();
    setInterval("setInfo()", 1000);
}

// 删除博客
function DeleteBlog() {
    $.ajax({
        type: "GET",
        url: "/Blog/Delete?id=" + $('#BlogID').val(),
        sync: false,
        success: function(data) {
            response = data;
            Info(data);
        },
        error: function(data) {
            response = data;
            Info(data);
        }
    });
}