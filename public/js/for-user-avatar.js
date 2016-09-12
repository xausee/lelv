window.onload = function () {
    var options = {
        imageBox: '.imageBox',
        thumbBox: '.thumbBox',
        spinner: '.spinner',
        imgSrc: '../public/img/icon_alert.gif'
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

function save() {
    var src = $("#cropped_avatar")[0].src;
    window.parent.document.getElementById("avatar").src = src;

    window.parent.document.getElementById("_ButtonCancel_0").click();
}