function search() {
    data = new FormData();
    var key = $("#key").val();
    data.append("key", $('#key').val());
    $.ajax({
        data: data,
        type: "POST",
        url: "/Admin/SearchForModule",
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
    $("[type=checkbox]").click(function (e) {
        var selectedBlog = "";
        var ids = [];
        $('input:checkbox[name=SearchResults]:checked').each(function (i) {
            var blogid = $(this).val()
            var blog = $("#" + blogid).html();
            ids.push(blogid)
            selectedBlog += blog;
        });
        selectedBlog = "<input style=\"display:none\" id=\"ChosedBlogIDs\" value=" + ids + ">" + selectedBlog;
        $("#ChosedBlogs").html(selectedBlog);
    });
});

function save() {
    data = new FormData();
    data.append("IDs", $("#ChosedBlogIDs").val());
    data.append("BlogType", $("#BlogType").val());

    $.ajax({
        data: data,
        type: "POST",
        url: "/Admin/PostEditHomeBlog",
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