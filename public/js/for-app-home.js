function search() {
    data = new FormData();
    var key = $("#key").val();
    data.append("key", $('#key').val());
    $.ajax({
        data: data,
        type: "POST",
        url: "/App/Search",
        cache: false,
        processData: false, // 告诉jQuery不要去处理发送的数据
        contentType: false, // 告诉jQuery不要去设置Content-Type请求头
        success: function (data) {
            $('#content').html(data);
        },
        error: function (data) {
            $('#content').html(data);
        }
    });
}