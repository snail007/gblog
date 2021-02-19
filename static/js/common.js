/*
 * 提取通用的通知消息方法
 * 这里只采用简单的用法，如果想要使用回调或者更多的用法，请查看lyear_js_notify.html页面
 * @param $msg 提示信息
 * @param $type 提示类型:'info', 'success', 'warning', 'danger'
 * @param $delay 毫秒数，例如：1000
 * @param $icon 图标，例如：'fa fa-user' 或 'glyphicon glyphicon-warning-sign'
 * @param $from 'top' 或 'bottom' 消息出现的位置
 * @param $align 'left', 'right', 'center' 消息出现的位置
 */
function notify($msg, $type, $delay, $icon, $from, $align) {
    showNotify($msg, $type, $delay, $icon, $from, $align)
}
function showNotify($msg, $type, $delay, $icon, $from, $align) {
    $type = $type || 'info';
    $delay = $delay || 1000;
    $from = $from || 'top';
    $align = $align || 'right';
    $enter = $type == 'danger' ? 'animated shake' : 'animated fadeInUp';

    jQuery.notify({
            icon: $icon,
            message: $msg
        },
        {
            element: 'body',
            type: $type,
            allow_dismiss: true,
            newest_on_top: true,
            showProgressbar: false,
            placement: {
                from: $from,
                align: $align
            },
            offset: 20,
            spacing: 10,
            z_index: 10800,
            delay: $delay,
            animate: {
                enter: $enter,
                exit: 'animated fadeOutDown'
            }
        });
}
