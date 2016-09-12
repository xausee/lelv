$(document).ready(function () {
    $("#Submit").click(function () {

    });

    $("#Cancel").click(function () {
        cancel();
    });
})

function cancel() {
    var uid = $("#RemoteUserID").val();
    console.log(uid)
    history.go(-1)
}

function postMessage(){
      data = new FormData();
        data.append("UserID", $('#UserID').val());
        $.ajax({
            data: data,
            type: "POST",
            url: "/User/AllBlogs",
            cache: false,
            processData: false,
            contentType: false,
            success: function (data) {
                $('#content').html(data)
            },
            error: function (data) {
                $('#content').html(data);
            }
        });
}