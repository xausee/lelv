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
        $.ajax({
            type: "GET",
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

    // 点击所有博客的事件响应
    $('#Collection').click(function(e) {
        $.ajax({
            type: "GET",
            url: "/User/Collection",
            success: function(data) {
                $('#content').html(data)
            },
            error: function(data) {
                $('#content').html(data);
            }
        });
    });

    // 点击粉丝事件响应
    $('#Watch').click(function(e) {
        $.ajax({
            type: "GET",
            url: "/User/Watches",
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

    // 点击粉丝事件响应
    $('#Fans').click(function(e) {
        $.ajax({
            type: "GET",
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

    // 点击消息事件响应
    //$('[name=Conversations]') $('#Conversations')
    $('[name=Conversations]').click(function(e) {
        $.ajax({
            type: "GET",
            url: "/User/Conversations",
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