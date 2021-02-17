$(document).ready(function () {
    $('.ajax-form').ajaxForm({
        dataType: "json",
        error: function (a, b, c) {
            alert("请求错误，请重试。响应码：" + a.status)
        },
        success: function (data) {
            if (data.code == 200) {
                if (data.url) {
                    setTimeout(function () {
                        location = data.url
                    }, 1500)
                }
            } else {
                alert(data.msg || "操作失败，请重试。")
            }
        }
    });
    if(parent!=self){
        parent.location=location;
    }
});