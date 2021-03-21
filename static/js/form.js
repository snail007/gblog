$(document).ready(function () {
    $("form .spinner-grow").each(function () {
        var that = $(this)
        var submit = that.parent();
        var form = that.closest("form")
        submit.click(function () {
            form.attr("spinner-grow", "#" + that.attr("id"))
        });
    });
    var call = function (typ, form, content, status) {
        var done = form.attr("ajax-done");
        var success = form.attr("ajax-success");
        var error = form.attr("ajax-error");
        var stop = false
        switch (typ) {
            case "done":
                if (done) {
                    stop = window[done](form, content, status)
                }
            case "success":
                if (success) {
                    stop = window[success](form, content, status)
                }
            case "error":
                if (error) {
                    stop = window[error](form, content, status)
                }
        }
        return stop;
    }
    $('.ajax-form').submit(function () {
        var form = $(this)
        $(this).ajaxSubmit({
            dataType: "json",
            beforeSend: function () {
                form.find("[type=submit]").attr("disabled", "disabled")
                var loadingID = form.attr("spinner-grow");
                form.find(loadingID).parent().attr("disabled", "disabled")
                form.find(loadingID).show()
            },
            complete: function (a, b) {
                var msg = a.responseText || ""
                if (call("done", form, msg, a.status)) {
                    return
                }
                form.find("[type=submit]").removeAttr("disabled")
                var loadingID = form.attr("spinner-grow");
                form.find(loadingID).parent().removeAttr("disabled")
                form.find(loadingID).hide()
            },
            error: function (a, b, c) {
                var msg = a.responseText || ""
                if (call("error", form, msg, a.status)) {
                    return
                }
                notify("请求错误，请重试。响应码：" + a.status, "danger")
            },
            success: function (data, status, jxhr) {
                var msg = jxhr.responseText || ""
                if (call("success", form, msg, status)) {
                    return
                }
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
    });
    // fileupload
    $(".fileupload").each(function () {
        var loader;
        var $input = $(this);
        var targetInput = $($input.attr("data-input"))
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
                showNotify(errorThrown | textStatus, 'danger');
            },
            done: function () {
                loader.destroy();
            }
        });
    });
});