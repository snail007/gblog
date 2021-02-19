$(document).ready(function () {
    $('.ajax-form').submit(function () {
        var form = $(this)
        $(this).ajaxSubmit({
            dataType: "json",
            beforeSend: function () {
                form.find(".spinner-grow").parent().attr("disabled", "disabled")
                form.find(".spinner-grow").show()
            },
            complete: function (a, b) {
                form.find(".spinner-grow").parent().removeAttr("disabled")
                form.find(".spinner-grow").hide()
            },
            error: function (a, b, c) {
                notify("请求错误，请重试。响应码：" + a.status, "danger")
            },
            success: function (data) {
                if (data.code == 200) {
                    notify(data.msg || "操作成功！", "success")
                    if (data.url) {
                        setTimeout(function () {
                            location = data.url
                        }, 1500)
                    }
                } else {
                    notify(data.msg || "操作失败，请重试。", "warning")
                }
            }
        })
        return false;
    })
    // fileupload
    $(".fileupload").each(function () {
        var loader;
        var $input=$(this);
        var targetInput=$($input.attr("data-input"))
        $input.fileupload({
            dataType: 'json',
            beforeSend: function () {
                loader = $('body').lyearloading({
                    opacity: 0.2,
                    spinnerSize: 'lg'
                });
            },
            success: function (result, textStatus, jqXHR) {
                showNotify("上传成功", 'success');
                targetInput.val(result.url)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                showNotify(errorThrown|textStatus, 'danger');
            },
            done: function () {
                loader.destroy();
            }
        });
    });
});