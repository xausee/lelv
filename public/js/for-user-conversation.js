$(document).ready(function() {
    $('#MsgPanel').scrollTop($('#MsgPanel')[0].scrollHeight);
});

$('#PostMessage').click(function(e) {
    PostMessage();
});

$('#MessageContent').keypress(function(e) {
    if (e.charCode == 13 || e.keyCode == 13) {
        $('#PostMessage').click()
        e.preventDefault()
    }
});

function PostMessage() {
    data = new FormData();
    data.append("Content", $('#MessageContent').val());
    data.append("ConversationID", $('#ConversationID').val());
    $('#MessageContent').val("");

    $.ajax({
        type: "Post",
        data: data,
        url: "/User/PostMessage",
        cache: false,
        processData: false,
        contentType: false,
        success: function(data) {
            var content = $('#Messages').html();
            $('#Messages').html(content + data);
            $('#MsgPanel').scrollTop($('#MsgPanel')[0].scrollHeight);
        },
        error: function(data) {
            var content = $('#Messages').html();
            $('#Messages').html(content + data);
        }
    });
}

var LoopPostMessage = function() {
    $.ajax({
        type: "GET",
        url: "/User/GetUnreadMessages?conversationID=" + $('#ConversationID').val(),
        cache: false,
        processData: false,
        contentType: false,
        success: function(data) {
            console.log(data)
            if (data != "") {
                var content = $('#Messages').html();
                $('#Messages').html(content + data);
                $('#MsgPanel').scrollTop($('#MsgPanel')[0].scrollHeight);
            }

            LoopPostMessage()
        },
        error: function(data) {
            var content = $('#Messages').html();
            $('#Messages').html(content + data);
        }
    });
}
LoopPostMessage();