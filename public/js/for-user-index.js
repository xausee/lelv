function strcnt(s, o, n) {
    var reg = new RegExp(o, "g");
    return s.replace(reg, n);
};

function recoverSummernoteCode(id) {
    var c = document.getElementById(id)
    var str = c.innerHTML
    str = strcnt(str, "&amp;", "&");
    //str = strcnt(str, "&nbsp;", " ");
    str = strcnt(str, "&lt;", "<");
    str = strcnt(str, "&gt;", ">");
    str = strcnt(str, "＇", "'");
    //str = strcnt(str, "<br>", "\n");
    c.innerHTML = str;
}

$(document).ready(function() {
    // var c = document.getElementById("introduction")
    // var str = c.innerHTML
    // str = strcnt(str, " ", "&nbsp;");
    // c.innerHTML = str;
    recoverSummernoteCode("content");

    // 点击所有博客的事件响应
    $('#AllBlogs').click(function(e) {
        data = new FormData();
        data.append("UserID", $('#UserID').val());
        $.ajax({
            data: data,
            type: "POST",
            url: "/User/AllBlogs",
            cache: false,
            processData: false,
            contentType: false,
            success: function(data) {
                $('#content').html(data)
            },
            error: function(data) {
                $('#content').html(data);
            }
        });
    });

    // 点击关注的事件响应
    $('#Watch').click(function(e) {
        if ($("#SigninedUserID").val() == "Guest") {
            var html = '<div class="alert alert-success text-center" style="margin-bottom:0;">\
                        <button type="button" class="close" data-dismiss="alert">×</button>\
                        账号未登陆，<a href="/User/SignIn?redirect=/u/' + $('#UserID').val() + '">登陆</a>后关注\
                    </div>'
            var info = document.getElementById("info");
            info.innerHTML = html;
            return;
        }

        data = new FormData();
        data.append("UserID", $('#UserID').val());
        data.append("Flag", $('#Watch').html());
        $.ajax({
            data: data,
            type: "POST",
            url: "/User/Watch",
            cache: false,
            processData: false,
            contentType: false,
            success: function(data) {
                $('#Watch').html(data)
            },
            error: function(data) {
                $('#Watch').html(data);
            }
        });
    });

    // 点击关注的事件响应
    $('#Message').click(function(e) {
        if ($("#SigninedUserID").val() == "Guest") {
            var html = '<div class="alert alert-success text-center" style="margin-bottom:0;">\
                        <button type="button" class="close" data-dismiss="alert">×</button>\
                        账号未登陆，<a href=/User/SignIn>登陆</a>后发送私信\
                    </div>'
            var info = document.getElementById("info");
            info.innerHTML = html;
            return;
        } else {
            window.location.href = "/User/ConversationWith?uid=" + $('#UserID').val();
        }
    });

    // 点击粉丝事件响应
    $('#Fans').click(function(e) {
        console.log($(this).html())
            //e.preventDefault();
        data = new FormData();
        data.append("UserID", $('#UserID').val());
        $.ajax({
            data: data,
            type: "POST",
            url: "/User/Fans",
            cache: false,
            processData: false,
            contentType: false,
            success: function(data) {
                $('#content').html(data)
            },
            error: function(data) {
                $('#content').html(data);
            }
        });
    });

});